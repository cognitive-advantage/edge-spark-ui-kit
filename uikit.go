package uikit

import (
	"embed"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/cognitive-advantage/edge-spark-ui-kit/renderer"
)

// TemplatesFS contains the shared UI template files.
//
//go:embed templates
var TemplatesFS embed.FS

// AssetsFS contains shared CSS/JS assets for theming and baseline UI behavior.
//
//go:embed assets
var AssetsFS embed.FS

// NewRenderer returns a Gin HTML renderer backed by embedded UI kit templates.
func NewRenderer(funcMap template.FuncMap, cacheEnabled bool) *renderer.Engine {
	templateFS, err := TemplateFS()
	if err != nil {
		panic(err)
	}

	return renderer.NewEngine(templateFS, funcMap, cacheEnabled)
}

// TemplateFS returns the embedded templates sub-filesystem.
func TemplateFS() (fs.FS, error) {
	return fs.Sub(TemplatesFS, "templates")
}

// AssetsHTTPFS returns embedded ui-kit assets as an http.FileSystem.
func AssetsHTTPFS() http.FileSystem {
	assetsFS, err := fs.Sub(AssetsFS, "assets")
	if err != nil {
		panic(err)
	}

	return http.FS(assetsFS)
}
