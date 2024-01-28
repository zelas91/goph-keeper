package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"log"
)

var (
	addr  *string
	dbURL *string
)

func init() {
	addr = flag.String("a", "localhost:8081", "endpoint start server")
	dbURL = flag.String("d", "host=localhost port=5432 user=keeper dbname=goph-keeper password=12345678 sslmode=disable", "url DB")
}

type Config struct {
	Addr  *string `env:"RUN_ADDRESS"`
	DBURL *string `env:"DATABASE_URI"`
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
	if cfg.DBURL == nil {
		cfg.DBURL = dbURL
	}

	flag.Parse()
	return &cfg
}
