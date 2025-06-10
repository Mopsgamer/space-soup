package internal

import (
	"fmt"
	"io/fs"
	"time"

	"github.com/Mopsgamer/space-soup/server/controller"
	"github.com/Mopsgamer/space-soup/server/controller/controller_http"
	"github.com/Mopsgamer/space-soup/server/controller/model_http"
	"github.com/Mopsgamer/space-soup/server/soup"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/static"
)

var AlgCache = soup.FileHashCacheMap{}

func useHttp(handler func(ctl controller_http.ControllerHttp) error) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		ctl := controller_http.ControllerHttp{
			Ctx: ctx,
		}

		return handler(ctl)
	}
}

func useHttpPage(
	templatePath string,
	bind *fiber.Map,
	redirect controller_http.RedirectCompute,
	layouts ...string,
) fiber.Handler {
	bindx := fiber.Map{
		"Title": "?",
	}
	bindx = controller.MapMerge(&bindx, bind)
	return useHttp(func(ctl controller_http.ControllerHttp) error {
		return ctl.RenderPage(
			templatePath,
			&bindx,
			redirect,
			layouts...,
		)
	})
}

var (
	tests          []soup.MovementTest
	testsPaginated [][]soup.MovementTest
	testsFailCount int
)

// Initialize gofiber application, including DB and view engine.
func NewApp(embedFS fs.FS) (app *fiber.App, err error) {
	tests, err = soup.CheckOrbitList(soup.FileContentEtalon, soup.FileContentInput)
	if err != nil {
		return
	}
	testsPaginated = soup.Paginate(tests)
	for _, test := range tests {
		if test.AssertionResult.Has(soup.TestResultFailed) {
			testsFailCount += 1
		}
	}

	engine := NewAppHtmlEngine(embedFS, "client/templates")

	app = fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())

	// static
	app.Get("/static/*", static.New("./client/static", static.Config{Browse: true}))
	app.Get("/partials/*", static.New("./client/templates/partials", static.Config{Browse: true}))

	// pages
	var noRedirect controller_http.RedirectCompute = func(ctl controller_http.ControllerHttp, bind *fiber.Map) string { return "" }

	app.Get("/", useHttpPage("homepage", &fiber.Map{"Title": "Home", "IsHomePage": true}, noRedirect, "partials/main"))
	app.Get("/manually", useHttpPage("manually", &fiber.Map{"Title": "Calculate manually", "IsManually": true}, noRedirect, "partials/main"))
	app.Get("/parse", useHttpPage("parse", &fiber.Map{"Title": "Upload file", "IsFile": true}, noRedirect, "partials/main"))
	app.Get("/tests.png", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		writerTo, err := soup.VisualizeAlphaDelta(tests, "")
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		_, err = writerTo.WriteTo(ctl.Ctx)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		return nil
	}))
	app.Get("/tests.xlsx", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		bytes, err := soup.NewFileExcelBytes(tests, false)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		_, err = ctl.Ctx.Write(bytes)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		return nil
	}))
	app.Get("/tests/:page", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.TablePage)
		err := ctl.BindAll(req)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		if err := req.Page.Validate(len(testsPaginated)); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}
		return ctl.RenderPage(
			"tests",
			&fiber.Map{
				"Title":          "Table",
				"IsTests":        true,
				"TestsFailCount": testsFailCount,
				"ExpandTable":    req.TableState,
				"Table":          testsPaginated[req.Page-1],
				"TableMax":       len(tests),
				"Page":           int(req.Page),
				"PageMax":        len(testsPaginated),
			},
			noRedirect,
			"partials/main",
		)
	}))
	app.Get("/tests", func(ctx fiber.Ctx) error { return ctx.Redirect().To("/tests/1") })
	app.Get("/terms", useHttpPage("terms", &fiber.Map{"Title": "Terms", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/privacy", useHttpPage("privacy", &fiber.Map{"Title": "Privacy", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/acknowledgements", useHttpPage("acknowledgements", &fiber.Map{"Title": "Acknowledgements"}, noRedirect, "partials/main"))
	app.Get("/alg/parse/cache/:hash.xlsx", useHttp(func(ctl controller_http.ControllerHttp) error {
		hash := ctl.Ctx.Params("hash", "x")
		cache, cacheExists := AlgCache[hash]
		if !cacheExists {
			return ctl.Ctx.SendStatus(fiber.StatusMovedPermanently)
		}

		if err := AlgCache.Live(hash); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		return ctl.Ctx.Send(cache.ExcelBytes)
	}))
	app.Get("/alg/parse/cache/:hash.png", useHttp(func(ctl controller_http.ControllerHttp) error {
		hash := ctl.Ctx.Params("hash", "x")
		cache, cacheExists := AlgCache[hash]
		if !cacheExists {
			return ctl.Ctx.SendStatus(fiber.StatusMovedPermanently)
		}

		if err := AlgCache.Live(hash); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		return ctl.Ctx.Send(cache.PlotImageBytes)
	}))

	// api
	app.Post("/tests/expand", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.TableSetMode)
		err := ctl.BindAll(req)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		newCookie := fiber.Cookie{Name: "expand-state", Expires: time.Now()}

		if req.ExpandTable != model_http.TableStateNormal {
			newCookie.Expires = newCookie.Expires.Add(time.Hour * 24 * 30)
			newCookie.Value = string(req.ExpandTable)
		}

		ctl.Ctx.Cookie(&newCookie)

		ctl.HTMXRefresh()
		return nil
	}))
	app.Post("/alg/parse/cache", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		formFile, err := ctl.Ctx.FormFile("document")
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		formFileBytes, err := soup.NewFormFileBytes(formFile)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		hash := soup.HashString(formFileBytes)
		cache, cacheExists := AlgCache[hash]
		if !cacheExists {
			tests, err := soup.NewMovementTestsFromFile(formFile, req.FileType)
			if err != nil {
				return ctl.RenderInternalError(err.Error(), "err-request")
			}
			newCache, err := soup.NewCache(tests, req.Range, "")
			if err != nil {
				return ctl.RenderInternalError(err.Error(), "err-request")
			}
			cache = newCache
		}
		AlgCache.Add(hash, cache)

		return ctl.Ctx.SendString(hash)
	}))
	app.Post("/alg/parse/:page", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		formFile, err := ctl.Ctx.FormFile("document")
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		formFileBytes, err := soup.NewFormFileBytes(formFile)
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		hash := soup.HashString(formFileBytes)
		cache, cacheExists := AlgCache[hash]
		if !cacheExists {
			tests, err := soup.NewMovementTestsFromFile(formFile, req.FileType)
			if err != nil {
				return ctl.RenderInternalError(err.Error(), "err-request")
			}
			newCache, err := soup.NewCache(tests, req.Range, "")
			if err != nil {
				return ctl.RenderInternalError(err.Error(), "err-request")
			}
			cache = newCache
		}
		AlgCache.Add(hash, cache)

		movementTestsPaginated := soup.Paginate(cache.TestList)
		if err := req.Page.Validate(len(movementTestsPaginated)); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}
		movementTests := movementTestsPaginated[req.Page-1]

		return ctl.Ctx.Render("partials/table", fiber.Map{
			"IsFile":      true,
			"FileHash":    hash,
			"ExpandTable": model_http.TableStateNormal,
			"Table":       movementTests,
			"TableMax":    len(cache.TestList),
			"Page":        int(req.Page),
			"PageMax":     len(movementTestsPaginated),
		})
	}))
	app.Post("/alg", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInput)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		input, err := req.Input()
		if err != nil {
			return ctl.RenderDanger(err.Error(), "err-calc")
		}
		result := soup.NewMovement(input)
		movementTestList := []soup.MovementTest{{
			Input:  input,
			Actual: result,
		}}

		return ctl.Ctx.Render("partials/table", fiber.Map{
			"ExpandTable": model_http.TableStateNormal,
			"Table":       movementTestList,
			"TableMax":    1,
			"Page":        1,
			"PageMax":     1,
		})
	}))

	app.Use(useHttpPage("partials/x", &fiber.Map{
		"Title":         fmt.Sprintf("%d", fiber.StatusNotFound),
		"StatusCode":    fiber.StatusNotFound,
		"StatusMessage": fiber.ErrNotFound.Message,
		"CenterContent": true,
	}, noRedirect, "partials/main"))

	return
}
