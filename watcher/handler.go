package watcher

import (
	"encoding/json"
	"github.com/synw/centcom"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb/events"
	"github.com/synw/terr"
)

func handle(cli *centcom.Cli, channel string, dev bool) {
	httpServer.ParseTemplates()
	if dev == false {
		return
	}
	p := map[string]string{"reload": "true"}
	payload, err := json.Marshal(p)
	if err != nil {
		tr := terr.New(err)
		events.Error("http", "Error encoding payload", tr)
	}
	cli.Publish(channel, payload)
	_, err = cli.Http.Publish(channel, payload)
	if err != nil {
		tr := terr.New(err)
		events.Error("http", "watcher.handler.handle", tr)
	}
}
