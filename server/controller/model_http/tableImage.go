package model_http

type TableImage struct {
	IsDownload  bool   `query:"download"`
	Range       string `form:"range"`
	Description string `form:"description"`
}
