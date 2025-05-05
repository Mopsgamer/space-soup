package model_http

import "errors"

var ErrInvalidPageRange = errors.New("invalid page range")

type PageInt int

// 1 < page <= max
func (p PageInt) Validate(max int) error {
	page := int(p)
	if page > max || page < 1 {
		return ErrInvalidPageRange
	}
	return nil
}

type TablePage struct {
	TableMode
	Page PageInt `query:"page"`
}
