package logger

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"
	"sync"
)

var (
	once   sync.Once
	logger *zap.SugaredLogger
)

func Shutdown() {
	if err := logger.Sync(); err != nil {
		log.Printf("logger sync %v", err)
	}
}
func New() *zap.SugaredLogger {
	once.Do(func() {
		file, err := os.ReadFile("cfg/config.json")
		if err != nil {
			log.Println(err)
			l, err := zap.NewDevelopment()
			if err != nil {
				log.Fatal(err)
			}
			logger = l.Sugar()
			return
		}

		var cfg zap.Config

		if err := json.Unmarshal(file, &cfg); err != nil {
			log.Fatal(err)
		}

		l, err := cfg.Build()
		if err != nil {
			log.Fatal(err)
		}
		logger = l.Sugar()
	})
	return logger
}
