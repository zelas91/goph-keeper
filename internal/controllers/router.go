package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/middleware"
	"github.com/zelas91/goph-keeper/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type Controllers struct {
	*auth
	log *zap.SugaredLogger
}

type services interface {
	userService
}

func NewControllers(log *zap.SugaredLogger, serv services) *Controllers {
	valid := validator.New()
	return &Controllers{
		log:  log,
		auth: newUserHandler(log, serv, valid),
	}
}

func (c *Controllers) InitRoutes() http.Handler {
	router := chi.NewRouter()
	if parser, ok := c.auth.service.(service.ParserService); ok {
		middleware.ValidationAuthorization(c.log, parser)
	}

	router.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
	router.Mount("/api", c.auth.InitRoutes())
	return router
}
