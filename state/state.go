package state

import (
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"net/http"
)

var HttpServer = &types.HttpServer{}
var Conf *types.Conf

func InitState(dev bool, verbosity int) *terr.Trace {
	Conf, tr := conf.GetConf(dev)
	if tr != nil {
		events.Terr("http", "state.InitState", "Unable to init http server config", tr)
		return tr
	}
	instance := &http.Server{}
	running := false
	HttpServer = &types.HttpServer{Conf.Domain, Conf.Addr, instance, running}
	httpServer.InitHttpServer(HttpServer, Conf.WsAddr, Conf.WsKey, Conf.Domain, false)
	return nil
}
