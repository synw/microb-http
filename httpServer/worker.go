package httpServer

import (
	"encoding/json"
	"github.com/oschwald/geoip2-golang"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/microb/libmicrob/redis"
	"github.com/synw/terr"
	"time"
)

func saveHits(key string) {
	for {
		geoDb, err := geoip2.Open("db/GeoLite2-City.mmdb")
		if err != nil {
			tr := terr.New("httpServer.worker.saveHits", err)
			events.Error("http", "Error opening geo ip database", tr, "fatal")
		}
		duration := time.Second * 10
		time.Sleep(duration)
		// get the data from Redis
		keys, err := redis.GetKeys(key)
		if err != nil {
			tr := terr.New("httpServer.worker.saveHits", err)
			tr.Fatal()
		}
		// process the data
		var vals []*types.Hit
		for _, key := range keys {
			var data map[string]interface{}
			err := json.Unmarshal(key.([]byte), &data)
			if err != nil {
				tr := terr.New("httpServer.worker.saveHits", err)
				tr.Fatal()
			}
			hit := getHit(data, geoDb)
			vals = append(vals, hit)
		}
		saveToDb(vals)
	}
}
