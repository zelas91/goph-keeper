package models

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Card     string `json:"card" validate:"required,credit_card"`
}
