package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/zelas91/goph-keeper/internal/logger"
)

type Repository struct {
	Auth       *auth
	CreditCard *creditCard
	Credential *credential
}

func New(log logger.Logger, db *sqlx.DB) *Repository {
	manager := newTm(log, db)
	return &Repository{
		Auth:       &auth{tm: manager},
		CreditCard: &creditCard{tm: manager},
		Credential: &credential{tm: manager},
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
