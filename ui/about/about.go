package about

import (
	"github.com/BigJk/project_gonzo/ui"
	"github.com/BigJk/project_gonzo/ui/style"
	"github.com/BigJk/project_gonzo/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	aboutStyle      = style.ListStyle.Copy().Align(lipgloss.Left).Padding(1, 2).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseWhite)
	versionStyle    = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseWhite).Margin(0, 2).Padding(0, 2).Foreground(style.BaseRed)
	backButtonStyle = lipgloss.NewStyle().Margin(0, 2).Foreground(style.BaseRed)
)

type Model struct {
	ui.MenuBase

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

func (m Model) View() string {
	title := style.TitleStyle.Render(ui.Title)

	version := versionStyle.Render("Version: 0.0.1 alpha")
	about := aboutStyle.Height(lipgloss.Height(ui.About)).Width(util.Min(m.Size.Width, 65)).Render(ui.About)
	back := backButtonStyle.Render("<- ESC")

	return lipgloss.JoinVertical(lipgloss.Top, title, version, about, back)
}
