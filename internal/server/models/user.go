package models

type User struct {
	Id       int    `json:"-"`
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
}
