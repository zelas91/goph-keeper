package session

import (
	"fmt"
	"net/http"
)

const (
	baseURL = "http://"
)

type Session struct {
	Jwt *http.Cookie
	Url string
}

func NewSession(url string) *Session {
	return &Session{Url: fmt.Sprintf("%s%s", baseURL, url)}
}

func (s *Session) IsAuth() bool {
	if s.Jwt != nil {
		return s.Jwt.Value != ""
	}
	return false
}
