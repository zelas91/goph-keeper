package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/middleware"
	"go.uber.org/zap"
	"net/http"
)

type Controllers struct {
	*auth
	log   *zap.SugaredLogger
	valid *validator.Validate
}

func New(log *zap.SugaredLogger, options ...func(c *Controllers)) *Controllers {
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

func (c *Controllers) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
	router.Route("/api", func(r chi.Router) {
		r.Mount("/", c.auth.initRoutes())
	})
	router.Route("/test", func(r chi.Router) {
		r.Use(middleware.AuthorizationHandler(c.log, c.auth.service))
		r.Post("/", func(writer http.ResponseWriter, request *http.Request) {
			writer.Write([]byte("asdasdasdasdasdas"))
			writer.WriteHeader(http.StatusCreated)
		})
	})
	return router
}
