package middleware

import (
	"github.com/zelas91/goph-keeper/internal/payload"
	"go.uber.org/zap"
	"net/http"
)

const (
	content     = "Content-Type"
	contentJSON = "application/json"
)

func ContentTypeJSON(log *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Header.Get(content) != contentJSON {
				log.Error("invalid content type")
				payload.NewErrorResponse(w, "invalid content type", http.StatusUnsupportedMediaType)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

}
