package model_http

import (
	"errors"
	"path/filepath"
	"strings"
)

var (
	ErrUnsupportedFileType = errors.New("unsupported file type")
)

type FileType string

func (ft FileType) DecideFileName(filename string) FileType {
	switch ft {
	case "":
		fallthrough
	case "auto":
	default:
		return ft
	}

	ext := FileType(strings.ToLower(filepath.Ext(filename)))
	return FileType(ext[1:])
}
