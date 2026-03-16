# edge-spark ui-kit

Shared UI composition layer for Edge Spark Go web apps.

## Goals

- Centralize common layout templates (shell, header, footer, breadcrumbs, flash).
- Keep application-specific pages in each app repository/module.
- Enforce clear separation of concerns:
  - domain logic in app domain packages
  - route/handler orchestration in app handlers
  - presentation primitives in this module

## Included

- Embedded templates under `templates/`
- Embedded theme assets under `assets/`
- Gin-compatible HTML renderer in `renderer/`
- Shared baseline view models in `viewmodel/`

## Serving Ui-Kit Assets

Mount assets once in your Gin server:

```go
r.StaticFS("/ui-kit", uikit.AssetsHTTPFS())
```

This exposes:

- `/ui-kit/css/uikit.css`
- `/ui-kit/js/theme.js`

The shared `layouts/baseof.html` already includes those files before app CSS.

## Quick Start

```go
import (
    "html/template"

    uikit "github.com/cognitive-advantage/edge-spark-ui-kit"
)

funcMap := template.FuncMap{}
r.HTMLRender = uikit.NewRenderer(funcMap, true)
```

Render a page with:

```go
c.HTML(http.StatusOK, "layouts/page.html", gin.H{
    "Title": "Dashboard",
    "AppName": "Edge Spark Factory",
    "NavItems": []map[string]string{
        {"Label": "Builds", "Path": "/build/builds"},
    },
    "ContentTemplate": "build/builds/page.html",
})
```

## CI Recommendations

- Run `go test ./...` in this module.
- Require downstream app integration tests against a pinned ui-kit version.
- Version with semver tags and update apps intentionally.

## Build

Build a distributable source bundle to repository root `./bin`:

```bash
make build
```
