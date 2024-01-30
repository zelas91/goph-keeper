package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var (
	addr      *string
	dbURL     *string
	cfgLogger *string
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	dbURL = flag.String("d", "host=localhost port=5432 user=keeper dbname=goph-keeper password=12345678 sslmode=disable", "url DB")
	cfgLogger = flag.String("l", "cfg/config.json", "config file logger")

}

type Config struct {
	Addr      *string `env:"RUN_ADDRESS"`
	DBurl     *string `env:"DATABASE_URI"`
	CfgLogger *string `env:"CONFIG_LOGGER"`
}

func NewConfig() *Config {
	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatalf("read env error=%v", err)
	}

	if cfg.Addr == nil {
		cfg.Addr = addr
	}
	if cfg.DBurl == nil {
		cfg.DBurl = dbURL
	}
	if cfg.CfgLogger == nil {
		cfg.CfgLogger = cfgLogger
	}

	flag.Parse()
	return &cfg
}
