package models

type User struct {
	ID       int    `json:"-"`
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}
