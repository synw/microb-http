package httpServer

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"html/template"
	"net/http"
	"strings"
	"time"
)

var isdev = false

type httpResponseWriter struct {
	http.ResponseWriter
	status *int
}

type authRequest struct {
	Client   string   `json:client`
	Channels []string `json:channels`
}

var conn *types.Conn
var key string
var domain string
var edit_channel string
var datasource *types.Datasource
var templates *template.Template
var basePath string = conf.GetBasePath()

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
	/*if err != nil {
		msg := "Can not find templates directory"
		tr := terr.New("httpServer.Stop", err)
		events.New("error", "http", "httpServer.parseTemplates", msg, tr)
		return nil, tr
	}*/
	tmps := template.Must(template.ParseGlob(path))
	templates = tmps
	return tmps, nil
}

func Init(server *types.HttpServer, ws bool, addr string, key string, dm string, ec string, ds *types.Datasource, dev bool) {
	isdev = dev
	domain = dm
	datasource = ds
	edit_channel = ec
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
	// routes
	r.Route("/centrifuge", func(r chi.Router) {
		r.Post("/auth", serveAuth)
	})
	r.Route("/", func(r chi.Router) {
		r.Get("/*", serveRequest)
	})
	// init
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
	server.Instance.ListenAndServe()
}

func Stop(server *types.HttpServer) *terr.Trace {
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	srv := server.Instance
	err := srv.Shutdown(ctx)
	if err != nil {
		tr := terr.New("httpServer.Stop", err)
		events.Error("http", stopMsg(), tr)
		return tr
	}
	events.State("http", stopMsg())
	return nil
}

func ParseTemplates() {
	_, _ = parseTemplates()
}

// internal methods

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

func handleRequest(req *http.Request) {
	/*path := req.URL.Path
	host := req.URL.Host
	header := req.Header
	ua := header["User-Agent"][0]
	lang := header["Accept-Language"][0]
	cl := header["ContentLength"]
	fmt.Println(path, host)
	fmt.Println(ua, lang, cl)*/
}

func serveRequest(resp http.ResponseWriter, req *http.Request) {
	handleRequest(req)

	url := req.URL.Path
	status := http.StatusOK
	resp = httpResponseWriter{resp, &status}
	page, tr := getPage(domain, url, conn, edit_channel)
	if tr != nil {
		tr = terr.Pass("serveRequest", tr)
		events.Error("http", "Error retrieving page "+url, tr, "warn")
		p := &types.Page{}
		render404(resp, p)
		tr.Print()
		return
	}
	renderTemplate(resp, page)
}

func renderTemplate(resp http.ResponseWriter, page *types.Page) {
	page.Content = template.HTML(page.Content)
	err := templates.ExecuteTemplate(resp, "index.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering template"
		err := errors.New("Can not render template")
		tr := terr.New("httpServer.renderTemplate", err)
		events.Error("http", msg, tr)
	}
}

func render404(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "404.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 404"
		err := errors.New(msg)
		tr := terr.New("httpServer.render400", err)
		events.Error("http", msg, tr)
	}
}

func render500(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "500.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 500"
		err := errors.New(msg)
		tr := terr.New("httpServer.render500", err)
		events.Error("http", msg, tr)
	}
}

func serveAuth(resp http.ResponseWriter, req *http.Request) {
	// this is used only for autoreload in dev
	if isdev == false {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var data authRequest
	err := decoder.Decode(&data)
	if err != nil {
		msg := "Error decoding data"
		err := errors.New(msg)
		tr := terr.New("httpServer.serveAuth", err)
		events.Error("http", msg, tr)
	}
	r := map[string]map[string]string{}
	for _, channel := range data.Channels {
		client := data.Client
		info := ""
		sign := auth.GenerateChannelSign(key, client, channel, info)
		s := map[string]string{
			"sign": sign,
			"info": info,
		}
		r[channel] = s
	}
	resp.Header().Set("Content-Type", "application/json")
	json_bytes, err := json.Marshal(r)
	if err != nil {
		msg := "Can not marshall json"
		err := errors.New(msg)
		tr := terr.New("httpServer.serveAuth", err)
		events.Error("http", msg, tr)
	}
	fmt.Fprintf(resp, "%s\n", json_bytes)
}

func stopMsg() string {
	msg := "Http server stopped"
	return msg
}

func startMsg(server *types.HttpServer) string {
	var msg string
	msg = "Http server started at " + server.Addr + " for domain " + skittles.BoldWhite(server.Domain)
	return msg
}

func getTemplate(name string) string {
	t := basePath + "/templates/" + name + ".html"
	return t
}
