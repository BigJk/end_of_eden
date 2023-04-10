package mainmenu

import (
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/menus"
	"github.com/BigJk/project_gonzo/menus/about"
	"github.com/BigJk/project_gonzo/menus/gameview"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

type Model struct {
	menus.MenuBase

	choices ChoicesModel
}

func NewModel() Model {
	model := Model{
		choices: NewChoicesModel(),
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
		updated, cmd := m.choices.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - (strings.Count(menus.Title, "\n") + menus.TitleStyle.GetVerticalFrameSize() + 1),
		})
		m.choices = updated.(ChoicesModel)
		return m, cmd
	}

	updated, cmd := m.choices.Update(msg)
	m.choices = updated.(ChoicesModel)

	switch m.choices.selected {
	case ChoiceContinue:
	case ChoiceNewGame:
		m.choices = m.choices.Clear()
		return gameview.New(m, game.NewSession()), cmd
	case ChoiceAbout:
		m.choices = m.choices.Clear()
		return about.New(m), cmd
	case ChoiceSettings:
	}

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, menus.TitleStyle.Render(menus.Title), m.choices.View())
}
