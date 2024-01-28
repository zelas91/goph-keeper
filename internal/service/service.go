package service

import (
	"context"
	"github.com/zelas91/goph-keeper/internal/models"
)

type Service struct {
	*auth
}

type ParserService interface {
	ParserToken(ctx context.Context, tokenString string) (*models.User, error)
}

func New(services ...func(c *Service)) *Service {
	ctl := &Service{}
	for _, serv := range services {
		serv(ctl)
	}
	return ctl
}
