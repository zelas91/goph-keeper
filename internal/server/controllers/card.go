package controllers

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"golang.org/x/net/context"
	"net/http"
)

type сreditCard struct {
	service cardService
	valid   *validator.Validate
	log     logger.Logger
}

type cardService interface {
	Create(ctx context.Context, card models.Card) error
	GetAll(ctx context.Context) ([]models.Card, error)
	Get(ctx context.Context, cardID int) (models.Card, error)
	Delete(ctx context.Context, cardID int) error
	Update(ctx context.Context, card models.Card) error
}

func (c *сreditCard) getAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *сreditCard) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *сreditCard) add() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *сreditCard) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (c *сreditCard) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
func (c *сreditCard) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/card", func(r chi.Router) {
		r.Get("/", c.getAll())
		r.Get("/{id}", c.get())
		r.Post("/", c.add())
		r.Put("/{id}", c.update())
		r.Delete("/{id}", c.delete())
	})
	return router
}
