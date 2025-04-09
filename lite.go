//go:build lite

package main

import (
	server "github.com/Mopsgamer/space-soup/server"
)

func main() {
	server.Serve(nil)
}
