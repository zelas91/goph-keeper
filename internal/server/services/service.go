package services

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Service struct {
	Auth       *auth
	CreditCard *creditCard
	Credential *credential
}

func New(options ...func(s *Service)) *Service {
	sv := &Service{}
	for _, opt := range options {
		opt(sv)
	}
	return sv
}

func WithAuthUseRepository(up userRepo) func(s *Service) {
	return func(s *Service) {
		s.Auth = &auth{repo: up, cache: cache.New(time.Minute*10, time.Minute*10)}
	}
}

func WithCardUseRepository(cr cardRepo) func(s *Service) {
	return func(s *Service) {
		s.CreditCard = &creditCard{repo: cr}
	}
}

func WithCredentialUseRepository(cr credentialRepo) func(s *Service) {
	return func(s *Service) {
		s.Credential = &credential{repo: cr}
	}
}
