package repository

import (
	"errors"
	"fmt"
	"github.com/zelas91/goph-keeper/internal/server/repository/entities"
	"golang.org/x/net/context"
)

type textData struct {
	tm transactionManager
}

func (t textData) Create(ctx context.Context, uc entities.TextData) error {
	query := `insert into text_data (large_text,user_id)
		values (large_text=:large_text,user_id=:user_id);`
	if _, err := t.tm.getConn(ctx).NamedExecContext(ctx, query, uc); err != nil {
		return fmt.Errorf("repo text create err: %v", err)
	}
	return nil
}

func (t textData) FindAllByUserID(ctx context.Context, userID int) ([]entities.TextData, error) {
	query := `select * from text_data where user_id=$1`
	var texts []entities.TextData
	if err := t.tm.getConn(ctx).SelectContext(ctx, &texts, query, userID); err != nil {
		return texts, fmt.Errorf("repo: get texts err %v", err)
	}
	return texts, nil
}

func (t textData) FindByIDAndUserID(ctx context.Context, textID, userID int) (entities.TextData, error) {
	query := `select * from text_data where id=$1 and user_id=$2`
	var text entities.TextData
	if err := t.tm.getConn(ctx).GetContext(ctx, &text, query, textID, userID); err != nil {
		return text, fmt.Errorf("repo: text get id=%d  err: %v", textID, err)
	}
	return text, nil
}

func (t textData) Delete(ctx context.Context, textID, userID int) error {
	query := `delete from text_data where id=$1 and user_id=$2`
	if _, err := t.tm.getConn(ctx).ExecContext(ctx, query, textID, userID); err != nil {
		return fmt.Errorf("repo text delete err: %v", err)
	}
	return nil
}

func (t textData) Update(ctx context.Context, text entities.TextData) error {
	err := t.tm.do(ctx, func(ctx context.Context) error {
		query := `select id from text_data where id=$1 for update;`
		if _, err := t.tm.getConn(ctx).ExecContext(ctx, query, text.ID); err != nil {
			return fmt.Errorf("repo text update block err :%v", err)
		}
		query = `update text_data set
				large_text=:large_text
			where
				id=:id and user_id=:user_id and version=:version;`
		result, err := t.tm.getConn(ctx).NamedExecContext(ctx, query, text)
		if err != nil {
			return fmt.Errorf("repo text update err: %v", err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("repo text update result err: %v", err)
		}
		if rowsAffected == 0 {
			return errors.New("the versions on the server and client do not match")
		}
		return nil
	})

	return err
}
