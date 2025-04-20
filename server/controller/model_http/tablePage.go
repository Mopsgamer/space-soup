package model_http

type TablePage struct {
	Page     int  `uri:"page"`
	Expanded bool `cookie:"expanded"`
}
