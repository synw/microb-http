package httpServer

import (
	"encoding/json"
	//"fmt"
	"github.com/oschwald/geoip2-golang"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/events"
	"github.com/synw/microb/redis"
	"github.com/synw/terr"
	"time"
)

func saveHits(key string) {
	for {
		//fmt.Println("Sleeping")

		duration := time.Second * 10
		time.Sleep(duration)
		// get the geo data
		//fmt.Println("Opening geo db")

		geoDb, err := geoip2.Open("db/geo/GeoLite2-City.mmdb")
		if err != nil {
			tr := terr.New(err)
			events.Error("http", "Error opening geo ip database", tr)
		}

		// get the data from Redis
		//fmt.Println("Getting keys")

		keys, err := redis.GetKeys(key)
		if err != nil {
			tr := terr.New(err)
			tr.Fatal("Can not get keys in Redis")
		}
		// process the data
		var vals []*types.Hit
		//fmt.Println("************************************** Processing", len(keys), "key")

		for _, key := range keys {
			var data map[string]interface{}
			err := json.Unmarshal(key.([]byte), &data)
			if err != nil {
				tr := terr.New(err)
				tr.Fatal("Can not unmarshal json")
			}
			hit := getHit(data, geoDb)
			vals = append(vals, hit)
		}
		//fmt.Println("Saving to db")

		saveToDb(vals)
		//fmt.Println("Data saved")
	}
}
