package repository

import (
	"errors"
	"github.com/lib/pq"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type auth struct {
	tm transactionManager
}

func (a *auth) Create(ctx context.Context, user entities.User) error {

	if _, err := a.tm.getConn(ctx).ExecContext(ctx,
		"INSERT INTO USERS (login, password) values($1, $2)", user.Login, user.Password); err != nil {
		if errPg := new(pq.PGError); errors.As(err, errPg) {
			return ErrDuplicate
		}
		return err
	}
	return nil
}
func (a *auth) FindUserByLogin(ctx context.Context, user entities.User) (entities.User, error) {
	if err := a.tm.getConn(ctx).GetContext(ctx, &user,
		"SELECT * FROM users WHERE login=$1", user.Login); err != nil {
		return user, err
	}
	return user, nil
}
