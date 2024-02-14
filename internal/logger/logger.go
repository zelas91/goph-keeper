package logger

import (
	"encoding/json"
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
)

var (
	once   sync.Once
	logger *zap.SugaredLogger
)

func New(pathCfg string) Logger {
	once.Do(func() {
		file, err := os.ReadFile(pathCfg)
		if err != nil {
			log.Println(err)
			cfg := zap.NewProductionConfig()
			cfg.OutputPaths = []string{"stdout", "noConfigLogger.log"}
			cfg.ErrorOutputPaths = []string{"stdout", "noConfigLogger.log"}
			cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
			l, err := cfg.Build()
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

type Logger interface {
	Infof(template string, args ...interface{})
	Info(args ...interface{})
	Debugf(template string, args ...interface{})
	Debug(args ...interface{})
	Warnf(template string, args ...interface{})
	Warn(args ...interface{})
	Errorf(template string, args ...interface{})
	Error(args ...interface{})
	Fatalf(template string, args ...interface{})
	Fatal(args ...interface{})
	With(args ...interface{}) *zap.SugaredLogger
}
