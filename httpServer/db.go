package httpServer

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/msgs"
	"github.com/synw/terr"
)

var database *gorm.DB

func connect(addr string) (*gorm.DB, *terr.Trace) {
	db, err := gorm.Open("sqlite3", addr)
	if err != nil {
		tr := terr.New("httpServer.db.connect", err)
		return db, tr
	}
	return db, nil
}

func initDb(addr string, hostname string) *terr.Trace {
	msgs.Status("Initializing hits database")
	db, tr := connect(addr)
	if tr != nil {
		tr := terr.Pass("httpServer.db.initDb", tr)
		return tr
	}
	db.AutoMigrate(&types.Hit{})
	database = db
	// run the worker to save the hits to db
	go saveHits("hit_" + hostname)
	return nil
}

func saveToDb(hits []map[string]interface{}) *terr.Trace {
	for _, datapoint := range hits {
		ip := datapoint["Ip"].(string)
		path := datapoint["Path"].(string)
		host := datapoint["Host"].(string)
		ua := datapoint["Ua"].(string)
		lang := datapoint["Lang"].(string)
		length := datapoint["Length"].(string)
		status := int(datapoint["Status"].(float64))
		hit := &types.Hit{
			Ip:     ip,
			Path:   path,
			Host:   host,
			Ua:     ua,
			Lang:   lang,
			Length: length,
			Status: status,
		}
		database.Create(hit)
	}
	return nil
}
