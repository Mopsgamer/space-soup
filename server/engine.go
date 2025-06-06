package internal

import (
	"fmt"
	"io/fs"
	"math"
	"net/http"
	"time"

	"github.com/Mopsgamer/space-soup/server/environment"
	"github.com/Mopsgamer/space-soup/server/soup"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

// Initialize the view engine.
func NewAppHtmlEngine(embedFS fs.FS, directory string) *html.Engine {
	var engine *html.Engine
	if embedFS == nil {
		engine = html.New(directory, ".html")
	} else {
		embedTemplates, _ := fs.Sub(embedFS, directory)
		engine = html.NewFileSystem(http.FS(embedTemplates), ".html")
	}

	if environment.BuildModeValue == environment.BuildModeDevelopment {
		engine.Reload(true)
	}

	engine.AddFuncMap(map[string]interface{}{
		"degrees": soup.DegreesFromRadians,
		"radians": soup.RadiansFromDegrees,
		"percent": func(current, total int) string {
			return fmt.Sprintf("%.2f%%", float64(current)/float64(total)*100)
		},
		"add": func(n ...int) (result int) {
			for _, v := range n {
				result += v
			}
			return
		},
		"seq": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
		},
		"concatString": func(v ...string) string {
			result := ""
			for _, str := range v {
				result += str
			}
			return result
		},
		"absMinus": func(a, b float64) float64 {
			return math.Abs(a - b)
		},
		"jsonTime": func(t time.Time) string {
			return t.Format("2006-01-02T15:04:05.000Z")
		},
		"isString": satisfies[string],
		"isMap":    satisfies[fiber.Map],
		"newMap": func(args ...any) fiber.Map {
			result := fiber.Map{}
			for i := 0; i < len(args)-1; i += 2 {
				k := args[i].(string)
				v := args[i+1]
				result[k] = v
			}
			return result
		},
		"newArr": func(args ...any) []any {
			return args
		},
	})

	return engine
}

func satisfies[T any](v any) bool {
	_, ok := v.(T)
	return ok
}
