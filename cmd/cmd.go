package cmd

import (
	"github.com/synw/microb-http/state"
	"github.com/synw/microb-http/state/mutate"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	datatypes "github.com/synw/microb/libmicrob/types"
)

func Dispatch(cmd *datatypes.Command) *datatypes.Command {
	com := &datatypes.Command{}
	if cmd.Name == "start" {
		res := Start(cmd, state.HttpServer)
		return res
	} else if cmd.Name == "stop" {
		return Stop(cmd, state.HttpServer)
	}
	return com
}

func Start(cmd *datatypes.Command, server *types.HttpServer) *datatypes.Command {
	tr := mutate.StartHttpServer(server)
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		events.Terr("http", "cmd.Start", "Error starting http service", tr)
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server started")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func Stop(cmd *datatypes.Command, server *types.HttpServer) *datatypes.Command {
	tr := mutate.StopHttpServer(server)
	if tr != nil {
		cmd.Trace = tr
		cmd.Status = "error"
		events.Terr("http", "cmd.Stop", "Error stopping http service", tr)
		return cmd
	}
	var resp []interface{}
	resp = append(resp, "Http server stopped")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
