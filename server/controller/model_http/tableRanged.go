package model_http

import (
	"errors"
)

type TableRanged struct {
	Range string `form:"range"`
}

var ErrInvalidIndiceRange = errors.New("invalid indice range")
