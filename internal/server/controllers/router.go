package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/middleware"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"net/http"
	"strconv"
)

type Controllers struct {
	*auth
	*сreditCard
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
func WithCardUseService(cs cardService) func(c *Controllers) {
	return func(c *Controllers) {
		c.сreditCard = &сreditCard{service: cs, valid: c.valid, log: c.log}
	}
}

func (c *Controllers) CreateRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
	router.Route("/api", func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Mount("/", c.auth.createRoutes())
		})
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthorizationHandler(c.log, c.auth.service))
			r.Post("/", func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(201)
				u := r.Context().Value(types.UserIDKey).(int)
				w.Write([]byte(strconv.Itoa(u)))
			})
		})
	})

	return router
}
