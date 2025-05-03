package model_http

type TablePage struct {
	TableExpanded
	Page int `uri:"page"`
}
