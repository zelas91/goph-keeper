package entities

import "time"

type BinaryFile struct {
	ID        int       `db:"id"`
	UserId    int       `db:"user_id"`
	Path      string    `db:"path"`
	FileName  string    `db:"file_name"`
	CreatedAt time.Time `db:"created_at"`
}
