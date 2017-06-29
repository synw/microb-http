package state

import (
	"fmt"
	"github.com/synw/centcom"
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb-http/watcher"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"net/http"
)

var HttpServer = &types.HttpServer{}
var Conf *types.Conf
var cli *centcom.Cli

func Init(dev bool, verbosity int, start bool) *terr.Trace {
	Conf, tr := conf.GetConf(dev)
	if tr != nil {
		events.Terr("http", "state.InitState", "Unable to init http server config", tr)
		return tr
	}
	// htto server
	instance := &http.Server{}
	running := false
	HttpServer = &types.HttpServer{Conf.Domain, Conf.Addr, instance, running}
	httpServer.Init(HttpServer, Conf.Ws, Conf.WsAddr, Conf.WsKey, Conf.Domain, start, Conf.EditChan, Conf.Datasource)
	// init ws cli
	cli = centcom.NewClient(Conf.WsAddr, Conf.WsKey)
	err := centcom.Connect(cli)
	if err != nil {
		tr := terr.New("state.Init", err)
		return tr
	}
	err = cli.CheckHttp()
	if err != nil {
		tr := terr.New("state.Init", err)
		return tr
	}
	// watcher for hot reload
	if dev == true {
		if verbosity > 1 {
			fmt.Println("Initializing files watcher ...")
		}
		go watcher.Start(Conf.Datasource.Path, cli, Conf.EditChan, verbosity, dev)
	}
	return nil
}
