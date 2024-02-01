package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"github.com/zelas91/goph-keeper/internal/server/repository"
	"golang.org/x/net/context"
	"net/http"
)

type auth struct {
	service userService
	valid   *validator.Validate
	log     logger.Logger
}

//go:generate mockgen -package mocks -destination=./mocks/mock_user_service.go -source=user.go -package=mock
type userService interface {
	ParserToken(ctx context.Context, tokenString string) (int, error)
	CreateUser(ctx context.Context, user models.User) error
	CreateToken(ctx context.Context, user models.User) (string, error)
}

func (a *auth) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				a.log.Errorf("sign in body close err :%v", err)
			}
		}()

		user, err := a.userFromRequestAndValid(r)
		if err != nil {
			a.log.Errorf("sign up get use err: %v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err = a.service.CreateUser(r.Context(), user); err != nil {
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
	}
}

func (a *auth) signIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			payload.NewErrorResponse(w, "body is empty", http.StatusBadRequest)
			return
		}

		defer func() {
			if err := r.Body.Close(); err != nil {
				a.log.Errorf("sign in body close err :%v", err)
			}
		}()

		user, err := a.userFromRequestAndValid(r)
		if err != nil {
			a.log.Errorf("sign in get use err: %v", err)
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
	}
}

func (a *auth) userFromRequestAndValid(r *http.Request) (user models.User, err error) {
	if err = json.NewDecoder(r.Body).Decode(&user); err != nil {
		return user, fmt.Errorf("signIn json decode err:%w", err)
	}
	if err = a.valid.Struct(user); err != nil {
		return user, fmt.Errorf("signIn json validate err:%w", err)

	}
	return
}

func (a *auth) createRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/signup", a.signUp())
		r.Post("/signin", a.signIn())
	})
	return router
}
