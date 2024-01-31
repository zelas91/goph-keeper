package entities

import "time"

type User struct {
	Id        int       `db:"id"`
	Login     string    `db:"login"`
	Password  string    `db:"password"`
	CreatedAt time.Time `db:"created_at"`
}
