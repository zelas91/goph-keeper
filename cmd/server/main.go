package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/controllers"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"github.com/zelas91/goph-keeper/internal/server/service"
	"net/http"
)

func main() {
	cfg := NewConfig()
	log := logger.New(*cfg.CfgLogger)
	log.Info("start ")

	db, err := repository.NewPostgresDB(*cfg.DBurl)
	if err != nil {
		log.Fatalf("db init err : %v", err)
	}

	repo := repository.New(log, db)

	serv := service.New(
		service.WithAuthUseRepository(repo),
	)

	handlers := controllers.New(log,
		controllers.WithAuthUseService(serv),
	)

	router := chi.NewRouter()
	router.Mount("/", handlers.CreateRoutes())

	http.ListenAndServe(*cfg.Addr, router)

}
