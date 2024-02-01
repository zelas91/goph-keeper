package repository

import (
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type creditCard struct {
	tm transactionManager
}

func (c creditCard) Create(ctx context.Context, card entities.Card) error {
	query := `
	INSERT INTO CARDS (number,expired_at,cvv,user_id)
	VALUES (:number,:expired_at,:cvv,:user_id);`
	if _, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, card); err != nil {
		return fmt.Errorf("repo card create err: %v", err)
	}
	return nil
}

func (c creditCard) FindCardsByUserID(ctx context.Context, userID int) ([]entities.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (c creditCard) FindCardByUserID(ctx context.Context, cardID, userID int) (entities.Card, error) {
	//TODO implement me
	panic("implement me")
}

func (c creditCard) Delete(ctx context.Context, cardID, userID int) error {
	//TODO implement me
	panic("implement me")
}

func (c creditCard) Update(ctx context.Context, card entities.Card) error {
	err := c.tm.do(ctx, func(ctx context.Context) error {
		query := `select id from cards where id=$1 for update;`
		if _, err := c.tm.getConn(ctx).ExecContext(ctx, query, card.ID); err != nil {
			return fmt.Errorf("repo card update block err :%v", err)
		}
		query = `UPDATE CARDS SET
				number=:number,
				cvv=:cvv,
				expired_at=:expired_at
			WHERE
				id=:id and user_id=:user_id;`
		if _, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, card); err != nil {
			return fmt.Errorf("repo card update err: %v", err)
		}
		return nil
	})

	return err
}
