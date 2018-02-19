package httpServer

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/acmacalister/skittles"
	"github.com/centrifugal/centrifuge-go"
	"github.com/centrifugal/centrifugo/libcentrifugo/auth"
	"github.com/pressly/chi"
	"github.com/pressly/chi/middleware"
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

var config, _ = conf.GetConf(true)

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
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	path = path + "/templates/*"
	if err != nil {
		msg := "Can not find templates directory"
		tr := terr.New("httpServer.Stop", err)
		events.New("error", "http", "httpServer.parseTemplates", msg, tr)
		return nil, tr
	}
	tmps := template.Must(template.ParseGlob(path))
	templates = tmps
	return tmps, nil
}

func Init(server *types.HttpServer, ws bool, addr string, key string, dm string, serve bool, ec string, ds *types.Datasource) {
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
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	path = path + "/static"
	if err != nil {
		msg := "Can not find static directory"
		tr := terr.New("httpServer.Init", err)
		events.New("error", "http", "httpServer.Init", msg, tr)
		return
	}
	// static
	r.FileServer("/static", http.Dir(path))
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
	// run
	if serve == true {
		go Run(server)
	}
}

func Run(server *types.HttpServer) {
	events.New("state", "http", "httpServer.Run", startMsg(server), nil)
	server.Running = true
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
		events.New("error", "http", "httpServer.Stop", stopMsg(), tr)
		return tr
	}
	server.Running = false
	events.New("state", "http", "httpServer.Stop", stopMsg(), nil)
	return nil
}

func ParseTemplates() {
	_, _ = parseTemplates()
}

// internal methods

func serveRequest(resp http.ResponseWriter, req *http.Request) {
	url := req.URL.Path
	status := http.StatusOK
	resp = httpResponseWriter{resp, &status}
	page, tr := getPage(domain, url, conn, edit_channel)
	if tr != nil {
		events.Terr("http", "httpServer.serveRequest", "Error retrieving page", tr)
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
		msg := "Error rendering template"
		events.Err("http", "httpServer.renderTemplate", msg, err)
	}
}

func render404(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "404.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 404"
		events.Err("http", "httpServer.render404", msg, err)
	}
}

func render500(resp http.ResponseWriter, page *types.Page) {
	err := templates.ExecuteTemplate(resp, "500.html", page)
	if err != nil {
		http.Error(resp, err.Error(), http.StatusInternalServerError)
		msg := "Error rendering 500"
		events.Err("http", "httpServer.render500", msg, err)
	}
}

func serveAuth(resp http.ResponseWriter, req *http.Request) {
	// this is used only for autoreload in dev
	if config.Dev != true {
		return
	}
	decoder := json.NewDecoder(req.Body)
	var data authRequest
	err := decoder.Decode(&data)
	if err != nil {
		msg := "Error decoding data"
		events.Err("http", "httpServer.serveAuth", msg, err)
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
		fmt.Println(err)
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

/*
func getDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	d := fmt.Sprintf("%s", path.Dir(filename))
	d = strings.Replace(d, "/httpServer", "", -1)
	return d
}*/

func getTemplate(name string) string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic("Error")
	}
	t := dir + "/templates/" + name + ".html"
	return t
}
