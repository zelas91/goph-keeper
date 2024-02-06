package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/helper"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"golang.org/x/net/context"
	"net/http"
)

type textData struct {
	service textDataService
	valid   *validator.Validate
	log     logger.Logger
}

//go:generate mockgen -package mocks -destination=./mocks/mock_text_data_service.go -source=text_data.go -package=mock
type textDataService interface {
	Create(ctx context.Context, text models.TextData) error
	Texts(ctx context.Context) ([]models.TextData, error)
	Text(ctx context.Context, textID int) (models.TextData, error)
	Delete(ctx context.Context, textID int) error
	Update(ctx context.Context, text models.TextData) error
}

func (t *textData) texts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		texts, err := t.service.Texts(r.Context())
		if err != nil {
			t.log.Errorf("text: get texts err %v", err)
			payload.NewErrorResponse(w, "texts: get texts err", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(texts); err != nil {
			t.log.Errorf("texts: encode err %v", err)
			payload.NewErrorResponse(w, "texts: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (t *textData) text() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			t.log.Errorf("text: get id from request err: %v", err)
			payload.NewErrorResponse(w, "text: get id from request err", http.StatusBadRequest)
			return
		}
		text, err := t.service.Text(r.Context(), id)
		if err != nil {
			t.log.Errorf("text: get text err: %v", err)
			payload.NewErrorResponse(w, "text: get text err", http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(text); err != nil {
			t.log.Errorf("text: encode err %v", err)
			payload.NewErrorResponse(w, "text: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (t *textData) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "create: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				t.log.Errorf("create: text in body close err :%v", err)
			}
		}()

		text, err := t.fromRequestAndValid(r)
		if err != nil {
			t.log.Errorf("create: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = t.service.Create(r.Context(), text); err != nil {
			t.log.Errorf("create: text save err: %v", err)
			payload.NewErrorResponse(w, "create: text save err", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (t *textData) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "update: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				t.log.Errorf("update: text in body close err :%v", err)
			}
		}()

		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			t.log.Errorf("update: get id text err: %v", err)
			payload.NewErrorResponse(w, "update: get id text err:", http.StatusBadRequest)
			return
		}

		text, err := t.fromRequestAndValid(r)
		if err != nil {
			t.log.Errorf("update: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if text.Version == 0 {
			t.log.Error("update: text version == 0")
			payload.NewErrorResponse(w, "update: text version == 0", http.StatusBadRequest)
			return
		}

		text.ID = id

		if err = t.service.Update(r.Context(), text); err != nil {
			t.log.Errorf("update: text save err: %v", err)
			payload.NewErrorResponse(w, "update: text save err", http.StatusInternalServerError)
		}
	}
}

func (t *textData) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			t.log.Errorf("get id text err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = t.service.Delete(r.Context(), id); err != nil {
			t.log.Errorf("delete text err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (t *textData) fromRequestAndValid(r *http.Request) (text models.TextData, err error) {
	if err = json.NewDecoder(r.Body).Decode(&text); err != nil {
		return text, fmt.Errorf("text  json decode err:%v", err)
	}

	if err = t.valid.Struct(text); err != nil {
		return text, fmt.Errorf("text validate err: %v", err)
	}
	return
}
func (t *textData) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", t.texts())
		r.Get("/{id}", t.text())
		r.Post("/", t.create())
		r.Put("/{id}", t.update())
		r.Delete("/{id}", t.delete())
	})
	return router
}
