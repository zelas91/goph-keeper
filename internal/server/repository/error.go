package repository

import (
	"errors"
)

var (
	ErrDuplicate = errors.New("login is already taken")
)
