package service

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/patrickmn/go-cache"
	"github.com/zelas91/goph-keeper/internal/server/models"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"github.com/zelas91/goph-keeper/internal/utils"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
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

//go:generate mockgen -package mocks -destination=./mocks/mock_user_repo.go -source=user.go -package=mock
type userRepo interface {
	CreateUser(ctx context.Context, user entities.User) error
	GetUser(ctx context.Context, user entities.User) (entities.User, error)
}

func (a *auth) CreateUser(ctx context.Context, user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("generate password hash err")
	}

	return a.repo.CreateUser(ctx, entities.User{Login: user.Login, Password: string(hashedPassword)})
}
func (a *auth) CreateToken(ctx context.Context, authUser models.User) (string, error) {
	user, err := a.repo.GetUser(ctx, utils.ModelUserInEntitiesUser(authUser))
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

func (a *auth) ParserToken(ctx context.Context, tokenString string) (int64, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("error parsing jwt")
		}
		return secret, nil
	})
	if err != nil && token == nil {
		return 0, err
	}

	if !token.Valid && !time.Now().Before(claims.ExpiresAt.Time) {
		return 0, errors.New("token not valid")
	}

	val, ok := a.cache.Get(tokenString)
	if ok {
		user := val.(entities.User)
		return user.ID, nil
	}

	user, err := a.repo.GetUser(ctx, entities.User{Login: claims.Login})
	if err != nil {
		return 0, err
	}

	a.cache.Set(tokenString, user, cache.DefaultExpiration)

	return user.ID, nil
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
