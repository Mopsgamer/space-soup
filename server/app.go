package internal

import (
	"bytes"
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

var VisualizationInMemoryFileMap = VisualizationCache{}

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

func useHttpPageTable(
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

func useVisualDeclRasc(tests []soup.MovementTest, ctl controller_http.ControllerHttp, file []byte) (hashStr string, err error) {
	hashStr = ""
	req := new(model_http.TableImage)
	err = ctl.BindAll(req)
	if err != nil {
		return
	}

	testsRanged := []soup.MovementTest{}
	if rangeList, errParse := expandrange.Parse(req.Range); errParse != nil {
		err = errors.Join(ErrInvalidIndiceRange, errParse)
		return
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

	imageWriterTo, err := soup.Visualize(soup.VisualizeConfig{
		XLabel:      "Right Ascension",
		YLabel:      "Declination",
		Tests:       testsRanged,
		Description: description,
		GetXY: func(m soup.Movement) (x float64, y float64) {
			x, y = m.Alpha, m.Delta
			return
		},
	})
	if err != nil {
		return
	}

	buf := []byte{}
	imageBuffer := bytes.NewBuffer(buf)
	_, err = imageWriterTo.WriteTo(imageBuffer)
	if err != nil {
		return
	}
	imageBytes := imageBuffer.Bytes()
	_, err = ctl.Ctx.Write(imageBytes)
	if err != nil {
		return
	}

	if file != nil {
		hashStr = HashString(file)
		VisualizationInMemoryFileMap.Add(hashStr, imageBytes)
		log.Info(hashStr)
	}
	return
}

var (
	tests          []soup.MovementTest
	testsPaginated [][]soup.MovementTest
)

// Initialize gofiber application, including DB and view engine.
func NewApp(embedFS fs.FS) (app *fiber.App, err error) {
	engine := NewAppHtmlEngine(embedFS, "client/templates")
	tests, testsPaginated, err = soup.CheckOrbitList()
	if err != nil {
		return
	}
	app = fiber.New(fiber.Config{
		Views:             engine,
		PassLocalsToViews: true,
	})

	app.Use(logger.New())

	// static
	app.Get("/static/*", static.New("./client/static", static.Config{Browse: true}))
	app.Get("/partials*", static.New("./client/templates/partials", static.Config{Browse: true}))

	// pages
	var noRedirect controller_http.RedirectCompute = func(ctl controller_http.ControllerHttp, bind *fiber.Map) string { return "" }

	app.Get("/", useHttpPage("homepage", &fiber.Map{"Title": "Home", "IsHomePage": true}, noRedirect, "partials/main"))
	app.Get("/manually", useHttpPage("manually", &fiber.Map{"Title": "Calculate manually", "IsManually": true}, noRedirect, "partials/main"))
	app.Get("/parse", useHttpPage("parse", &fiber.Map{"Title": "Upload file", "IsFile": true}, noRedirect, "partials/main"))
	app.Get("/table/image", useHttp(func(ctl controller_http.ControllerHttp) error {
		_, err = useVisualDeclRasc(tests, ctl, nil)
		return err
	}))
	app.Get("/table/:page", useHttpPageTable("table", &fiber.Map{"Title": "Table"}, noRedirect, "partials/main"))
	app.Get("/table", func(ctx fiber.Ctx) error { return ctx.Redirect().To("/table/1") })
	app.Post("/table/expand", useHttp(func(ctl controller_http.ControllerHttp) error {
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
	app.Get("/terms", useHttpPage("terms", &fiber.Map{"Title": "Terms", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/privacy", useHttpPage("privacy", &fiber.Map{"Title": "Privacy", "CenterContent": true}, noRedirect, "partials/main"))
	app.Get("/acknowledgements", useHttpPage("acknowledgements", &fiber.Map{"Title": "Acknowledgements"}, noRedirect, "partials/main"))

	app.Get("/alg/parse/image/:hash.png", useHttp(func(ctl controller_http.ControllerHttp) error {
		hash := ctl.Ctx.Params("hash", "x")
		imageCache, ok := VisualizationInMemoryFileMap[hash]
		if !ok {
			return ctl.Ctx.SendStatus(fiber.StatusMovedPermanently)
		}
		VisualizationInMemoryFileMap[hash].Live()
		log.Info(hash)
		return ctl.Ctx.Send(imageCache.Bytes)
	}))
	app.Post("/alg/parse/image", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		pFile, err := ctl.Ctx.FormFile("document")
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-calc")
		}

		file, err := pFile.Open()
		if err != nil {
			return err
		}
		buf := []byte{}
		if _, err = file.Read(buf); err != nil {
			file.Close()
			return err
		}
		file.Close()

		movementTestList, err := req.MovementTestList(*pFile)
		if err != nil {
			return ctl.RenderDanger(err.Error(), "err-calc")
		}

		hashStr, err := useVisualDeclRasc(movementTestList, ctl, buf)
		if err != nil {
			return err
		}

		log.Info(hashStr)
		ctl.HTMXRedirect("/alg/parse/image/" + hashStr + ".png")
		return nil
	}))
	app.Post("/alg/parse", useHttp(func(ctl controller_http.ControllerHttp) error {
		req := new(model_http.OrbitInputFile)
		if err := ctl.BindAll(req); err != nil {
			return ctl.RenderInternalError(err.Error(), "err-request")
		}

		pFile, err := ctl.Ctx.FormFile("document")
		if err != nil {
			return ctl.RenderInternalError(err.Error(), "err-calc")
		}

		movementTestList, err := req.MovementTestList(*pFile)
		if err != nil {
			return ctl.RenderDanger(err.Error(), "err-calc")
		}

		return ctl.Ctx.Render("partials/table", fiber.Map{
			"ExpandTable": false,
			"Table":       movementTestList,
			"Page":        1,
			"PageMax":     1,
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
			"ExpandTable": false,
			"Table":       movementTestList,
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
