package controller_http

import (
	"html/template"
	"strings"

	"github.com/Mopsgamer/space-soup/server/controller"
	"github.com/Mopsgamer/space-soup/server/environment"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"github.com/google/safehtml"
)

// Should return redirect path or empty string.
type RedirectCompute func(ctl ControllerHttp, bind *fiber.Map) string

// Render a page using a template.
// Special
func (r ControllerHttp) RenderPage(templatePath string, bind *fiber.Map, redirect RedirectCompute, layouts ...string) error {
	bindx := r.MapPage(bind)
	if path := redirect(r, bind); path != "" {
		return r.Ctx.Redirect().To(path)
	}
	if title, ok := (*bind)["Title"].(string); ok {
		bindx["Title"] = environment.AppName + " - " + title
	}
	return r.Ctx.Render(templatePath, bindx, layouts...)
}

func (ctl ControllerHttp) MapPage(bind *fiber.Map) fiber.Map {
	bindx := fiber.Map{
		"AppName":     environment.AppName,
		"GitHubRepo":  environment.GitHubRepo,
		"DenoJson":    environment.DenoJson,
		"GoMod":       environment.GoMod,
		"GitHash":     environment.GitHash,
		"GitHashLong": environment.GitHashLong,
	}

	bindx = controller.MapMerge(&bindx, bind)
	return bindx
}

func (ctl ControllerHttp) RenderString(template string, bind any) (string, error) {
	return controller.RenderString(ctl.Ctx.App(), template, bind)
}

func escapeNewLines(message string) string {
	message = safehtml.HTMLEscaped(message).String()
	message = strings.ReplaceAll(message, "\n", "<br />")
	return message
}

func wrapRenderNotice(r ControllerHttp, templateName, message, id string) error {
	return r.Ctx.Render(templateName, fiber.Map{
		"Id":      id,
		"Message": template.HTML(escapeNewLines(message)),
	})
}

// Renders the danger message html element.
func (ctl ControllerHttp) RenderInternalError(message string, id string) error {
	log.Error(message)
	return ctl.RenderDanger(message, id)
}

// Renders the danger message html element.
func (ctl ControllerHttp) RenderDanger(message, id string) error {
	return wrapRenderNotice(ctl, "partials/danger", message, id)
}

// Renders the warning message html element.
func (ctl ControllerHttp) RenderWarning(message, id string) error {
	return wrapRenderNotice(ctl, "partials/warning", message, id)
}

// Renders the success message html element.
func (ctl ControllerHttp) RenderSuccess(message, id string) error {
	return wrapRenderNotice(ctl, "partials/success", message, id)
}
