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
	options      Options

	cache map[string]*template.Template
	mu    sync.RWMutex
}

// Options controls how the renderer discovers and executes templates.
type Options struct {
	BaseTemplate      string
	SupportPrefixes   []string
	IncludeDirs       []string
	DomainSupportDirs []string
}

// NewEngine creates a shared UI renderer.
func NewEngine(templateFS fs.FS, funcMap template.FuncMap, cacheEnabled bool) *Engine {
	return NewEngineWithOptions(templateFS, funcMap, cacheEnabled, DefaultOptions())
}

// NewEngineWithOptions creates a shared UI renderer with custom conventions.
func NewEngineWithOptions(templateFS fs.FS, funcMap template.FuncMap, cacheEnabled bool, options Options) *Engine {
	copiedFuncMap := template.FuncMap{}
	for key, fn := range funcMap {
		copiedFuncMap[key] = fn
	}

	if options.BaseTemplate == "" {
		options = DefaultOptions()
	}

	return &Engine{
		templateFS:   templateFS,
		funcMap:      copiedFuncMap,
		cacheEnabled: cacheEnabled,
		options:      options,
		cache:        make(map[string]*template.Template),
	}
}

// DefaultOptions returns the default template conventions used by the shared ui-kit templates.
func DefaultOptions() Options {
	return Options{
		BaseTemplate:      "layouts/baseof.html",
		SupportPrefixes:   []string{"partials/", "components/"},
		IncludeDirs:       []string{"partials", "components", "layouts"},
		DomainSupportDirs: nil,
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

	if i.engine.isSupportTemplate(i.name) {
		return tmpl.ExecuteTemplate(w, i.name, i.data)
	}

	return tmpl.ExecuteTemplate(w, i.engine.options.BaseTemplate, i.data)
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

	if !e.isSupportTemplate(normalized) {
		if _, err := fs.Stat(e.templateFS, e.options.BaseTemplate); err == nil {
			files = append(files, e.options.BaseTemplate)
		}
	}

	for _, includeDir := range e.options.IncludeDirs {
		if e.isSupportTemplate(normalized) && includeDir == path.Dir(e.options.BaseTemplate) {
			continue
		}

		foundFiles, err := walkHTMLFiles(e.templateFS, includeDir)
		if err == nil {
			files = append(files, foundFiles...)
		}
	}

	for _, domainDir := range e.options.DomainSupportDirs {
		domainFiles, err := walkDomainSupportFiles(e.templateFS, domainDir)
		if err == nil {
			files = append(files, domainFiles...)
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

func walkDomainSupportFiles(fsys fs.FS, root string) ([]string, error) {
	var files []string

	err := fs.WalkDir(fsys, root, func(current string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(d.Name(), ".html") {
			return nil
		}

		baseName := path.Base(current)
		if baseName == "row_template.html" || strings.HasPrefix(baseName, "_") {
			files = append(files, current)
		}

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

func (e *Engine) isSupportTemplate(name string) bool {
	for _, prefix := range e.options.SupportPrefixes {
		if strings.HasPrefix(name, prefix) {
			return true
		}
	}

	return false
}
