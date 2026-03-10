package viewmodel

// NavItem defines a top-level navigation entry.
type NavItem struct {
	ID    string
	Label string
	Path  string
	Order int
}

// NavbarSlot defines a supported navbar contribution area.
type NavbarSlot string

const (
	NavbarSlotStart NavbarSlot = "start"
	NavbarSlotEnd   NavbarSlot = "end"
)

// NavbarSlotContribution defines a renderable item for navbar slots.
type NavbarSlotContribution struct {
	ID       string
	Slot     NavbarSlot
	Template string
	Order    int
	Data     any
}

// ToolbarSlot defines a supported toolbar contribution area.
type ToolbarSlot string

const (
	ToolbarSlotStart ToolbarSlot = "start"
	ToolbarSlotEnd   ToolbarSlot = "end"
)

// ToolbarSlotContribution defines a renderable item for toolbar slots.
type ToolbarSlotContribution struct {
	ID       string
	Slot     ToolbarSlot
	Template string
	Order    int
	Data     any
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
