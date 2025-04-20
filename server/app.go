package internal

import (
	"fmt"
	"io/fs"
	"strconv"

	"github.com/Mopsgamer/space-soup/server/controller"
	"github.com/Mopsgamer/space-soup/server/controller/controller_http"
	"github.com/Mopsgamer/space-soup/server/controller/model_http"
	"github.com/Mopsgamer/space-soup/server/soup"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

// Initialize gofiber application, including DB and view engine.
func NewApp(embedFS fs.FS) (app *fiber.App, err error) {
	engine := NewAppHtmlEngine(embedFS, "client/templates")
	app = fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())

	UseHttp := func(handler func(ctl controller_http.ControllerHttp) error) fiber.Handler {
		return func(ctx fiber.Ctx) error {
			ctl := controller_http.ControllerHttp{
				Ctx: ctx,
			}

			return handler(ctl)
		}
	}

	UseHttpPage := func(
		templatePath string,
		bind *fiber.Map,
		redirect controller_http.RedirectCompute,
		layouts ...string,
	) fiber.Handler {
		bindx := fiber.Map{
			"Title": "?",
		}
		bindx = controller.MapMerge(&bindx, bind)
		return UseHttp(func(ctl controller_http.ControllerHttp) error {
			return ctl.RenderPage(
				templatePath,
				&bindx,
				redirect,
				layouts...,
			)
		})
	}

	table, err := soup.CheckOrbitList()
	if err != nil {
		return
	}

	UseHttpPageTable := func(
		templatePath string,
		bind *fiber.Map,
		redirect controller_http.RedirectCompute,
		layouts ...string,
	) fiber.Handler {
		bindx := fiber.Map{
			"Title": "?",
		}
		bindx = controller.MapMerge(&bindx, bind)
		return UseHttp(func(ctl controller_http.ControllerHttp) error {
			page64, err := strconv.ParseInt(ctl.Ctx.Params("page", "1"), 0, 64)
			page := int(page64)
			if err != nil || page > len(table) || page < 1 {
				return ctl.Ctx.Next()
			}
			bindx["Table"] = table[page-1]
			bindx["Page"] = page
			bindx["PageMax"] = len(table)
			return ctl.RenderPage(
				templatePath,
				&bindx,
				redirect,
				layouts...,
			)
		})
	}

	// static
	app.Get("/static/*", static.New("./client/static", static.Config{Browse: true}))
	app.Get("/partials*", static.New("./client/templates/partials", static.Config{Browse: true}))

	// pages
	var noRedirect controller_http.RedirectCompute = func(ctl controller_http.ControllerHttp, bind *fiber.Map) string { return "" }

	app.Get("/", UseHttpPage("homepage", &fiber.Map{"Title": "Home", "IsHomePage": true}, noRedirect, "partials/main"))
	app.Get("/calc", UseHttpPage("calc", &fiber.Map{"Title": "Calculate", "IsCalc": true}, noRedirect, "partials/main"))
	app.Get("/table/:page", UseHttpPageTable("table", &fiber.Map{"Title": "Table"}, noRedirect, "partials/main"))
	app.Get("/table", func(ctx fiber.Ctx) error { return ctx.Redirect().To("/table/1") })
	app.Get("/test-table/:page", UseHttpPageTable("test-table", &fiber.Map{"Title": "Test table"}, noRedirect, "partials/main"))
	app.Get("/test-table", func(ctx fiber.Ctx) error { return ctx.Redirect().To("/test-table/1") })
	app.Get("/terms", UseHttpPage("terms", &fiber.Map{"Title": "Terms", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/privacy", UseHttpPage("privacy", &fiber.Map{"Title": "Privacy", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/acknowledgements", UseHttpPage("acknowledgements", &fiber.Map{"Title": "Acknowledgements"}, noRedirect, "partials/main"))

	// calc
	app.Post("/process", UseHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInput)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError("err-request")
		}

		meteor, err := req.Movement()
		if err != nil {
			return ctl.RenderDanger(err.Error(), "err-calc")
		}

		return ctl.Ctx.Render("partials/meteoroid_movement", fiber.Map{
			"Meteor": meteor,
		})
	}))

	app.Use(UseHttpPage("partials/x", &fiber.Map{
		"Title":         fmt.Sprintf("%d", fiber.StatusNotFound),
		"StatusCode":    fiber.StatusNotFound,
		"StatusMessage": fiber.ErrNotFound.Message,
		"CenterContent": true,
	}, noRedirect, "partials/main"))

	return
}
