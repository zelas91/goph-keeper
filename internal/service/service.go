package service

import (
	"context"
	"github.com/zelas91/goph-keeper/internal/models"
)

type Service struct {
	*auth
}

type repository interface {
	userRepo
}
type ParserService interface {
	ParserToken(ctx context.Context, tokenString string) (*models.User, error)
}

func NewService(repo repository) *Service {
	return &Service{auth: newUserService(repo)}
}
