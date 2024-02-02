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

type credential struct {
	service credentialService
	valid   *validator.Validate
	log     logger.Logger
}

//go:generate mockgen -package mocks -destination=./mocks/mock_credential_service.go -source=credential.go -package=mock
type credentialService interface {
	Create(ctx context.Context, user models.UserCredentials) error
	Credentials(ctx context.Context) ([]models.UserCredentials, error)
	Credential(ctx context.Context, credentialID int) (models.UserCredentials, error)
	Delete(ctx context.Context, credentialID int) error
	Update(ctx context.Context, credential models.UserCredentials) error
}

func (c *credential) credentials() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		credentials, err := c.service.Credentials(r.Context())
		if err != nil {
			c.log.Errorf("credentials: get credentials err %v", err)
			payload.NewErrorResponse(w, "credentials: get credentials err", http.StatusInternalServerError)
			return
		}
		if err = json.NewEncoder(w).Encode(credentials); err != nil {
			c.log.Errorf("credentials: encode err %v", err)
			payload.NewErrorResponse(w, "credentials: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (c *credential) credential() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("credential: get id from request err: %v", err)
			payload.NewErrorResponse(w, "credential: get id from request err", http.StatusBadRequest)
			return
		}
		credential, err := c.service.Credential(r.Context(), id)
		if err != nil {
			c.log.Errorf("credential: get credential err: %v", err)
			payload.NewErrorResponse(w, "credential: get credential err", http.StatusNotFound)
			return
		}
		if err = json.NewEncoder(w).Encode(credential); err != nil {
			c.log.Errorf("credential: encode err %v", err)
			payload.NewErrorResponse(w, "credential: encode err", http.StatusInternalServerError)
			return
		}
	}
}

func (c *credential) create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "create: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				c.log.Errorf("create: credentials in body close err :%v", err)
			}
		}()

		credential, err := c.credentialFromRequestAndValid(r)
		if err != nil {
			c.log.Errorf("create: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = c.service.Create(r.Context(), credential); err != nil {
			c.log.Errorf("create: credential save err: %v", err)
			payload.NewErrorResponse(w, "create: credential save err", http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func (c *credential) update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "update: body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				c.log.Errorf("update: credentials in body close err :%v", err)
			}
		}()

		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("update: get id credentials err: %v", err)
			payload.NewErrorResponse(w, "update: get id credentials err:", http.StatusBadRequest)
			return
		}

		uc, err := c.credentialFromRequestAndValid(r)
		if err != nil {
			c.log.Errorf("update: decode or validation err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		uc.ID = id

		if err = c.service.Update(r.Context(), uc); err != nil {
			c.log.Errorf("update: credentials save err: %v", err)
			payload.NewErrorResponse(w, "update: credentials save err", http.StatusInternalServerError)
		}
	}
}

func (c *credential) delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := helper.IDFromContext(r.Context())
		if err != nil {
			c.log.Errorf("get id credential err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err = c.service.Delete(r.Context(), id); err != nil {
			c.log.Errorf("delete credential err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (c *credential) credentialFromRequestAndValid(r *http.Request) (credential models.UserCredentials, err error) {
	if err = json.NewDecoder(r.Body).Decode(&credential); err != nil {
		return credential, fmt.Errorf("credential  json decode err:%v", err)
	}

	if err = c.valid.Struct(credential); err != nil {
		return credential, fmt.Errorf("credential validate err: %v", err)
	}
	return
}
func (c *credential) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Get("/", c.credentials())
		r.Get("/{id}", c.credential())
		r.Post("/", c.create())
		r.Put("/{id}", c.update())
		r.Delete("/{id}", c.delete())
	})
	return router
}
