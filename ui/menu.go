package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

// Menu is a tea.Model that keeps track of its size. It is intended to
// trigger the re-distribution of the tea.WindowSizeMsg for nested models.
// It also contains some common functionality for menus.
type Menu interface {
	tea.Model
	HasSize() bool
}

// MenuBase is the base Menu implementation.
type MenuBase struct {
	Size      tea.WindowSizeMsg
	LastMouse tea.MouseMsg
	Zones     *zone.Manager
}

// NewMenuBase returns a new MenuBase.
func NewMenuBase() MenuBase {
	return MenuBase{}
}

// WithSize returns a new MenuBase with the given size.
func (m MenuBase) WithSize(size tea.WindowSizeMsg) MenuBase {
	m.Size = size
	return m
}

// WithZones returns a new MenuBase with the given zone manager.
func (m MenuBase) WithZones(zones *zone.Manager) MenuBase {
	m.Zones = zones
	return m
}

// HasSize returns true if the menu has a saved size.
func (m MenuBase) HasSize() bool {
	return m.Size.Width > 0
}

// ZoneBackground returns the background color for a zone.
func (m MenuBase) ZoneBackground(zone string, hover lipgloss.Color, normal lipgloss.Color) lipgloss.Color {
	if m.Zones == nil {
		return normal
	}

	if m.Zones.Get(zone).InBounds(m.LastMouse) {
		return hover
	}

	return normal
}

// ZoneInBounds returns true if the last mouse position is in the bounds of the given zone.
func (m MenuBase) ZoneInBounds(zone string) bool {
	if m.Zones == nil {
		return false
	}
	return m.Zones.Get(zone).InBounds(m.LastMouse)
}

// ZoneMark returns the given string marked with the given zone.
func (m MenuBase) ZoneMark(zone string, str string) string {
	if m.Zones == nil {
		return str
	}
	return m.Zones.Mark(zone, str)
}
