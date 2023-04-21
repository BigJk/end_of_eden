package mainmenu

import (
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui"
	"github.com/BigJk/project_gonzo/ui/about"
	"github.com/BigJk/project_gonzo/ui/gameview"
	"github.com/BigJk/project_gonzo/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"
	"time"
)

type Model struct {
	ui.MenuBase

	zones   *zone.Manager
	choices ChoicesModel
}

func NewModel(zones *zone.Manager) Model {
	model := Model{
		zones:   zones,
		choices: NewChoicesModel(zones),
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
			Height: msg.Height - (strings.Count(ui.Title, "\n") + style.TitleStyle.GetVerticalFrameSize() + 1),
		})
		m.choices = updated.(ChoicesModel)
		return m, cmd
	}

	updated, cmd := m.choices.Update(msg)
	m.choices = updated.(ChoicesModel)

	switch m.choices.selected {
	case ChoiceContinue:
	case ChoiceNewGame:
		_ = os.Mkdir("./logs", 0777)
		f, err := os.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		m.choices = m.choices.Clear()
		return gameview.New(m, m.zones, game.NewSession(
			game.WithLogging(log.New(f, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
			lo.Ternary(os.Getenv("PG_DEBUG") == "1", game.WithDebugEnabled("127.0.0.1:8272"), nil),
		)), cmd
	case ChoiceAbout:
		m.choices = m.choices.Clear()
		return about.New(m, m.zones), cmd
	case ChoiceSettings:
	case ChoiceExit:
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) View() string {
	return lipgloss.JoinVertical(lipgloss.Top, style.TitleStyle.Render(ui.Title), m.choices.View())
}
