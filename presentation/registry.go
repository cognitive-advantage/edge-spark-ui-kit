package presentation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cognitive-advantage/edge-spark-ui-kit/viewmodel"
)

// Registry stores generic presentation composition data for shared shell rendering.
type Registry struct {
	navItems        []viewmodel.NavItem
	navByID         map[string]struct{}
	navByPath       map[string]struct{}
	navbarSlots     []viewmodel.NavbarSlotContribution
	navbarSlotByID  map[string]struct{}
	toolbarSlots    []viewmodel.ToolbarSlotContribution
	toolbarSlotByID map[string]struct{}
}

// NewRegistry creates a presentation registry with empty state.
func NewRegistry() *Registry {
	return &Registry{
		navItems:        make([]viewmodel.NavItem, 0),
		navByID:         map[string]struct{}{},
		navByPath:       map[string]struct{}{},
		navbarSlots:     make([]viewmodel.NavbarSlotContribution, 0),
		navbarSlotByID:  map[string]struct{}{},
		toolbarSlots:    make([]viewmodel.ToolbarSlotContribution, 0),
		toolbarSlotByID: map[string]struct{}{},
	}
}

// AddNavItem registers a top-level nav item with dedupe and normalization.
func (r *Registry) AddNavItem(item viewmodel.NavItem) error {
	trimmedID := strings.TrimSpace(item.ID)
	if trimmedID == "" {
		return fmt.Errorf("nav id is required")
	}

	trimmedPath := strings.TrimSpace(item.Path)
	if trimmedPath == "" {
		return fmt.Errorf("nav path is required")
	}

	trimmedLabel := strings.TrimSpace(item.Label)
	if trimmedLabel == "" {
		return fmt.Errorf("nav label is required")
	}

	if _, exists := r.navByID[trimmedID]; exists {
		return fmt.Errorf("duplicate nav id: %s", trimmedID)
	}
	if _, exists := r.navByPath[trimmedPath]; exists {
		return fmt.Errorf("duplicate nav path: %s", trimmedPath)
	}

	r.navByID[trimmedID] = struct{}{}
	r.navByPath[trimmedPath] = struct{}{}
	r.navItems = append(r.navItems, viewmodel.NavItem{
		ID:    trimmedID,
		Label: trimmedLabel,
		Path:  trimmedPath,
		Order: item.Order,
	})

	return nil
}

// AddNavbarSlot registers a navbar slot contribution.
func (r *Registry) AddNavbarSlot(contribution viewmodel.NavbarSlotContribution) error {
	trimmedID := strings.TrimSpace(contribution.ID)
	if trimmedID == "" {
		return fmt.Errorf("navbar slot id is required")
	}

	templateName := strings.TrimSpace(contribution.Template)
	if templateName == "" {
		return fmt.Errorf("navbar slot template is required")
	}

	if contribution.Slot != viewmodel.NavbarSlotStart && contribution.Slot != viewmodel.NavbarSlotEnd {
		return fmt.Errorf("unsupported navbar slot: %s", contribution.Slot)
	}

	if _, exists := r.navbarSlotByID[trimmedID]; exists {
		return fmt.Errorf("duplicate navbar slot id: %s", trimmedID)
	}
	r.navbarSlotByID[trimmedID] = struct{}{}

	r.navbarSlots = append(r.navbarSlots, viewmodel.NavbarSlotContribution{
		ID:       trimmedID,
		Slot:     contribution.Slot,
		Template: templateName,
		Order:    contribution.Order,
		Data:     contribution.Data,
	})

	return nil
}

// AddToolbarSlot registers a toolbar slot contribution.
func (r *Registry) AddToolbarSlot(contribution viewmodel.ToolbarSlotContribution) error {
	trimmedID := strings.TrimSpace(contribution.ID)
	if trimmedID == "" {
		return fmt.Errorf("toolbar slot id is required")
	}

	templateName := strings.TrimSpace(contribution.Template)
	if templateName == "" {
		return fmt.Errorf("toolbar slot template is required")
	}

	if contribution.Slot != viewmodel.ToolbarSlotStart && contribution.Slot != viewmodel.ToolbarSlotEnd {
		return fmt.Errorf("unsupported toolbar slot: %s", contribution.Slot)
	}

	if _, exists := r.toolbarSlotByID[trimmedID]; exists {
		return fmt.Errorf("duplicate toolbar slot id: %s", trimmedID)
	}
	r.toolbarSlotByID[trimmedID] = struct{}{}

	r.toolbarSlots = append(r.toolbarSlots, viewmodel.ToolbarSlotContribution{
		ID:       trimmedID,
		Slot:     contribution.Slot,
		Template: templateName,
		Order:    contribution.Order,
		Data:     contribution.Data,
	})

	return nil
}

// HasToolbarSlot reports whether a toolbar contribution id is already registered.
func (r *Registry) HasToolbarSlot(id string) bool {
	_, exists := r.toolbarSlotByID[strings.TrimSpace(id)]
	return exists
}

// NavItems returns nav items sorted by order then label.
func (r *Registry) NavItems() []viewmodel.NavItem {
	items := make([]viewmodel.NavItem, len(r.navItems))
	copy(items, r.navItems)
	sort.Slice(items, func(left, right int) bool {
		if items[left].Order == items[right].Order {
			return items[left].Label < items[right].Label
		}
		return items[left].Order < items[right].Order
	})
	return items
}

// NavbarSlots returns slot contributions sorted by order then id.
func (r *Registry) NavbarSlots(slot viewmodel.NavbarSlot) []viewmodel.NavbarSlotContribution {
	filtered := make([]viewmodel.NavbarSlotContribution, 0, len(r.navbarSlots))
	for _, contribution := range r.navbarSlots {
		if contribution.Slot != slot {
			continue
		}
		filtered = append(filtered, contribution)
	}

	sort.Slice(filtered, func(left, right int) bool {
		if filtered[left].Order == filtered[right].Order {
			return filtered[left].ID < filtered[right].ID
		}
		return filtered[left].Order < filtered[right].Order
	})

	return filtered
}

// ToolbarSlots returns slot contributions sorted by order then id.
func (r *Registry) ToolbarSlots(slot viewmodel.ToolbarSlot) []viewmodel.ToolbarSlotContribution {
	filtered := make([]viewmodel.ToolbarSlotContribution, 0, len(r.toolbarSlots))
	for _, contribution := range r.toolbarSlots {
		if contribution.Slot != slot {
			continue
		}
		filtered = append(filtered, contribution)
	}

	sort.Slice(filtered, func(left, right int) bool {
		if filtered[left].Order == filtered[right].Order {
			return filtered[left].ID < filtered[right].ID
		}
		return filtered[left].Order < filtered[right].Order
	})

	return filtered
}
