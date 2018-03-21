package conf

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/synw/microb-http/types"
	"github.com/synw/microb/libmicrob/events"
	"github.com/synw/terr"
	"os"
	"path"
	"path/filepath"
)

var Conf *types.Conf

func GetBasePath() string {
	filename, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		msg := "Can not find base directory"
		tr := terr.New("conf.GetConfigPath", err)
		events.New("error", "http", "conf.GetConfigPath", msg, tr)
		return ""
	}
	cp := fmt.Sprintf("%s", path.Dir(filename)) + "/microb-http"
	return cp
}

func GetConf() (*types.Conf, *terr.Trace) {
	// set some defaults for conf
	viper.SetConfigName("http_config")
	cp := GetBasePath()
	viper.AddConfigPath(cp)
	viper.SetDefault("domain", "localhost")
	viper.SetDefault("addr", "localhost:8080")
	viper.SetDefault("centrifugo_addr", "localhost:8001")
	viper.SetDefault("websockets", false)
	viper.SetDefault("datasource", map[string]interface{}{})
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		var conf *types.Conf
		switch err.(type) {
		case viper.ConfigParseError:
			er := errors.New("Error parsing config " + err.Error())
			tr := terr.New("conf.getConf", er)
			return conf, tr
		default:
			err := errors.New("Unable to locate config file at path " + cp)
			tr := terr.New("conf.getConf", err)
			return conf, tr
		}
	}
	domain := viper.GetString("domain")
	addr := viper.GetString("addr")
	caddr := viper.GetString("centrifugo_addr")
	key := viper.GetString("centrifugo_key")
	ws := viper.Get("websockets").(bool)
	ds := viper.Get("datasource").(map[string]interface{})
	datasource := &types.Datasource{}
	for k, v := range ds {
		el := v.(string)
		if k == "name" {
			datasource.Name = el
		} else if k == "path" {
			datasource.Path = el
		} else if k == "user" {
			datasource.User = el
		} else if k == "password" {
			datasource.Pwd = el
		} else if k == "database" {
			datasource.Db = el
		}
	}
	ec := "$edit_" + domain
	viper.SetDefault("edit_channel", ec)
	ech := viper.GetString("edit_channel")
	conf := types.Conf{domain, addr, caddr, key, ws, datasource, ech, true}
	Conf = &conf
	return &conf, nil
}
