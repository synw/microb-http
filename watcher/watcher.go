package watcher

import (
	"fmt"
	wa "github.com/radovskyb/watcher"
	"github.com/synw/centcom"
	"github.com/synw/terr"
	"time"
)

var w = wa.New()

func Start(path string, cli *centcom.Cli, channel string, verbosity int, dev bool) {
	if dev == true {
		err := w.AddRecursive("../microb-http/templates")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("../microb-http/static/js")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("../microb-http/static/css")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("../microb-http/static/content")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
	} else {
		err := w.AddRecursive("templates")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("static/js")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("static/css")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
		err = w.AddRecursive("static/content")
		if err != nil {
			tr := terr.New("wa.Start", err)
			tr.Fatal("start watcher")
		}
	}
	w.FilterOps(wa.Write, wa.Create, wa.Move, wa.Remove, wa.Rename)
	if verbosity > 1 {
		fmt.Println("Watching files :")
		for _, f := range w.WatchedFiles() {
			fmt.Printf("%s\n", f.Name())
		}
	}
	// lauch listener
	go func() {
		for {
			select {
			case _ = <-w.Event:
				//fmt.Println("EVENT", event.Path)
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
	err := w.Start(time.Millisecond * 100)
	if err != nil {
		tr := terr.New("watcher.Start", err)
		tr.Fatal("start watcher")
	}
}

func Stop() {
	w.Close()
}
