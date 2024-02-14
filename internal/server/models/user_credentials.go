package models

type UserCredentials struct {
	ID       int    `json:"id"`
	Version  int    `json:"version"`
	UserId   int    `json:"-"`
	Login    string `json:"login" validate:"required,min=4"`
	Password string `json:"password" validate:"required,min=8"`
}
