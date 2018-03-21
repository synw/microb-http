package manifest

import (
	"github.com/synw/microb-http/state"
	"github.com/synw/microb-http/state/mutate"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/types"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["start"] = start()
	return cmds
}

func initService(dev bool, start bool) error {
	tr := state.Init(dev, start)
	if tr != nil {
		return tr.ToErr()
	}
	return nil
}

func start() *types.Cmd {
	cmd := &types.Cmd{Name: "start", Service: "http", Exec: runStart}
	return cmd
}

func runStart(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	server := state.HttpServer
	tr := mutate.StartHttpServer(server)
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		events.Terr("http", "cmd.Start", "Error starting http service", tr)
		c <- cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	c <- cmd
}

/*
func http() *types.Cmd {
	cmd := &types.Cmd{Name: "http", Exec: runHttp}
	return cmd
}

func runHttp(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	server := state.HttpServer
	arg := cmd.Args[0].(string)
	if arg == "start" {
		tr := mutate.StartHttpServer(server)
		if tr != nil {
			cmd.Trace = tr
			cmd.Status = "error"
			events.Terr("http", "cmd.Start", "Error starting http service", tr)
			c <- cmd
			return
		}
		var resp []interface{}
		resp = append(resp, "Http server started")
		cmd.Status = "success"
		cmd.ReturnValues = resp
	} else if arg == "stop" {
		tr := mutate.StopHttpServer(server)
		if tr != nil {
			cmd.Trace = tr
			cmd.Status = "error"
			events.Terr("http", "cmd.Stop", "Error stoping http service", tr)
			c <- cmd
			return
		}
		var resp []interface{}
		resp = append(resp, "Http server stop")
		cmd.Status = "success"
		cmd.ReturnValues = resp
	}
	c <- cmd
}*/
