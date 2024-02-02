package repository

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type credential struct {
	tm transactionManager
}

func (c credential) Create(ctx context.Context, uc entities.UserCredentials) error {
	query := `insert into user_credentials (login,password,user_id)
		values (:login,:password,:user_id);`
	if _, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, uc); err != nil {
		return fmt.Errorf("repo credentials create err: %v", err)
	}
	return nil
}

func (c credential) FindCredentialsByUserID(ctx context.Context, userID int) ([]entities.UserCredentials, error) {
	query := `select * from user_credentials where user_id=$1`
	var ucs []entities.UserCredentials
	if err := c.tm.getConn(ctx).SelectContext(ctx, &ucs, query, userID); err != nil {
		return ucs, fmt.Errorf("repo: get credentials err %v", err)
	}
	return ucs, nil
}

func (c credential) FindCredentialByUserID(ctx context.Context, ucID, userID int) (entities.UserCredentials, error) {
	query := `select * from user_credentials where id=$1 and user_id=$2`
	var uc entities.UserCredentials
	if err := c.tm.getConn(ctx).GetContext(ctx, &uc, query, ucID, userID); err != nil {
		return uc, fmt.Errorf("repo: credentials get id=%d  err: %v", ucID, err)
	}
	return uc, nil
}

func (c credential) Delete(ctx context.Context, ucID, userID int) error {
	query := `delete from user_credentials where id=$1 and user_id=$2`
	if _, err := c.tm.getConn(ctx).ExecContext(ctx, query, ucID, userID); err != nil {
		return fmt.Errorf("repo credentials delete err: %v", err)
	}
	return nil
}

func (c credential) Update(ctx context.Context, uc entities.UserCredentials) error {
	err := c.tm.do(ctx, func(ctx context.Context) error {
		query := `select id from user_credentials where id=$1 for update;`
		if _, err := c.tm.getConn(ctx).ExecContext(ctx, query, uc.ID); err != nil {
			return fmt.Errorf("repo credentials update block err :%v", err)
		}
		query = `update user_credentials set
				login=:login,
				password=:password
			where
				id=:id and user_id=:user_id;`
		if _, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, uc); err != nil {
			return fmt.Errorf("repo credentials update err: %v", err)
		}
		return nil
	})

	return err
}
