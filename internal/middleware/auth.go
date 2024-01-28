package middleware

import (
	"context"
	"github.com/zelas91/goph-keeper/internal/payload"
	"github.com/zelas91/goph-keeper/internal/service"
	"github.com/zelas91/goph-keeper/internal/types"
	"go.uber.org/zap"
	"net/http"
)

func ValidationAuthorization(log *zap.SugaredLogger, auth service.ParserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				log.Errorf("not found jwt (err : %v)", err)
				payload.NewErrorResponse(w, "not found jwt", http.StatusUnauthorized)
				return
			}
			user, err := auth.ParserToken(r.Context(), cookie.Value)
			if err != nil {
				log.Errorf("parse token err : %v", err)
				payload.NewErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), types.UserIDKey, user.ID))
			next.ServeHTTP(w, r)
		})
	}
}
