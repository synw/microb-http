package manifest

import (
	"github.com/synw/microb-http/cmd"
	"github.com/synw/microb-http/state"
	"github.com/synw/microb/libmicrob/types"
	"github.com/synw/terr"
)

var Service *types.Service = &types.Service{
	"http",
	[]string{"start", "stop"},
	ini,
	dispatch,
}

func ini(dev bool, verbosity int) *terr.Trace {
	return state.Init(dev, verbosity)
}

func dispatch(c *types.Command) *types.Command {
	return cmd.Dispatch(c)
}
