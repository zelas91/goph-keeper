package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/goph-keeper/internal/controllers"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/repository"
	"github.com/zelas91/goph-keeper/internal/service"
	"net/http"
)

func main() {
	log := logger.New()
	log.Info("start ")
	cfg := NewConfig()
	db, err := repository.NewPostgresDB(*cfg.DBURL)
	if err != nil {
		log.Fatalf("db init err : %v", err)

	}
	repo := repository.NewRepository(log, db)
	serv := service.NewService(repo)
	handlers := controllers.NewControllers(log, serv)

	router := chi.NewRouter()
	router.Mount("/", handlers.InitRoutes())

	http.ListenAndServe(":9095", router)

}
