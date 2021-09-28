package main // mtcoind

import (
	ol "log"

	"gopkg.in/ini.v1"
)

const __CONF__ = "config"

var config *config_T
type database_T struct {
	user string
	password string
	name string
}

type server_T struct {
	ip string
	port string
	webPort string
	webHost string
	mode bool
}

type config_T struct {
	logLevelKey string
	period int64
	sms int
	server_T
	database_T
	nodes []string
}

func readConf() {
	ol.Println(faciDir + __CONF__)

	cfg, err := ini.Load(faciDir + __CONF__)
	if err != nil {
		ol.Fatal(err)
	}

	config = &config_T{}
	config.logLevelKey = cfg.Section(ini.DEFAULT_SECTION).Key("log_level_key").String()
	config.period, err = cfg.Section(ini.DEFAULT_SECTION).Key("period").Int64()
	if err != nil {
		ol.Fatal(err)
	}
	config.sms, err = cfg.Section(ini.DEFAULT_SECTION).Key("sms").Int()
	if err != nil {
		ol.Fatal(err)
	}
	if config.sms < 0 || config.sms > 1 {
		ol.Fatal("invalid sms config")
	}
	config.ip = cfg.Section("server").Key("ip").String()
	config.port = cfg.Section("server").Key("port").String()
	config.webPort = cfg.Section("server").Key("web_port").String()
	config.webHost = cfg.Section("server").Key("web_host").String()
	config.mode, err = cfg.Section("server").Key("mode").Bool()
	if err != nil {
		ol.Fatal(err)
	}

	config.user = cfg.Section("database").Key("user").String()
	config.password = cfg.Section("database").Key("password").String()
	config.name = cfg.Section("database").Key("name").String()

	config.nodes = cfg.Section("filnodes").Key("nodes").Strings(",")
}
