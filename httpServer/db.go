package httpServer

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/oschwald/geoip2-golang"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/msgs"
	"github.com/synw/terr"
	"net"
)

var database *gorm.DB

func connectHits(addr string) (*gorm.DB, *terr.Trace) {
	db, err := gorm.Open("sqlite3", addr)
	if err != nil {
		tr := terr.New(err)
		return db, tr
	}
	return db, nil
}

func initHitsDb(addr string, hostname string) *terr.Trace {
	msgs.Status("Initializing hits database")
	db, tr := connectHits(addr)
	if tr != nil {
		tr := tr.Pass()
		return tr
	}
	db.AutoMigrate(&types.Hit{})
	database = db
	// run the worker to save the hits to db
	go saveHits("hit_" + hostname)
	return nil
}

func getGeoIp(ipin string, db *geoip2.Reader) (*geoip2.City, error) {
	if dev == true {
		ipin = "81.2.69.142"
	}
	ip := net.ParseIP(ipin)
	record, err := db.City(ip)
	if err != nil {
		return &geoip2.City{}, err
	}
	return record, nil
}

func getHit(datapoint map[string]interface{}, db *geoip2.Reader) *types.Hit {
	ip := datapoint["Ip"].(string)
	path := datapoint["Path"].(string)
	host := datapoint["Host"].(string)
	ua := datapoint["Ua"].(string)
	lang := datapoint["Lang"].(string)
	length := datapoint["Length"].(string)
	status := int(datapoint["Status"].(float64))
	geo, _ := getGeoIp(ip, db)
	city := geo.City.Names["en"]
	country := geo.Country.Names["en"]
	sub := geo.Subdivisions[0].Names["en"]
	code := geo.Country.IsoCode
	tz := geo.Location.TimeZone
	lat := geo.Location.Latitude
	long := geo.Location.Longitude
	ano := geo.Traits.IsAnonymousProxy
	sat := geo.Traits.IsSatelliteProvider
	hit := &types.Hit{
		Ip:            ip,
		Path:          path,
		Host:          host,
		Ua:            ua,
		Lang:          lang,
		Length:        length,
		Status:        status,
		City:          city,
		Subdivision:   sub,
		CountryName:   country,
		CountryCode:   code,
		Timezone:      tz,
		Latitude:      lat,
		Longitude:     long,
		IsProxy:       ano,
		IsSatProvider: sat,
	}
	return hit
}

func saveToDb(hits []*types.Hit) *terr.Trace {
	tx := database.Begin()
	for _, hit := range hits {
		tx.Create(hit)
	}
	tx.Commit()
	return nil
}
