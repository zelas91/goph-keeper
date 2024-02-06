package controllers

import (
	"github.com/go-chi/chi/v5"
	middleware2 "github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/middleware"
	"github.com/zelas91/goph-keeper/internal/utils/validation"
	"net/http"
)

type Controllers struct {
	auth       *auth
	card       *сreditCard
	credential *credential
	textData   *textData
	binary     *binaryFile
	log        logger.Logger
	valid      *validator.Validate
}

func New(log logger.Logger, options ...func(c *Controllers)) *Controllers {
	ctl := &Controllers{
		log:   log,
		valid: validation.NewValidator(log),
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
		c.card = &сreditCard{service: cs, valid: c.valid, log: c.log}
	}
}

func WithUserCredentialUseService(cs credentialService) func(c *Controllers) {
	return func(c *Controllers) {
		c.credential = &credential{service: cs, valid: c.valid, log: c.log}
	}
}

func WithTextUseService(td textDataService) func(c *Controllers) {
	return func(c *Controllers) {
		c.textData = &textData{service: td, valid: c.valid, log: c.log}
	}
}

func WithBinaryFileUseService(bs binaryFileService) func(c *Controllers) {
	return func(c *Controllers) {
		c.binary = &binaryFile{service: bs, valid: c.valid, log: c.log}
	}
}
func (c *Controllers) CreateRoutes() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.ContentTypeJSON(c.log), middleware2.Recoverer)
	router.Route("/api", func(r chi.Router) {
		r.Mount("/", c.auth.createRoutes())
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthorizationHandler(c.log, c.auth.service))
			r.Group(func(r chi.Router) {
				r.Mount("/card", c.card.createRoutes())
				r.Mount("/credential", c.credential.createRoutes())
				r.Mount("/text", c.textData.createRoutes())
				r.Mount("/file", c.binary.createRoutes())
			})
		})
	})
	return router
}
