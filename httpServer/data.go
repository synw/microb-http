package httpServer

import (
	"bufio"
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
	"html/template"
	"os"
	"strings"
)

func getPage(domain string, url string, conn *types.Conn, edit_channel string) (*types.Page, *terr.Trace) {
	// remove eventual trailing slash
	if strings.HasSuffix(url, "/") {
		url = url[:len(url)-1]
	}
	if url == "" {
		url = "/index"
	}
	filepath := basePath + "/static/content/" + getFilepath(url)
	title, content, tr := getContent(filepath)
	if tr != nil {
		tr := tr.Pass("httpServer.data.getPage")
		var p *types.Page
		return p, tr
	}
	page := &types.Page{domain, url, title, template.HTML(content), conn, edit_channel, ""}
	return page, nil
}

func getContent(filepath string) (string, string, *terr.Trace) {
	f, err := os.Open(filepath)
	if err != nil {
		tr := terr.New(err)
		return "", "", tr
	}
	scanner := bufio.NewScanner(f)
	line := 1
	title := ""
	content := ""
	for scanner.Scan() {
		if line == 1 {
			firstLine := scanner.Text()
			prefix := "<!-- Title:"
			if strings.HasPrefix(firstLine, prefix) {
				title = firstLine
				title = strings.Replace(title, prefix, "", -1)
				suffix := " -->"
				title = strings.Replace(title, suffix, "", -1)
			}
		} else {
			content = content + scanner.Text()
		}
		line = line + 1
	}
	return title, content, nil
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
