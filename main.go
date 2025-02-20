package main

import (
	"github.com/Mopsgamer/space-soup/server"
	"github.com/Mopsgamer/space-soup/server/environment"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v3/log"
)

func main() {
	environment.Load()
	if app, err := internal.NewApp(); err == nil {
		log.Fatal(app.Listen(":" + environment.Port))
	}
}
