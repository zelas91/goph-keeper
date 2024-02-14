package repository

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/zelas91/goph-keeper/internal/logger"
	"golang.org/x/net/context"
)

const txKey = "tx"

type transactionManager interface {
	do(ctx context.Context, fn func(ctx context.Context) error) error
	getConn(ctx context.Context) conn
}

type tm struct {
	db  *sqlx.DB
	log logger.Logger
}

func newTm(log logger.Logger, db *sqlx.DB) transactionManager {
	return &tm{
		db:  db,
		log: log,
	}
}
func (t *tm) do(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := t.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err = tx.Rollback(); err != nil && !errors.Is(sql.ErrTxDone, err) {
			t.log.Errorf("Rollback err: %v", err)
			return
		}
	}()

	if err = fn(context.WithValue(ctx, txKey, tx)); err != nil {
		return err
	}
	return tx.Commit()
}

func (t *tm) getConn(ctx context.Context) conn {
	txByCtx := ctx.Value(txKey)
	if txByCtx == nil {
		return t.db
	}
	tx, ok := txByCtx.(*sqlx.Tx)
	if ok {
		return tx
	}
	return t.db
}

type conn interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}
