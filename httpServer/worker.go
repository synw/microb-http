package httpServer

import (
	"encoding/json"
	"github.com/synw/microb/libmicrob/redis"
	"github.com/synw/terr"
	"time"
)

func saveHits(key string) {
	for {
		duration := time.Second * 2
		time.Sleep(duration)
		// get the data from Redis
		keys, err := redis.GetKeys(key)
		if err != nil {
			tr := terr.New("httpServer.worker.saveHits", err)
			tr.Fatal()
		}
		// process the data
		var vals []map[string]interface{}
		for _, key := range keys {
			var data map[string]interface{}
			err := json.Unmarshal(key.([]byte), &data)
			if err != nil {
				tr := terr.New("httpServer.worker.saveHits", err)
				tr.Fatal()
			}
			vals = append(vals, data)
		}
		saveToDb(vals)
	}
}
