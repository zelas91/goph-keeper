package entities

import (
	"time"
)

type Card struct {
	ID        int       `db:"id"`
	Version   int       `db:"version"`
	Number    string    `db:"number"`
	ExpiredAt string    `db:"expired_at"`
	Cvv       string    `db:"cvv"`
	UserId    int       `db:"user_id"`
	UpdateAt  time.Time `db:"update_at"`
	CreatedAt time.Time `db:"created_at"`
}
