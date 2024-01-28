package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/zelas91/goph-keeper/internal/models"
	"github.com/zelas91/goph-keeper/internal/payload"
	"net/http"
)

type userService interface {
	CreateUser(ctx context.Context, user models.User) error
	CreateToken(ctx context.Context, user models.User) (string, error)
	ParserToken(ctx context.Context, tokenString string) (*models.User, error)
}

type UserHandler struct {
	service userService
}

func NewUserHandler(service userService) *UserHandler {
	return &UserHandler{service: service}
}

func (u *UserHandler) signUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			fmt.Println("invalid content type")
			payload.NewErrorResponse(w, "invalid content type", http.StatusUnsupportedMediaType)
			return
		}
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			fmt.Printf("sigUp json decode err :%v", err)
			payload.NewErrorResponse(w, err.Error(), http.StatusBadRequest)
			return
		}
		u.service.CreateUser(r.Context(), user)
		defer func() {
			if err := r.Body.Close(); err != nil {
				fmt.Printf("sign up body close err :%v", err)
			}
		}()
		w.WriteHeader(http.StatusOK)
	}
}

func (u *UserHandler) InitRoutes() http.Handler {
	router := chi.NewRouter()
	router.Route("/", func(r chi.Router) {
		r.Post("/register", u.signUp())
	})
	return router
}
