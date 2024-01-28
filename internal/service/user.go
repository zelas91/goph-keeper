package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"github.com/zelas91/goph-keeper/internal/models"
	"github.com/zelas91/goph-keeper/internal/repository/entities"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var secret = []byte("secret_key")

type auth struct {
	repo  userRepo
	cache *cache.Cache
}
type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func WithAuthUseService(up userRepo) func(s *Service) {
	return func(s *Service) {
		s.auth = &auth{repo: up, cache: cache.New(time.Minute*10, time.Minute*10)}
	}
}

//go:generate mockgen -package mocks -destination=./mocks/mock_user.go -source=user.go -package=mock
type userRepo interface {
	CreateUser(ctx context.Context, login, password string) error
	GetUser(ctx context.Context, user models.User) (entities.User, error)
}

func (a *auth) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("generate password hash err")
	}

	return a.repo.CreateUser(ctx, user.Login, string(hashedPassword))
}
func (a *auth) CreateToken(ctx context.Context, authUser models.User) (string, error) {
	user, err := a.repo.GetUser(ctx, authUser)
	if err != nil {
		return "", err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(authUser.Password)); err != nil {
		return "", err
	}
	token, err := generateJwt(user.Login)
	if err != nil {
		return "", err
	}
	a.cache.Set(token, user, cache.DefaultExpiration)
	return token, err
}

func (a *auth) ParserToken(ctx context.Context, tokenString string) (*models.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing jwt")
		}
		return secret, nil
	})
	if err != nil && token == nil {
		return nil, err
	}
	if !token.Valid && !time.Now().Before(claims.ExpiresAt.Time) {
		return nil, errors.New("token not valid")
	}
	val, ok := a.cache.Get(tokenString)

	if ok {
		user := val.(models.User)
		return &user, nil
	}
	user, err := a.repo.GetUser(ctx, models.User{Login: claims.Login})
	if err != nil {
		return nil, err
	}
	a.cache.Set(tokenString, user, cache.DefaultExpiration)
	return &models.User{
		Login:    user.Login,
		Password: user.Password,
		ID:       user.ID,
	}, nil
}

func generateJwt(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		},
		Login: login,
	})
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
