package mutate

import (
	"errors"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
)

func StartHttpServer(server *types.HttpServer) *terr.Trace {
	if server.State.Current() == "start" {
		err := errors.New("The http server is already running")
		tr := terr.New("state.mutate.StartHttpServer", err)
		return tr
	}
	go httpServer.Run(server)
	err := server.State.Event("start")
	if err != nil {
		tr := terr.New("state.mutate.StartHttpServer", err)
		return tr
	}
	return nil
}

func StopHttpServer(server *types.HttpServer) *terr.Trace {
	if server.State.Current() == "stop" {
		err := errors.New("The http server is not running")
		tr := terr.New("state.mutate.StopHttpServer", err)
		return tr
	}
	tr := httpServer.Stop(server)
	if tr != nil {
		return tr
	}
	err := server.State.Event("stop")
	if err != nil {
		tr := terr.New("state.mutate.StopHttpServer", err)
		return tr
	}
	return nil
}
