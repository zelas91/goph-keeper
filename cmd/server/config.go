package main

import (
	"flag"
	"log"

	"github.com/caarlos0/env/v6"
)

var (
	addr             *string
	dbURL            *string
	cfgLogger        *string
	basePathSaveFile *string
	secretKey        *string
	buildCommit      = "N/A"
	buildDate        = "N/A"
)

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
	dbURL = flag.String("d", "host=localhost port=5432 user=keeper dbname=goph-keeper password=12345678 sslmode=disable", "url DB")
	cfgLogger = flag.String("l", "cfg/config.json", "config file logger")
	basePathSaveFile = flag.String("s", "save_file", "dir save file")
	secretKey = flag.String("ek", "", "encrypt secret key")
}

type Config struct {
	Addr             *string `env:"RUN_ADDRESS"`
	DBurl            *string `env:"DATABASE_URI"`
	CfgLogger        *string `env:"CONFIG_LOGGER"`
	BasePathSaveFile *string `env:"BASE_PATH_SAVE"`
	SecretKey        *string `env:"ENCRYPT_SECRET_KEY"`
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
	if cfg.SecretKey == nil {
		cfg.SecretKey = secretKey
	}

	flag.Parse()
	return &cfg
}
