package service

import (
	"context"
	"fmt"
	"github.com/zelas91/goph-keeper/internal/models"
)

type UserService struct {
}

func NewUserService() *UserService {
	return &UserService{}
}

func (u UserService) CreateUser(ctx context.Context, user models.User) error {
	fmt.Println(user)
	return nil
}
func (u UserService) CreateToken(ctx context.Context, user models.User) (string, error) {
	return "", nil
}
func (u UserService) ParserToken(ctx context.Context, tokenString string) (*models.User, error) {
	return nil, nil
}
