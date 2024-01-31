package middleware

import (
	"github.com/zelas91/goph-keeper/internal/logger"
	"github.com/zelas91/goph-keeper/internal/server/payload"
	"github.com/zelas91/goph-keeper/internal/server/types"
	"golang.org/x/net/context"
	"net/http"
)

type tokenParser interface {
	ParserToken(ctx context.Context, tokenString string) (int, error)
}

func AuthorizationHandler(log logger.Logger, parser tokenParser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				log.Errorf("not found jwt (err : %v)", err)
				payload.NewErrorResponse(w, "not found jwt", http.StatusUnauthorized)
				return
			}
			id, err := parser.ParserToken(r.Context(), cookie.Value)
			if err != nil {
				log.Errorf("parse token err : %v", err)
				payload.NewErrorResponse(w, err.Error(), http.StatusUnauthorized)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), types.UserIDKey, id))
			next.ServeHTTP(w, r)
		})
	}
}
