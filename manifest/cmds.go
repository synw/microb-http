package manifest

import (
	"errors"
	"github.com/synw/microb-http/state"
	"github.com/synw/microb-http/state/mutate"
	//"github.com/synw/microb/events"
	"github.com/synw/microb/types"
	"github.com/synw/terr"
)

func getCmds() map[string]*types.Cmd {
	cmds := make(map[string]*types.Cmd)
	cmds["start"] = start()
	cmds["stop"] = stop()
	cmds["status"] = status()
	return cmds
}

func initService(dev bool, start bool) *terr.Trace {
	tr := state.Init(dev, start)
	if tr != nil {
		tr = tr.Pass()
		return tr
	}
	return nil
}

func status() *types.Cmd {
	cmd := &types.Cmd{Name: "status", Service: "http", Exec: runStatus, NoLog: true}
	return cmd
}

func runStatus(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	var resp []interface{}
	var msg string
	if state.HttpServer.State.Current() == "start" {
		msg = "The http server is running"
	} else {
		msg = "The http server is not running"
	}
	resp = append(resp, msg)
	cmd.Status = "success"
	cmd.ReturnValues = resp
	c <- cmd
}

func start() *types.Cmd {
	cmd := &types.Cmd{Name: "start", Service: "http", Exec: runStart}
	return cmd
}

func runStart(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	server := state.HttpServer
	if state.HttpServer.State.Current() == "start" {
		msg := "The http server is already started"
		err := errors.New(msg)
		tr := terr.New(err)
		//events.Error("http", msg, tr, "warning")
		cmd.Status = "error"
		cmd.Trace = tr
		cmd.ErrMsg = tr.Error()
		c <- cmd
		return
	}
	tr := mutate.StartHttpServer(server)
	if tr != nil {
		msg := "Error starting http service"
		//events.Error("http", msg, tr, "warning")
		cmd.Status = "error"
		cmd.ErrMsg = msg
		cmd.Trace = tr
		c <- cmd
		return
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	c <- cmd
}

func stop() *types.Cmd {
	cmd := &types.Cmd{Name: "stop", Service: "http", Exec: runStop}
	return cmd
}

func runStop(cmd *types.Cmd, c chan *types.Cmd, args ...interface{}) {
	server := state.HttpServer
	if state.HttpServer.State.Current() == "stop" {
		msg := "The http server is not running"
		err := errors.New(msg)
		tr := terr.New(err)
		cmd.Status = "error"
		cmd.Trace = tr
		c <- cmd
		return
	}
	tr := mutate.StopHttpServer(server)
	if tr != nil {
		msg := "Error stoping http service"
		cmd.Status = "error"
		cmd.ErrMsg = msg
		cmd.Trace = tr
		c <- cmd
		return
	}
	var resp []interface{}
	resp = append(resp, "Http server stopped")
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
