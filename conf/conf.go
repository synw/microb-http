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
	if name == "dev" {
		viper.SetConfigName("dev_config")
	} else {
		viper.SetConfigName("config")
	}
	cp := getConfigPath()
	viper.AddConfigPath(cp)
	viper.SetDefault("domain", "localhost")
	viper.SetDefault("addr", "localhost:8080")
	viper.SetDefault("centrifugo_addr", "localhost:8001")
	// get the actual conf
	err := viper.ReadInConfig()
	if err != nil {
		var conf *types.Conf
		switch err.(type) {
		case viper.ConfigParseError:
			tr := terr.New("conf.getConf", err)
			return conf, tr
		default:
			err := errors.New("Unable to locate config file")
			tr := terr.New("conf.getConf", err)
			return conf, tr
		}
	}
	domain := viper.GetString("domain")
	url := viper.GetString("addr")
	addr := viper.GetString("centrifugo_addr")
	key := viper.GetString("centrifugo_key")
	conf := types.Conf{domain, url, addr, key}
	return &conf, nil
}
