package about

import (
	"fmt"
	"github.com/BigJk/end_of_eden/internal/git"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/system/localization"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
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
		if (msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft) && m.zones.Get("back").InBounds(msg) {
			audio.Play("btn_menu")

			return m.parent, nil
		}
	}

	return m, nil
}

func (m Model) View() string {
	title := style.TitleStyle.Render(ui.Title)

	aboutText := ui.AboutText()
	version := versionStyle.Render(fmt.Sprintf(localization.G("ui.about.version_fmt", "Version: %s (%s)"), git.Tag, git.CommitHash))
	about := aboutStyle.Height(lipgloss.Height(aboutText)).Width(ui.Min(m.Size.Width, 65)).Render(aboutText)
	back := m.zones.Mark("back", style.HeaderStyle.Render(localization.G("ui.about.back", "Back")))

	return lipgloss.JoinVertical(lipgloss.Top, title, version, about, back)
}
