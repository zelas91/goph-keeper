package repository

import (
	"errors"
	"github.com/lib/pq"
	"github.com/zelas91/goph-keeper/internal/models"
	"github.com/zelas91/goph-keeper/internal/repository/entities"
	"golang.org/x/net/context"
)

type auth struct {
	tm transactionManager
}

func newAuth(tm transactionManager) *auth {
	return &auth{tm: tm}
}

func (a *auth) CreateUser(ctx context.Context, user entities.User) error {

	if _, err := a.tm.getConn(ctx).ExecContext(ctx,
		"INSERT INTO USERS (login, password) values($1, $2)", user.Login, user.Password); err != nil {
		if errPg := new(pq.PGError); errors.As(err, errPg) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}
func (a *auth) GetUser(ctx context.Context, authUser models.User) (entities.User, error) {
	var user entities.User
	if err := a.tm.getConn(ctx).GetContext(ctx, &user,
		"SELECT * FROM users WHERE login=$1", authUser.Login); err != nil {
		return user, err
	}
	return user, nil
}
