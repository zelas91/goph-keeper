package service

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type Service struct {
	*auth
}

func New(options ...func(c *Service)) *Service {
	ctl := &Service{}
	for _, opt := range options {
		opt(ctl)
	}
	return ctl
}

func WithAuthUseRepository(up userRepo) func(s *Service) {
	return func(s *Service) {
		s.auth = &auth{repo: up, cache: cache.New(time.Minute*10, time.Minute*10)}
	}
}
