package internal

import (
	"time"

	"github.com/Mopsgamer/space-soup/internal/environment"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

// Initialize the view engine.
func NewAppHtmlEngine() *html.Engine {
	engine := html.New("./web/templates", ".html")

	if environment.Environment == environment.EnvironmentDevelopment {
		engine.Reload(true)
	}

	engine.AddFuncMap(map[string]interface{}{
		"concatString": func(v ...string) string {
			result := ""
			for _, str := range v {
				result += str
			}
			return result
		},
		"jsonTime": func(t time.Time) string {
			return t.Format("2006-01-02T15:04:05.000Z")
		},
		"isString": satisfies[string],
		"isMap":    satisfies[fiber.Map],
		"newMap": func(args ...any) fiber.Map {
			result := fiber.Map{}
			for i := 0; i < len(args)-1; i = i + 2 {
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
