package mutate

import (
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
)

func StartHttpServer(server *types.HttpServer) *terr.Trace {
	if server.State.Current() == "start" {
		tr := terr.New("The http server is already running", "warning")
		return tr
	}
	go httpServer.Run(server)
	err := server.State.Event("start")
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	return nil
}

func StopHttpServer(server *types.HttpServer) *terr.Trace {
	if server.State.Current() == "stop" {
		tr := terr.New("The http server is not running", "warning")
		return tr
	}
	tr := httpServer.Stop(server)
	if tr != nil {
		tr.Pass()
		return tr
	}
	err := server.State.Event("stop")
	if err != nil {
		tr := terr.New(err)
		return tr
	}
	return nil
}
