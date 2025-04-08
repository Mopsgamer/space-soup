//go:build !lite

package main

import (
	"embed"

	server "github.com/Mopsgamer/space-soup/server"
)

//go:embed client/static/** client/templates/**
var embedFS embed.FS

func main() {
	server.Serve(embedFS)
}
