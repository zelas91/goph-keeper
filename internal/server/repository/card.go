package repository

import (
	"errors"
	"fmt"

	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type creditCard struct {
	tm transactionManager
}

func (c creditCard) Create(ctx context.Context, card entities.Card) error {
	query := `insert into cards (number,expired_at,cvv,user_id)
		values (:number,:expired_at,:cvv,:user_id);`
	if _, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, card); err != nil {
		return fmt.Errorf("repo card create err: %v", err)
	}
	return nil
}

func (c creditCard) FindAllByUserID(ctx context.Context, userID int) ([]entities.Card, error) {
	query := `select * from cards where user_id=$1`
	var cards []entities.Card
	if err := c.tm.getConn(ctx).SelectContext(ctx, &cards, query, userID); err != nil {
		return cards, fmt.Errorf("repo: get cards err %v", err)
	}
	return cards, nil
}

func (c creditCard) FindByIDAndUserID(ctx context.Context, cardID, userID int) (entities.Card, error) {
	query := `select * from cards where id=$1 and user_id=$2`
	var card entities.Card
	if err := c.tm.getConn(ctx).GetContext(ctx, &card, query, cardID, userID); err != nil {
		return card, fmt.Errorf("repo: card get id=%d  err: %v", cardID, err)
	}
	return card, nil
}

func (c creditCard) Delete(ctx context.Context, cardID, userID int) error {
	query := `delete from cards where id=$1 and user_id=$2`
	if _, err := c.tm.getConn(ctx).ExecContext(ctx, query, cardID, userID); err != nil {
		return fmt.Errorf("repo card delete err: %v", err)
	}
	return nil
}

func (c creditCard) Update(ctx context.Context, card entities.Card) error {
	err := c.tm.do(ctx, func(ctx context.Context) error {
		query := `select id from cards where id=$1 for update;`
		if _, err := c.tm.getConn(ctx).ExecContext(ctx, query, card.ID); err != nil {
			return fmt.Errorf("repo card update block err :%v", err)
		}
		query = `update cards set
				number=:number,
				cvv=:cvv,
				expired_at=:expired_at
			where
				id=:id and user_id=:user_id and version=:version;`
		result, err := c.tm.getConn(ctx).NamedExecContext(ctx, query, card)
		if err != nil {
			return fmt.Errorf("repo card update err: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("repo card update result err: %v", err)
		}
		if rowsAffected == 0 {
			return errors.New("the versions on the server and client do not match")
		}
		return nil
	})

	return err
}
