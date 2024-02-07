package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/controllers"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"github.com/zelas91/goph-keeper/internal/server/services"
	"net/http"
)

func main() {

	cfg := NewConfig()
	log := logger.New(*cfg.CfgLogger)
	log.Info("start ", cfg.SecretKey)

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

	http.ListenAndServe(*cfg.Addr, router)

}
