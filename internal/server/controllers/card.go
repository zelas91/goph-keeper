package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"golang.org/x/net/context"
)

type сreditCard struct {
	service cardService
	valid   *validator.Validate
	log     logger.Logger
}

//go:generate mockgen -package mocks -destination=./mocks/mock_card_service.go -source=card.go -package=mock
type cardService interface {
	Create(ctx context.Context, card models.Card) error
	Cards(ctx context.Context) ([]models.Card, error)
	Card(ctx context.Context, cardID int) (models.Card, error)
	Delete(ctx context.Context, cardID int) error
	Update(ctx context.Context, card models.Card) error
}

func (c *сreditCard) cards() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cards, err := c.service.Cards(r.Context())
		if err != nil {
			c.log.Errorf("cards: get cards err %v", err)
			payload.NewErrorResponse(w, "cards: get cards err", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(cards); err != nil {
			c.log.Errorf("cards: encode err %v", err)
			payload.NewErrorResponse(w, "cards: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (c *сreditCard) card() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("card: get id from request err: %v", err)
			payload.NewErrorResponse(w, "card: get id from request err", http.StatusBadRequest)
			return
		}
		card, err := c.service.Card(r.Context(), id)
		if err != nil {
			c.log.Errorf("card: get card err: %v", err)
			payload.NewErrorResponse(w, "card: get card err", http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(card); err != nil {
			c.log.Errorf("card: encode err %v", err)
			payload.NewErrorResponse(w, "card: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (c *сreditCard) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "create: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				c.log.Errorf("create: card in body close err :%v", err)
			}
		}()

		card, err := c.fromRequestAndValid(r)
		if err != nil {
			c.log.Errorf("create: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = c.service.Create(r.Context(), card); err != nil {
			c.log.Errorf("create: card save err: %v", err)
			payload.NewErrorResponse(w, "create: card save err", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (c *сreditCard) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "update: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				c.log.Errorf("update: card in body close err :%v", err)
			}
		}()

		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("update: get id card err: %v", err)
			payload.NewErrorResponse(w, "update: get id card err:", http.StatusBadRequest)
			return
		}

		card, err := c.fromRequestAndValid(r)
		if err != nil {
			c.log.Errorf("update: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if card.Version == 0 {
			c.log.Error("update: card version == 0")
			payload.NewErrorResponse(w, "update: card version == 0", http.StatusBadRequest)
			return
		}

		card.ID = id

		if err = c.service.Update(r.Context(), card); err != nil {
			c.log.Errorf("update: card save err: %v", err)
			payload.NewErrorResponse(w, "update: card save err", http.StatusInternalServerError)
		}
	}
}

func (c *сreditCard) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("get id card err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = c.service.Delete(r.Context(), id); err != nil {
			c.log.Errorf("delete card err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (c *сreditCard) fromRequestAndValid(r *http.Request) (card models.Card, err error) {
	if err = json.NewDecoder(r.Body).Decode(&card); err != nil {
		return card, fmt.Errorf("card  json decode err:%v", err)
	}

	if err = c.valid.Struct(card); err != nil {
		return card, fmt.Errorf("card validate err: %v", err)
	}
	return
}
func (c *сreditCard) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", c.cards())
		r.Get("/{id}", c.card())
		r.Post("/", c.create())
		r.Put("/{id}", c.update())
		r.Delete("/{id}", c.delete())
	})
	return router
}
