package models

type TextData struct {
	ID      int    `json:"id"`
	Version int    `json:"version"`
	UserId  int    `json:"-"`
	Text    string `json:"text"`
}
