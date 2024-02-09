package main

import (
	"errors"
	"golang.org/x/net/context"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/controllers"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"github.com/zelas91/goph-keeper/internal/server/services"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	_ = cancel

	cfg := NewConfig()
	log := logger.New(*cfg.CfgLogger)
	db, err := repository.NewPostgresDB(*cfg.DBurl)
	if err != nil {
		log.Fatalf("db init err : %v", err)
	}

	repo := repository.New(log, db)

	serv := services.New(
		services.WithAuthUseRepository(repo.Auth),
		services.WithCardUseRepository(repo.CreditCard),
		services.WithCredentialUseRepository(repo.Credential),
		services.WithTextUseRepository(repo.TextData),
		services.WithBinaryFileUseRepository(repo.BinaryFile, log, *cfg.BasePathSaveFile),
	)

	handlers := controllers.New(log,
		controllers.WithAuthUseService(serv.Auth),
		controllers.WithCardUseService(serv.CreditCard),
		controllers.WithUserCredentialUseService(serv.Credential),
		controllers.WithTextUseService(serv.TextData),
		controllers.WithBinaryFileUseService(serv.BinaryFile),
	)

	router := chi.NewRouter()
	router.Mount("/", handlers.CreateRoutes())
	server := http.Server{Addr: *cfg.Addr, Handler: router}
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("ListenAndServe %v", err)
		}
	}()

	log.Infof("start server (version - %s, date build - %s)", buildCommit, buildDate)
	<-ctx.Done()

	ctxTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err = server.Shutdown(ctxTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("shutdown server %v", err)
	}

	if err = db.Close(); err != nil {
		log.Error(err)
	}

	log.Info("server stop")

}
