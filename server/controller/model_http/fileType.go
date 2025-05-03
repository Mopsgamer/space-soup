package model_http

import "errors"

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type FileType string

func (ft FileType) ShouldDecide() bool {
	switch string(ft) {
	case "":
		return true
	case "auto":
		return true
	default:
		return false
	}
}
