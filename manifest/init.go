package manifest

import (
	"github.com/synw/microb/libmicrob/types"
)

var Service *types.Service = &types.Service{
	"http",
	getCmds(),
	initService,
}
