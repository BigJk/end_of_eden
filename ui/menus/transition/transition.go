package transition

import (
	"github.com/BigJk/end_of_eden/system/gen"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
	"time"
)

const (
	charFull      = "●"
	charHalfLeft  = "◐"
	charHalfRight = "◑"
	charEmpty     = "○"
)

type TickMsg string

type Model struct {
	ui.MenuBase

	parent  tea.Model
	created time.Time
	line    string
}

func New(parent tea.Model) Model {
	return Model{parent: parent, line: gen.GetRandom("loading_lines")}
}

func (m Model) Start() tea.Msg {
	m.created = time.Now()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case root.ModelGettingVisibleMsg:
		m.created = time.Now()
	case TickMsg:
		if !m.created.IsZero() && time.Since(m.created) > time.Millisecond*800 {
			return m.parent, tea.Sequence(root.GettingVisible(), func() tea.Msg {
				return tea.WindowSizeMsg{Width: m.Size.Width, Height: m.Size.Height}
			})
		}
	}

	return m, tea.Tick(time.Second/30, func(t time.Time) tea.Msg {
		return TickMsg("tick")
	})
}

func (m Model) View() string {
	elapsed := time.Since(m.created)

	spinner := strings.Split(strings.Repeat(charEmpty, 10), "")
	pos := elapsed.Milliseconds() / 80 % int64(len(spinner))
	spinner[pos] = charFull
	spinner[(pos+1)%int64(len(spinner))] = charHalfLeft
	spinner[((pos-1)+int64(len(spinner)))%int64(len(spinner))] = charHalfRight

	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, m.line+"\n", style.RedText.Render(strings.Join(spinner, ""))))
}
