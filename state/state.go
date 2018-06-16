package state

import (
	"github.com/looplab/fsm"
	"github.com/synw/centcom"
	"github.com/synw/microb-http/conf"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/state/mutate"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb-http/watcher"
	"github.com/synw/microb/events"
	"github.com/synw/terr"
	"net/http"
)

var HttpServer = &types.HttpServer{}
var Conf *types.Conf
var cli *centcom.Cli
var BasePath string = conf.GetBasePath()

func Init(dev bool, start bool) *terr.Trace {
	Conf, tr := conf.GetConf()
	if tr != nil {
		tr = tr.Pass()
		events.Error("http", "Unable to init http server config", tr)
		return tr
	}
	// htto server
	instance := &http.Server{}
	var st *fsm.FSM
	HttpServer = &types.HttpServer{Conf.Domain, Conf.Addr, instance, st}
	httpServer.Init(HttpServer, Conf.Ws, Conf.WsAddr, Conf.WsKey, Conf.Domain, Conf.EditChan, Conf.Datasource, Conf.Dev, Conf.Mail, Conf.CsrfKey, Conf.HitsDbAddr)
	// init ws cli
	cli = centcom.NewClient(Conf.WsAddr, Conf.WsKey)
	err := centcom.Connect(cli)
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	err = cli.CheckHttp()
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	// initialize the state machine
	HttpServer.State = setState()
	if start == true {
		tr := mutate.StartHttpServer(HttpServer)
		if tr != nil {
			tr := terr.New(err)
			return tr
		}
	}
	// watcher for templates reparsing on change
	// and hot reload for development
	go watcher.Start(BasePath, Conf.Datasource.Path, cli, Conf.EditChan, dev)
	return nil
}

func setState() *fsm.FSM {
	st := fsm.NewFSM(
		"stop",
		fsm.Events{
			{Name: "start", Src: []string{"stop"}, Dst: "start"},
			{Name: "stop", Src: []string{"start"}, Dst: "stop"},
		},
		fsm.Callbacks{},
	)
	return st
}
