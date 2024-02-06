package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var (
	addr             *string
	dbURL            *string
	cfgLogger        *string
	basePathSaveFile *string
	secretKey        = ""
	buildCommit      = "N/A"
	buildDate        = "N/A"
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	dbURL = flag.String("d", "host=localhost port=5432 user=keeper dbname=goph-keeper password=12345678 sslmode=disable", "url DB")
	cfgLogger = flag.String("l", "cfg/config.json", "config file logger")
	basePathSaveFile = flag.String("s", "save_file", "dir save file")
}

type Config struct {
	Addr             *string `env:"RUN_ADDRESS"`
	DBurl            *string `env:"DATABASE_URI"`
	CfgLogger        *string `env:"CONFIG_LOGGER"`
	BasePathSaveFile *string `env:"BASE_PATH_SAVE"`
	Version          string
	BuildData        string
	SecretKey        string
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

	if cfg.BasePathSaveFile == nil {
		cfg.BasePathSaveFile = basePathSaveFile
	}

	cfg.BuildData = buildDate
	cfg.SecretKey = secretKey
	cfg.Version = buildCommit

	flag.Parse()
	return &cfg
}
