package types

import (
	"net/http"
)

type HttpServer struct {
	Domain   string
	Addr     string
	Instance *http.Server
	Running  bool
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
}

type Conn struct {
	Addr      string
	Timestamp string
	User      string
	Token     string
}

type Page struct {
	Domain   string
	Url      string
	Title    string
	Content  string
	Conn     *Conn
	EditChan string
}
