package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
	"time"

	"github.com/Mopsgamer/space-soup/server/controller"
	"github.com/Mopsgamer/space-soup/server/controller/controller_http"
	"github.com/Mopsgamer/space-soup/server/controller/model_http"
	"github.com/Mopsgamer/space-soup/server/soup"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
	expandrange "github.com/n0madic/expand-range"
)

var (
	ErrInvalidIndiceRange = errors.New("invalid indice range")
	ErrInvalidPageRange   = errors.New("invalid page range")
)

// Initialize gofiber application, including DB and view engine.
func NewApp(embedFS fs.FS) (app *fiber.App, err error) {
	engine := NewAppHtmlEngine(embedFS, "client/templates")
	tests, testsPaginated, err := soup.CheckOrbitList()
	if err != nil {
		return
	}
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

	UseVisualDeclRasc := func() fiber.Handler {
		return UseHttp(func(ctl controller_http.ControllerHttp) error {
			req := new(model_http.TableImage)
			err := ctl.BindAll(req)
			if err != nil {
				return err
			}

			testsRanged := []soup.MovementTest{}
			if rangeList, err := expandrange.Parse(req.Range); err != nil {
				return errors.Join(ErrInvalidIndiceRange, err)
			} else if req.Range == "" {
				testsRanged = tests
			} else {
				for _, i := range rangeList {
					testsRanged = append(testsRanged, tests[i])
				}
			}

			if req.IsDownload {
				ctl.Ctx.Set("Content-Type", "application/octet-stream")
				ctl.Ctx.Set("Content-Disposition", "attachment; filename=orbits-decl-rasc.png")
			}

			if len(req.Description) > 400 {
				req.Description = req.Description[:400]
			}
			description := strings.ReplaceAll(req.Description, "\n", " | ")

			writerTo, err := soup.VisualizeDeclRasc(soup.VisualizeConfig{
				Tests:       testsRanged,
				Description: description,
			})
			if err != nil {
				return err
			}

			_, err = writerTo.WriteTo(ctl.Ctx.RequestCtx().Response.BodyWriter())
			return err
		})
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
			req := new(model_http.TablePage)
			err := ctl.BindAll(req)
			if err != nil {
				return err
			}
			if req.Page > len(testsPaginated) || req.Page < 1 {
				return ErrInvalidPageRange
			}
			bindx["ExpandTable"] = req.ExpandTable
			bindx["Table"] = testsPaginated[req.Page-1]
			bindx["Page"] = req.Page
			bindx["PageMax"] = len(testsPaginated)
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
	app.Get("/manually", UseHttpPage("manually", &fiber.Map{"Title": "Calculate manually", "IsManually": true}, noRedirect, "partials/main"))
	app.Get("/file", UseHttpPage("file", &fiber.Map{"Title": "Upload file", "IsFile": true}, noRedirect, "partials/main"))
	app.Get("/table/image", UseVisualDeclRasc())
	app.Get("/table/:page", UseHttpPageTable("table", &fiber.Map{"Title": "Table"}, noRedirect, "partials/main"))
	app.Get("/table", func(ctx fiber.Ctx) error { return ctx.Redirect().To("/table/1") })
	app.Post("/table/expand", UseHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.TablePage)
		err := ctl.BindAll(req)
		if err != nil {
			return err
		}
		if req.ExpandTable {
			ctl.Ctx.Cookie(&fiber.Cookie{
				Name:    "expanded",
				Expires: time.Now(),
				Value:   "",
			})
		} else {
			ctl.Ctx.Cookie(&fiber.Cookie{
				Name:    "expanded",
				Expires: time.Now().Add(time.Hour * 24 * 30),
				Value:   "true",
			})
		}

		ctl.HTMXRefresh()
		return nil
	}))
	app.Get("/terms", UseHttpPage("terms", &fiber.Map{"Title": "Terms", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/privacy", UseHttpPage("privacy", &fiber.Map{"Title": "Privacy", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/acknowledgements", UseHttpPage("acknowledgements", &fiber.Map{"Title": "Acknowledgements"}, noRedirect, "partials/main"))

	// calc
	app.Post("/process/file", UseHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		pFile, err := ctl.Ctx.FormFile("document")
		if err != nil {
			log.Info(string(ctl.Ctx.BodyRaw()))
			return ctl.RenderInternalError(err.Error(), "err-calc")
		}

		movementList, err := req.MovementList(*pFile)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-calc")
		}

		return ctl.Ctx.Render("partials/table", fiber.Map{
			"ExpandTable": false,
			"Table":       movementList,
			"Page":        1,
			"PageMax":     1,
		})
	}))
	app.Post("/process", UseHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInput)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		input, err := req.Input()
		if err != nil {
			return ctl.RenderDanger(err.Error(), "err-calc")
		}
		result := soup.NewMovement(input)

		return ctl.Ctx.Render("partials/table", fiber.Map{
			"ExpandTable": false,
			"Table": []soup.MovementTest{{
				Input:  input,
				Actual: result,
			}},
			"Page":    1,
			"PageMax": 1,
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
