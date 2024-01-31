package models

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
