package watcher

import (
	"fmt"
	wa "github.com/radovskyb/watcher"
	"github.com/synw/centcom"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"time"
)

var w = wa.New()

func Start(path string, cli *centcom.Cli, channel string, verbosity int, dev bool) {
	iserr := false
	if dev == true {
		_ = w.AddRecursive("../microb-http/templates")
		_ = w.AddRecursive("../microb-http/static/js")
		_ = w.AddRecursive("../microb-http/static/css")
		_ = w.AddRecursive("../microb-http/static/content")
	}
	err := w.AddRecursive("templates")
	if err != nil {
		tr := terr.New("wa.Start", err)
		tr.Printf("start watcher")
		events.Terr("http", "watcher.Start", "Error finding templates", tr)
		iserr = true
	}
	err = w.AddRecursive("static/js")
	if err != nil {
		tr := terr.New("wa.Start", err)
		tr.Printf("start watcher")
		events.Terr("http", "watcher.Start", "Error finding static/js", tr)
		iserr = true
	}
	err = w.AddRecursive("static/css")
	if err != nil {
		tr := terr.New("wa.Start", err)
		tr.Printf("start watcher")
		events.Terr("http", "watcher.Start", "Error finding static/css", tr)
		iserr = true
	}
	err = w.AddRecursive("static/content")
	if err != nil {
		tr := terr.New("wa.Start", err)
		tr.Printf("start watcher")
		events.Terr("http", "watcher.Start", "Error finding static/content", tr)
		iserr = true
	}
	w.FilterOps(wa.Write, wa.Create, wa.Move, wa.Remove, wa.Rename)
	if verbosity > 1 && iserr == false {
		fmt.Println("Watching files :")
		for path, f := range w.WatchedFiles() {
			fmt.Printf("%s %s\n", f.Name(), path)
		}
	}
	// lauch listener
	go func() {
		for {
			select {
			case e := <-w.Event:
				if verbosity > 2 {
					fmt.Println("Change detected in", e.Path, "reloading")
				}
				handle(cli, channel)
			case err := <-w.Error:
				fmt.Println("Watcher error", err)
			case <-w.Closed:
				fmt.Println("Watcher closed")
				return
			}
		}
	}()
	// start listening
	err = w.Start(time.Millisecond * 200)
	if err != nil {
		tr := terr.New("watcher.Start", err)
		tr.Printf("start watcher")
		events.Terr("http", "watcher.Start", "Error starting the watcher", tr)
	}
}

func Stop() {
	w.Close()
}
