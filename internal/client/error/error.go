package error

import "errors"

var (
	ErrInvalidCommand = errors.New("invalid command usage")
)
