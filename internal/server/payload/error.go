package payload

import (
	"encoding/json"
	"net/http"
)

type errorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func NewErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(errorMessage{Message: message, StatusCode: statusCode}); err != nil {
		return
	}
}
