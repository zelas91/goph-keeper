package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/models"
	"github.com/zelas91/goph-keeper/internal/payload"
	"github.com/zelas91/goph-keeper/internal/repository"
	"go.uber.org/zap"
	"net/http"
)

type auth struct {
	service userService
	valid   *validator.Validate
	log     *zap.SugaredLogger
}
type userService interface {
	CreateUser(ctx context.Context, user models.User) error
	CreateToken(ctx context.Context, user models.User) (string, error)
}

func newUserHandler(log *zap.SugaredLogger, service userService, valid *validator.Validate) *auth {
	return &auth{service: service, valid: valid, log: log}
}

func (a *auth) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			a.log.Errorf("sigUp json decode err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				a.log.Errorf("sign up body close err :%v", err)
			}
		}()

		if err := a.valid.Struct(user); err != nil {
			a.log.Errorf("sigUp json validate err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := a.service.CreateUser(r.Context(), user); err != nil {
			if errors.Is(err, repository.ErrDuplicate) {
				a.log.Errorf("sigUp user duplicate err :%v", err)
				payload.NewErrorResponse(w, err.Error(), http.StatusConflict)
				return
			}
			a.log.Errorf("sigUp create user err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusInternalServerError)
			return
		}
		token, err := a.service.CreateToken(r.Context(), user)
		if err != nil {
			a.log.Errorf("sigUp create token err :%v", err)
			payload.NewErrorResponse(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		cookies := http.Cookie{
			Path:  "/",
			Name:  "jwt",
			Value: token,
		}
		http.SetCookie(w, &cookies)

		w.WriteHeader(http.StatusOK)
	}
}

func (a *auth) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			a.log.Errorf("signIn json decode err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := a.valid.Struct(user); err != nil {
			a.log.Errorf("signIn json validate err:%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		token, err := a.service.CreateToken(r.Context(), user)
		if err != nil {
			a.log.Errorf("signIn create token err:%v", err)
			payload.NewErrorResponse(w, "invalid login or password", http.StatusUnauthorized)
			return
		}
		cookies := http.Cookie{
			Path:  "/",
			Name:  "jwt",
			Value: token,
		}
		http.SetCookie(w, &cookies)
		w.WriteHeader(http.StatusOK)
	}
}
func (a *auth) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/signup", a.signUp())
		r.Post("/signin", a.signIn())
	})
	return router
}
