package model_http

import "github.com/Mopsgamer/space-soup/server/soup"

type OrbitInputFile struct {
	TablePage
	TableRanged
	FileType soup.FileType `form:"file-type"`
}
