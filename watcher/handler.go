package watcher

import (
	"encoding/json"
	"github.com/synw/centcom"
	"github.com/synw/microb-http/httpServer"
	"github.com/synw/microb/libmicrob/events"
)

func handle(cli *centcom.Cli, channel string) {
	httpServer.ParseTemplates()
	p := map[string]string{"reload": "true"}
	payload, err := json.Marshal(p)
	if err != nil {
		events.Err("http", "watcher.handler.handle", "error encoding payload", err)
	}
	cli.Publish(channel, payload)
	_, err = cli.Http.Publish(channel, payload)
	if err != nil {
		events.Err("http", "watcher.handler.handle", "error publishing to channel", err)
	}
}
