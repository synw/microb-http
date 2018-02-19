package conf

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"github.com/synw/microb-http/types"
	"github.com/synw/terr"
	"path"
	"runtime"
)

var Conf *types.Conf

func GetConf(dev bool) (*types.Conf, *terr.Trace) {
	name := "normal"
	if dev {
		name = "dev"
	}
	return getConf(name)
}

func getConfigPath() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("No caller information")
	}
	cp := fmt.Sprintf("%s", path.Dir(filename))
	return cp
}

func getConf(name string) (*types.Conf, *terr.Trace) {
	// set some defaults for conf
	dev := false
	if name == "dev" {
		viper.SetConfigName("http_dev_config")
		dev = true
	} else {
		viper.SetConfigName("http_config")
	}
	cp := getConfigPath()
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
			err := errors.New("Unable to locate config file")
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
	conf := types.Conf{domain, addr, caddr, key, ws, datasource, ech, dev}
	Conf = &conf
	return &conf, nil
}
