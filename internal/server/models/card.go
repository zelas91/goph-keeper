package models

type Card struct {
	ID        int    `json:"id"`
	Version   int    `json:"version "`
	Number    string `json:"number" validate:"required,credit_card"`
	ExpiredAt string `json:"expired_at" validate:"required,expired_credit_card"`
	Cvv       string `json:"cvv" validate:"required,len=3"`
	UserId    int    `json:"-"`
}
