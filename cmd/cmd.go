package cmd

import (
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb-http/state"
	"github.com/synw/microb-http/state/mutate"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/events"
	datatypes "github.com/synw/microb/types"
)

func Dispatch(cmd *datatypes.Command) *datatypes.Command {
	com := &datatypes.Command{}
	if cmd.Name == "start" {
		res := start(cmd, state.HttpServer)
		return res
	} else if cmd.Name == "stop" {
		return stop(cmd, state.HttpServer)
	} else if cmd.Name == "parse_templates" {
		return parseTemplates(cmd)
	}
	return com
}

func parseTemplates(cmd *datatypes.Command) *datatypes.Command {
	httpServer.ParseTemplates()
	var resp []interface{}
	resp = append(resp, "Templates parsed")
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}

func start(cmd *datatypes.Command, server *types.HttpServer) *datatypes.Command {
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

func stop(cmd *datatypes.Command, server *types.HttpServer) *datatypes.Command {
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
