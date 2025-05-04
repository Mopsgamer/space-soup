package model_http

type TableImage struct {
	TableRanged
	IsDownload  bool   `query:"download"`
	Description string `form:"description"`
}
