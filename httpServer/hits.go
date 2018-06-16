package httpServer

import (
	"encoding/json"
	"github.com/Tomasen/realip"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/redis"
	"net/http"
)

func processHit(r *http.Request, status int) {
	hit := getRequestData(r, status)
	go saveDatapoint(hit)
}

func saveDatapoint(hit *types.Hit) error {
	// get data
	data, _ := json.Marshal(hit)
	// send to Redis
	key := "hit_" + redis.Hostname
	err := redis.PushKey(key, data)
	if err != nil {
		return err
	}
	return nil
}

func getRequestData(r *http.Request, status int) *types.Hit {
	path := r.URL.Path
	host := r.URL.Host
	header := r.Header
	ua := ""
	lang := ""
	length := ""
	for k, v := range header {
		if k == "User-Agent" {
			ua = v[0]
		} else if k == "Accept-Language" {
			lang = v[0]
		} else if k == "ContentLength" {
			length = v[0]
		}
	}
	ip := realip.FromRequest(r)
	hit := &types.Hit{
		Ip:     ip,
		Path:   path,
		Host:   host,
		Ua:     ua,
		Lang:   lang,
		Length: length,
		Status: status,
	}
	return hit
}
