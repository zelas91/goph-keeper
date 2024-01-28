package service

import (
	"context"
)

type Service struct {
	*auth
}

type ParserService interface {
	ParserToken(ctx context.Context, tokenString string) (int64, error)
}

func New(services ...func(c *Service)) *Service {
	ctl := &Service{}
	for _, serv := range services {
		serv(ctl)
	}
	return ctl
}
