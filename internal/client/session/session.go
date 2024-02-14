package session

import (
	"net/http"
)

type Session struct {
	Jwt  *http.Cookie
	Host string
}

func NewSession(addr string) *Session {
	return &Session{Host: addr}
}

func (s *Session) IsAuth() bool {
	if s.Jwt != nil {
		return s.Jwt.Value != ""
	}
	return false
}

func (s *Session) CleanToken() {
	s.Jwt = nil
}
func (s *Session) GetJwt() *http.Cookie {
	if s.Jwt == nil {
		return &http.Cookie{}
	}
	return s.Jwt
}
