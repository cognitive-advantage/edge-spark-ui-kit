package viewmodel

// NavItem defines a top-level navigation entry.
type NavItem struct {
	Label string
	Path  string
}

// BreadcrumbItem defines a breadcrumb entry.
type BreadcrumbItem struct {
	Label string
	Path  string
}

// FlashMessage defines a status banner rendered near the page content.
type FlashMessage struct {
	Kind    string
	Message string
}

// ShellData is the baseline payload expected by shared shell templates.
// Applications can embed or extend this shape in their own page view models.
type ShellData struct {
	Title       string
	AppName     string
	NavItems    []NavItem
	Breadcrumbs []BreadcrumbItem
	Flash       *FlashMessage
}
