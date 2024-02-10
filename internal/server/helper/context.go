package helper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/go-chi/chi/v5"
	"golang.org/x/net/context"
)

func IDFromContext(ctx context.Context) (id int, err error) {
	idStr := chi.URLParamFromCtx(ctx, "id")
	if idStr == "" {
		return id, errors.New("id not found")
	}
	id, err = strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("convert id=%s, to int err: %w", idStr, err)
	}
	if id == 0 {
		return id, errors.New("id incorrect")
	}
	return
}
