package model_http

import "errors"

var ErrInvalidPageRange = errors.New("invalid page range")

type TablePage struct {
	TableExpanded
	Page int `uri:"page"`
}
