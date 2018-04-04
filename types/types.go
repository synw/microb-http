package types

import (
	"github.com/jinzhu/gorm"
	"github.com/looplab/fsm"
	"html/template"
	"net/http"
)

type HttpServer struct {
	Domain   string
	Addr     string
	Instance *http.Server
	State    *fsm.FSM
}

type Datasource struct {
	Name string
	Type string
	Path string
	User string
	Pwd  string
	Db   string
}

type Conf struct {
	Domain     string
	Addr       string
	WsAddr     string
	WsKey      string
	Ws         bool
	Datasource *Datasource
	EditChan   string
	Dev        bool
	Mail       bool
	CsrfKey    string
	HitsDbAddr string
}

type Conn struct {
	Addr      string
	Timestamp string
	User      string
	Token     string
}

type Hit struct {
	gorm.Model
	Ip     string `json:ip`
	Path   string `json:path`
	Host   string `json:host`
	Ua     string `json:ua`
	Lang   string `json:lang`
	Length string `json:length`
	Status int    `json:status`
}

type Page struct {
	Domain   string
	Url      string
	Title    string
	Content  template.HTML
	Conn     *Conn
	EditChan string
	Token    string
}
