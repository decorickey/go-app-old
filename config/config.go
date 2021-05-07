package config

import (
	"gopkg.in/ini.v1"
	"log"
)

type List struct {
	LogFile   string
	DbName    string
	SqlDriver string
	WebPort   int
	Bucket    string
}

var Config List

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	Config = List{
		LogFile:   cfg.Section("logging").Key("log_file").String(),
		DbName:    cfg.Section("db").Key("name").String(),
		SqlDriver: cfg.Section("db").Key("driver").String(),
		WebPort:   cfg.Section("web").Key("port").MustInt(),
		Bucket:    cfg.Section("aws").Key("bucket").String(),
	}
}
