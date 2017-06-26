package info

import (
	"github.com/synw/microb/libmicrob/datatypes"
	"github.com/synw/microb/services/httpServer/state"
)

func Http(cmd *datatypes.Command) *datatypes.Command {
	var resp []interface{}
	if state.HttpServer.Running == true {
		resp = append(resp, "Http server is up")
	} else {
		resp = append(resp, "Http server is down")
	}
	cmd.Status = "success"
	cmd.ReturnValues = resp
	return cmd
}
