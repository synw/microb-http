package mutate

import (
	"errors"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
)

func StartHttpServer(server *types.HttpServer) *terr.Trace {
	if server.Running == true {
		err := errors.New("The http server is already running")
		tr := terr.New("state.mutate.StartHttpServer", err)
		return tr
	}
	go httpServer.Run(server)
	return nil
}

func StopHttpServer(server *types.HttpServer) *terr.Trace {
	if server.Running == false {
		err := errors.New("The http server is not running")
		tr := terr.New("state.mutate.StopHttpServer", err)
		return tr
	}
	tr := httpServer.Stop(server)
	if tr != nil {
		return tr
	}
	return nil
}
