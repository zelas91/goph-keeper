package repository

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Repository struct {
	*auth
}

func New(log *zap.SugaredLogger, db *sqlx.DB) *Repository {
	tm := newTm(log, db)
	return &Repository{auth: newAuth(tm)}
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
