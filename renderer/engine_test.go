package renderer

import (
	"io/fs"
	"testing"
	"testing/fstest"
)

func TestCollectFilesIncludesBaseLayoutForPages(t *testing.T) {
	t.Parallel()

	kit := fstest.MapFS{
		"layouts/baseof.html":       {Data: []byte("{{ define \"layouts/baseof.html\" }}{{ end }}")},
		"layouts/page.html":         {Data: []byte("{{ define \"layouts/page.html\" }}{{ end }}")},
		"partials/header.html":      {Data: []byte("{{ define \"partials/header.html\" }}{{ end }}")},
		"components/flash.html":     {Data: []byte("{{ define \"components/flash.html\" }}{{ end }}")},
		"pages/dashboard.html":      {Data: []byte("{{ define \"pages/dashboard.html\" }}{{ end }}")},
		"partials/footer.html":      {Data: []byte("{{ define \"partials/footer.html\" }}{{ end }}")},
		"components/table.html":     {Data: []byte("{{ define \"components/table.html\" }}{{ end }}")},
		"partials/breadcrumbs.html": {Data: []byte("{{ define \"partials/breadcrumbs.html\" }}{{ end }}")},
	}

	engine := NewEngine(fs.FS(kit), nil, false)
	files, err := engine.collectFiles("pages/dashboard.html")
	if err != nil {
		t.Fatalf("collect files failed: %v", err)
	}

	if !contains(files, "layouts/baseof.html") {
		t.Fatalf("expected layouts/baseof.html in parsed files, got: %v", files)
	}
}

func TestCollectFilesSkipsBaseLayoutForPartials(t *testing.T) {
	t.Parallel()

	kit := fstest.MapFS{
		"layouts/baseof.html":  {Data: []byte("{{ define \"layouts/baseof.html\" }}{{ end }}")},
		"partials/header.html": {Data: []byte("{{ define \"partials/header.html\" }}{{ end }}")},
	}

	engine := NewEngine(fs.FS(kit), nil, false)
	files, err := engine.collectFiles("partials/header.html")
	if err != nil {
		t.Fatalf("collect files failed: %v", err)
	}

	if contains(files, "layouts/baseof.html") {
		t.Fatalf("did not expect layouts/baseof.html for partial rendering, got: %v", files)
	}
}

func contains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
