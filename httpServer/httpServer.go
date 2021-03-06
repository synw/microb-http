package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb-mail/mail"
	"github.com/synw/microb/events"
	"github.com/synw/microb/msgs"
	"github.com/synw/terr"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var dev bool
var conn *types.Conn
var key string
var domain string
var edit_channel string
var datasource *types.Datasource
var templates *template.Template
var basePath string = conf.GetBasePath()
var csrfKey string

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

type authRequest struct {
	Client   string   `json:client`
	Channels []string `json:channels`
}

func Init(server *types.HttpServer, ws bool, addr string, key string, dm string, ec string, ds *types.Datasource, isDev bool, isMail bool, icsrfKey string, hitsDb string) {
	domain = dm
	datasource = ds
	edit_channel = ec
	dev = isDev
	csrfKey = icsrfKey
	if ws {
		initWs(addr, key)
	}
	// templates init
	templates, _ = parseTemplates()
	// routing
	r := chi.NewRouter()
	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	path := basePath + "/static"
	// static
	fileServer(r, "/static", http.Dir(path))
	// init hits db
	tr := initHitsDb(hitsDb, domain)
	if tr != nil {
		events.Fatal("httpServer.Init", "Can not initialize hits database", tr)
	}
	// routes
	// mail service
	if isMail == true {
		r.Route("/mail", func(r chi.Router) {
			r.Post("/post", mail.ProcessMailForm)
			r.Get("/ok", serveRequest)
			r.Get("/", mail.ServeMailForm)
		})
	}
	// authentication for edit channel
	r.Route("/centrifuge", func(r chi.Router) {
		r.Post("/auth", serveAuth)
	})
	// pages
	r.Route("/", func(r chi.Router) {
		r.Get("/*", serveRequest)
	})
	// init http
	httpServer := &http.Server{
		Addr:         server.Addr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      r,
	}
	server.Instance = httpServer
}

func Run(server *types.HttpServer) {
	events.State("http", startMsg(server))
	http.ListenAndServe(":8080", server.Instance.Handler)
}

func Stop(server *types.HttpServer) *terr.Trace {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := server.Instance
	err := srv.Shutdown(ctx)
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	events.State("http", stopMsg())
	return nil
}

func ParseTemplates() {
	_, _ = parseTemplates()
}

// internal methods

func getToken(user string, timestamp string, secret string) string {
	info := ""
	token := auth.GenerateClientToken(secret, user, timestamp, info)
	return token
}

func initWs(addr string, k string) {
	key = k
	user := "microb_http"
	timestamp := centrifuge.Timestamp()
	token := getToken(user, timestamp, key)
	conn = &types.Conn{addr, timestamp, user, token}
}

func parseTemplates() (*template.Template, *terr.Trace) {
	path := basePath + "/templates/*"
	tmps := template.Must(template.ParseGlob(path))
	templates = tmps
	return tmps, nil
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit URL parameters.")
	}
	fs := http.StripPrefix(path, http.FileServer(root))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

func serveRequest(resp http.ResponseWriter, r *http.Request) {
	url := r.URL.Path
	status := http.StatusOK
	resp = httpResponseWriter{resp, &status}
	go processHit(r, status)
	page, tr := getPage(domain, url, conn, edit_channel)
	if tr != nil {
		tr = tr.Pass("serveRequest")
		events.Warning("http", "Error retrieving page "+url, tr)
		p := &types.Page{}
		render404(resp, p)
		return
	}
	renderTemplate(resp, page)
}

func renderTemplate(resp http.ResponseWriter, page *types.Page) {
	page.Content = template.HTML(page.Content)
	err := templates.ExecuteTemplate(resp, "index.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering template: " + err.Error()
		tr := terr.New("Can not render template")
		events.Error("http", msg, tr)
	}
}

func render404(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "404.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 404"
		err := errors.New(msg)
		tr := terr.New(err)
		events.Error("http", msg, tr)
	}
}

func render500(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "500.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 500"
		err := errors.New(msg)
		tr := terr.New(err)
		events.Error("http", msg, tr)
	}
}

func serveAuth(w http.ResponseWriter, r *http.Request) {
	// this is used only for autoreload in dev
	if dev == false {
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data authRequest
	err := decoder.Decode(&data)
	if err != nil {
		msg := "Error decoding data"
		err := errors.New(msg)
		tr := terr.New(err)
		events.Error("http", msg, tr)
	}
	resp := map[string]map[string]string{}
	for _, channel := range data.Channels {
		client := data.Client
		info := ""
		sign := auth.GenerateChannelSign(key, client, channel, info)
		s := map[string]string{
			"sign": sign,
			"info": info,
		}
		resp[channel] = s
	}
	w.Header().Set("Content-Type", "application/json")
	json_bytes, err := json.Marshal(resp)
	if err != nil {
		msg := "Can not marshall json"
		err := errors.New(msg)
		tr := terr.New(err)
		events.Error("http", msg, tr)
	}
	fmt.Fprintf(w, "%s\n", json_bytes)
}

func stopMsg() string {
	msg := "Http server stopped"
	return msg
}

func startMsg(server *types.HttpServer) string {
	var msg string
	msg = "Http server started at " + server.Addr + " for domain "
	msg = msg + msgs.Bold(server.Domain)
	return msg
}

func getTemplate(name string) string {
	t := basePath + "/templates/" + name + ".html"
	return t
}
