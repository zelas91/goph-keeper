package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/middleware"
	"net/http"
)

type Controllers struct {
	*auth
	log   logger.Logger
	valid *validator.Validate
}

func New(log logger.Logger, options ...func(c *Controllers)) *Controllers {
	ctl := &Controllers{
		log:   log,
		valid: validator.New(),
	}
	for _, opt := range options {
		opt(ctl)
	}
	return ctl
}

func WithAuthUseService(us userService) func(c *Controllers) {
	return func(c *Controllers) {
		c.auth = &auth{service: us, valid: c.valid, log: c.log}
	}
}

func (c *Controllers) CreateRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
	router.Route("/api", func(r chi.Router) {
		r.Mount("/", c.auth.createRoutes())
	})
	router.Route("/test", func(r chi.Router) {
		r.Use(middleware.AuthorizationHandler(c.log, c.auth.service))
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.WriteHeader(http.StatusCreated)
			writer.Write([]byte("asdasdasdasdasdas"))
		})
	})
	return router
}
