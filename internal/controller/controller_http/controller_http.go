package controller_http

import (
	"github.com/gofiber/fiber/v3"
)

type ControllerHttp struct {
	Ctx fiber.Ctx
}

type Response interface {
	HandleHtmx(ctl ControllerHttp) error
}
