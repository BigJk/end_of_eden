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

// Model is the root model of the game. It holds the current model stack and
// the zone manager. The top model of the internal stack is the current model
// and will  be rendered.
type Model struct {
	zones           *zone.Manager
	stack           []tea.Model
	size            tea.WindowSizeMsg
	tooltips        map[string]Tooltip
	transitionModel func(parent tea.Model) tea.Model
}

// New creates a new root model.
func New(zones *zone.Manager, root tea.Model) Model {
	return Model{
		zones:    zones,
		stack:    []tea.Model{root},
		tooltips: map[string]Tooltip{},
	}
}

// PushModel pushes a new model on the stack.
func (m Model) PushModel(model tea.Model) Model {
	if m.transitionModel != nil {
		m.stack = append(m.stack, m.transitionModel(model))
	} else {
		m.stack = append(m.stack, model)
	}
	m.tooltips = map[string]Tooltip{}
	return m
}

// SetRoot sets the root (last) model of the stack.
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
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case TooltipMsg:
		m.tooltips[msg.ID] = Tooltip(msg)
	case TooltipDeleteMsg:
		delete(m.tooltips, string(msg))
	case TooltipClearMsg:
		m.tooltips = map[string]Tooltip{}
	case PushModelMsg:
		for _, model := range msg {
			m = m.PushModel(model)
		}
		cmds = append(cmds, GettingVisible())
	case PushTransitionFuncMsg:
		m.transitionModel = msg
	}

	curIndex := len(m.stack) - 1

	var cmd tea.Cmd
	m.stack[curIndex], cmd = m.stack[curIndex].Update(msg)

	if cmd != nil {
		cmds = append(cmds, cmd)
	}

	if m.stack[curIndex] == nil {
		// If we remove the top model, we need to send a window size message to the new top model
		// to avoid the layout to be broken.
		cmds = append(cmds,
			func() tea.Msg {
				return tea.WindowSizeMsg{
					Width:  m.size.Width,
					Height: m.size.Height,
				}
			},
			GettingVisible(),
		)
		m.stack = m.stack[:len(m.stack)-1]
	} else if menu, ok := m.stack[curIndex].(ui.Menu); ok && !menu.HasSize() {
		cmds = append(cmds, func() tea.Msg {
			return tea.WindowSizeMsg{
				Width:  m.size.Width,
				Height: m.size.Height,
			}
		})
	}

	return m, tea.Batch(cmds...)
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

// CheckLuaErrors checks if there are any lua errors and pushes them to the stack as lua error model.
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
			return PushModelMsg([]tea.Model{lua_error.New(zones, item)})
		}
	})...)
}
