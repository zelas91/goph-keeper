package models

import "time"

type Card struct {
	ID        int       `json:"id"`
	Version   int       `json:"version"`
	Number    string    `json:"number" validate:"required,credit_card"`
	ExpiredAt string    `json:"expired_at" validate:"required"`
	Cvv       string    `json:"cvv" validate:"required"`
	UserId    int       `json:"userid"`
	UpdateAt  time.Time `json:"update_at"`
	CreatedAt time.Time `json:"created_at"`
}
