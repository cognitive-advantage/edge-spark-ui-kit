package uikit

import (
	"embed"
	"html/template"
	"io/fs"

	"github.com/cognitive-advantage/edge-spark-ui-kit/renderer"
)

// TemplatesFS contains the shared UI template files.
//
//go:embed templates
var TemplatesFS embed.FS

// NewRenderer returns a Gin HTML renderer backed by embedded UI kit templates.
func NewRenderer(funcMap template.FuncMap, cacheEnabled bool) *renderer.Engine {
	templateFS, err := fs.Sub(TemplatesFS, "templates")
	if err != nil {
		panic(err)
	}

	return renderer.NewEngine(templateFS, funcMap, cacheEnabled)
}
