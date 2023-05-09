package root

import (
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/menus/lua_error"
	"github.com/BigJk/end_of_eden/ui/overlay"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

type PushModelMsg tea.Model

func Push(model tea.Model) tea.Cmd {
	return func() tea.Msg {
		return PushModelMsg(model)
	}
}

type ToolTip struct {
	ID      string
	Content string
	X       int
	Y       int
}

type ToolTipMsg ToolTip

func TooltipCreate(tip ToolTip) tea.Cmd {
	return func() tea.Msg {
		return ToolTipMsg(tip)
	}
}

type ToolTipDeleteMsg string

func TooltipDelete(id string) tea.Cmd {
	return func() tea.Msg {
		return ToolTipDeleteMsg(id)
	}
}

type Model struct {
	zones    *zone.Manager
	stack    []tea.Model
	size     tea.WindowSizeMsg
	tooltips map[string]ToolTip
}

func New(zones *zone.Manager, root tea.Model) Model {
	return Model{
		zones:    zones,
		stack:    []tea.Model{root},
		tooltips: map[string]ToolTip{},
	}
}

func (m Model) PushModel(model tea.Model) Model {
	m.stack = append(m.stack, model)
	m.tooltips = map[string]ToolTip{}
	return m
}

func (m Model) SetRoot(model tea.Model) Model {
	if len(m.stack) == 0 {
		return m.PushModel(model)
	}

	m.stack[0] = model
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case ToolTipMsg:
		m.tooltips[msg.ID] = ToolTip(msg)
	case ToolTipDeleteMsg:
		delete(m.tooltips, string(msg))
	case PushModelMsg:
		m = m.PushModel(msg)
	}

	curIndex := len(m.stack) - 1

	var cmd tea.Cmd
	m.stack[curIndex], cmd = m.stack[curIndex].Update(msg)

	if menu, ok := m.stack[curIndex].(ui.Menu); ok && !menu.HasSize() {
		return m, tea.Batch(cmd, func() tea.Msg {
			return m.size
		})
	}

	if m.stack[curIndex] == nil {
		m.stack = m.stack[:len(m.stack)-1]
	}

	return m, cmd
}

func (m Model) View() string {
	if len(m.stack) == 0 {
		return "stack empty!"
	}

	view := m.zones.Scan(m.stack[len(m.stack)-1].View())

	for _, v := range m.tooltips {
		view = overlay.PlaceOverlay(v.X, v.Y, v.Content, view)
	}

	return view
}

func CheckLuaErrors(zones *zone.Manager, s *game.Session) tea.Cmd {
	var errors []game.LuaError

	errChan := s.LuaErrors()
	if len(errChan) == 0 {
		return nil
	}

	for r := range errChan {
		errors = append(errors, r)

		if len(errChan) == 0 {
			break
		}
	}

	return tea.Sequence(lo.Map(errors, func(item game.LuaError, index int) tea.Cmd {
		return func() tea.Msg {
			return PushModelMsg(lua_error.New(zones, item))
		}
	})...)
}
