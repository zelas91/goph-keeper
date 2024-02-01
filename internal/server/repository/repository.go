package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zelas91/goph-keeper/internal/logger"
)

type Repository struct {
	Auth       *auth
	CreditCard *creditCard
}

func New(log logger.Logger, db *sqlx.DB) *Repository {
	tm := newTm(log, db)
	return &Repository{
		Auth:       &auth{tm: tm},
		CreditCard: &creditCard{tm: tm},
	}
}

func NewPostgresDB(url string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
