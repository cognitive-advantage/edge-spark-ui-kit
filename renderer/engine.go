package renderer

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/render"
)

// Engine renders UI kit templates in a Hugo-style pattern:
// parse the requested template plus shared layout/partials/components.
type Engine struct {
	templateFS   fs.FS
	funcMap      template.FuncMap
	cacheEnabled bool

	cache map[string]*template.Template
	mu    sync.RWMutex
}

// NewEngine creates a shared UI renderer.
func NewEngine(templateFS fs.FS, funcMap template.FuncMap, cacheEnabled bool) *Engine {
	copiedFuncMap := template.FuncMap{}
	for key, fn := range funcMap {
		copiedFuncMap[key] = fn
	}

	return &Engine{
		templateFS:   templateFS,
		funcMap:      copiedFuncMap,
		cacheEnabled: cacheEnabled,
		cache:        make(map[string]*template.Template),
	}
}

// Instance satisfies gin's HTMLRender interface.
func (e *Engine) Instance(name string, data interface{}) render.Render {
	return &instance{
		name:   name,
		data:   data,
		engine: e,
	}
}

type instance struct {
	name   string
	data   interface{}
	engine *Engine
}

func (i *instance) Render(w http.ResponseWriter) error {
	i.WriteContentType(w)

	tmpl, err := i.engine.parse(i.name)
	if err != nil {
		return err
	}

	if isSupportTemplate(i.name) {
		return tmpl.ExecuteTemplate(w, i.name, i.data)
	}

	return tmpl.ExecuteTemplate(w, "layouts/baseof.html", i.data)
}

func (i *instance) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}

func (e *Engine) parse(name string) (*template.Template, error) {
	cacheKey := name
	if e.cacheEnabled {
		e.mu.RLock()
		cached := e.cache[cacheKey]
		e.mu.RUnlock()
		if cached != nil {
			return cached, nil
		}
	}

	files, err := e.collectFiles(name)
	if err != nil {
		return nil, err
	}

	tmpl := template.New("").Funcs(e.funcMap)
	parsed, err := tmpl.ParseFS(e.templateFS, files...)
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	if e.cacheEnabled {
		e.mu.Lock()
		e.cache[cacheKey] = parsed
		e.mu.Unlock()
	}

	return parsed, nil
}

func (e *Engine) collectFiles(name string) ([]string, error) {
	normalized := path.Clean(strings.TrimPrefix(name, "/"))
	if normalized == "." {
		return nil, fmt.Errorf("template name is required")
	}

	if _, err := fs.Stat(e.templateFS, normalized); err != nil {
		return nil, fmt.Errorf("template not found: %s", normalized)
	}

	files := []string{normalized}

	if !isSupportTemplate(normalized) {
		if _, err := fs.Stat(e.templateFS, "layouts/baseof.html"); err == nil {
			files = append(files, "layouts/baseof.html")
		}
	}

	partials, err := walkHTMLFiles(e.templateFS, "partials")
	if err == nil {
		files = append(files, partials...)
	}

	components, err := walkHTMLFiles(e.templateFS, "components")
	if err == nil {
		files = append(files, components...)
	}

	if !isSupportTemplate(normalized) {
		layouts, err := walkHTMLFiles(e.templateFS, "layouts")
		if err == nil {
			files = append(files, layouts...)
		}
	}

	return dedupe(files), nil
}

func walkHTMLFiles(fsys fs.FS, root string) ([]string, error) {
	var files []string

	err := fs.WalkDir(fsys, root, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		files = append(files, current)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

func dedupe(items []string) []string {
	seen := make(map[string]struct{}, len(items))
	result := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func isSupportTemplate(name string) bool {
	return strings.HasPrefix(name, "partials/") || strings.HasPrefix(name, "components/")
}
