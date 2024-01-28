package repository

import (
	"errors"
	"time"
)

var (
	ErrDuplicate = errors.New("login is already taken")
)

type ErrHTTPClient struct {
	Msg        string
	StatusCode int
	RetryTime  time.Duration
}

func (e *ErrHTTPClient) Error() string {
	return e.Msg
}
