package entities

import (
	"time"
)

type Card struct {
	ID        int       `db:"id"`
	Version   int       `db:"version"`
	Number    []byte    `db:"number"`
	ExpiredAt []byte    `db:"expired_at"`
	Cvv       []byte    `db:"cvv"`
	UserId    int       `db:"user_id"`
	UpdateAt  time.Time `db:"update_at"`
	CreatedAt time.Time `db:"created_at"`
}
