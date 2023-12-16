package about

import (
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/git"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/BigJk/end_of_eden/util"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

var (
	aboutStyle   = style.ListStyle.Copy().Align(lipgloss.Left).Padding(1, 2).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseWhite)
	versionStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseWhite).Margin(0, 2).Padding(0, 2).Foreground(style.BaseRed)
)

type Model struct {
	ui.MenuBase

	zones  *zone.Manager
	parent tea.Model
}

func New(parent tea.Model, zones *zone.Manager) Model {
	return Model{zones: zones, parent: parent}
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
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft && m.zones.Get("back").InBounds(msg) {
			audio.Play("btn_menu")

			return m.parent, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := style.TitleStyle.Render(ui.Title)

	version := versionStyle.Render(fmt.Sprintf("Version: %s (%s)", git.Tag, git.CommitHash))
	about := aboutStyle.Height(lipgloss.Height(ui.About)).Width(util.Min(m.Size.Width, 65)).Render(ui.About)
	back := m.zones.Mark("back", style.HeaderStyle.Render("Back"))

	return lipgloss.JoinVertical(lipgloss.Top, title, version, about, back)
}
