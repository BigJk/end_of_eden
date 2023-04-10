package about

import (
	"github.com/BigJk/project_gonzo/menus"
	"github.com/BigJk/project_gonzo/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	menus.MenuBase

	parent tea.Model
}

func New(parent tea.Model) Model {
	return Model{parent: parent}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case tea.KeyMsg:
		if msg.Type == tea.KeyEscape {
			return m.parent, nil
		}
	}

	return m, nil
}

var aboutStyle = menus.ListStyle.Copy().
	Align(lipgloss.Left).
	Padding(1, 2).
	Border(lipgloss.NormalBorder(), false, false, false, true).
	BorderForeground(menus.BaseWhite)

func (m Model) View() string {
	title := menus.TitleStyle.Render(menus.Title)

	version := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(menus.BaseWhite).
		Margin(0, 2).
		Padding(0, 2).
		Foreground(menus.BaseRed).
		Render("Version: 0.0.1 alpha")

	about := aboutStyle.Height(lipgloss.Height(menus.About)).Width(util.Min(m.Size.Width, 65)).Render(menus.About)

	back := lipgloss.NewStyle().Margin(0, 2).Foreground(menus.BaseRed).Render("<- ESC")

	return lipgloss.JoinVertical(lipgloss.Top, title, version, about, back)
}
