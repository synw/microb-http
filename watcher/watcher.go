package watcher

import (
	wa "github.com/radovskyb/watcher"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/terr"
	"time"
)

var w = wa.New()

func Start(basePath string, path string, cli *centcom.Cli, channel string, dev bool) {
	msgs.Status("Initializing files watcher")
	iserr := false
	err := w.AddRecursive(basePath + "/templates")
	if err != nil {
		tr := terr.New("watcher.Start", err)
		events.Error("http", "Error finding templates", tr)
		iserr = true
	}
	err = w.AddRecursive(basePath + "/static/js")
	if err != nil {
		tr := terr.New("watcher.Start", err)
		events.Error("http", "Error finding static/js", tr)
		iserr = true
	}
	err = w.AddRecursive(basePath + "/static/css")
	if err != nil {
		tr := terr.New("watcher.Start", err)
		events.Error("http", "Error finding static/css", tr)
		iserr = true
	}
	err = w.AddRecursive(basePath + "/static/content")
	if err != nil {
		tr := terr.New("watcher.Start", err)
		events.Error("http", "Error finding static/content", tr)
		iserr = true
	}
	w.FilterOps(wa.Write, wa.Create, wa.Move, wa.Remove, wa.Rename)
	if iserr == false {
		msgs.State("Dev mode is enabled: watching files for change")
		/*for path, f := range w.WatchedFiles() {
			fmt.Printf("%s %s\n", f.Name(), path)
		}*/
	}
	// lauch listener
	go func() {
		for {
			select {
			case e := <-w.Event:
				msgs.Msg("Change detected in " + e.Path + ": reloading")
				handle(cli, channel, dev)
			case err := <-w.Error:
				msgs.Msg("Watcher error " + err.Error())
			case <-w.Closed:
				msgs.Msg("Watcher closed")
				return
			}
		}
	}()
	// start listening
	err = w.Start(time.Millisecond * 200)
	if err != nil {
		tr := terr.New("watcher.Start", err)
		events.Error("http", "Error starting the watcher", tr)
	}
}

func Stop() {
	w.Close()
}
