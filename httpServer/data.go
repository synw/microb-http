package httpServer

import (
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
	"io/ioutil"
	"strings"
)

func getPage(domain string, url string, conn *types.Conn, edit_channel string) (*types.Page, *terr.Trace) {
	filepath := "static/content/" + getFilepath(url)
	content, tr := getContent(filepath)
	if tr != nil {
		tr := terr.Pass("httpServer.data.getPage", tr)
		var p *types.Page
		return p, tr
	}
	title := url
	page := &types.Page{domain, url, title, content, conn, edit_channel}
	terr.Debug(page)
	return page, nil
}

func getContent(filepath string) (string, *terr.Trace) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		tr := terr.New("httpServer.data.getContent", err)
		return "", tr
	}
	return string(b), nil
}

func getFilepath(url string) string {
	s := strings.Split(url, "/")
	var file string
	var addr string
	if len(s) == 2 {
		if s[1] == "" {
			file = "index.html"
		} else {
			file = s[1] + ".html"
		}
		addr = file
	} else if len(s) > 2 {
		last := len(s) - 1
		path := ""
		file := ""
		for i, el := range s {
			if i < last {
				var add string
				if path == "" {
					add = el
				} else {
					add = "/" + el
				}
				path = add
			} else {
				file = el + ".html"
			}
		}
		addr = path + "/" + file
	}
	return addr
}
