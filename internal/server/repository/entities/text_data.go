package entities

import "time"

type TextData struct {
	ID        int       `db:"id"`
	Version   int       `db:"version"`
	UserId    int       `db:"user_id"`
	UpdateAt  time.Time `db:"update_at"`
	CreatedAt time.Time `db:"created_at"`
	Text      string    `db:"large_text"`
}
