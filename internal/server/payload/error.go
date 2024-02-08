package payload

import (
	"encoding/json"
	"net/http"
)

type ErrorMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func NewErrorResponse(w http.ResponseWriter, message string, statusCode int) {

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ErrorMessage{Message: message, StatusCode: statusCode}); err != nil {
		return
	}
}
