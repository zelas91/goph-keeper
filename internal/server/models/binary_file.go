package models

type BinaryFile struct {
	ID       int    `json:"id"`
	UserId   int    `json:"-"`
	Path     string `json:"-"`
	FileName string `json:"file_name" validate:"required"`
	Size     int    `json:"size" validate:"required"`
}

type AnswerBinaryFile struct {
	Confirm bool `json:"confirm"`
}
