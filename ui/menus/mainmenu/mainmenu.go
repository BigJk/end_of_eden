package mainmenu

import (
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/menus/about"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/style"
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
		audio.Play("btn_menu")

		if saved, err := os.ReadFile("./session.save"); err == nil {
			f, err := os.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}

			session := game.NewSession(
				game.WithLogging(log.New(f, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
				lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
			)

			err = session.GobDecode(saved)
			if err != nil {
				log.Println("Error loading save:", err)
			} else {
				m.choices = m.choices.Clear()
				return gameview.New(m, m.zones, session), cmd
			}
		}

	case ChoiceNewGame:
		audio.Play("btn_menu")

		_ = os.Mkdir("./logs", 0777)
		f, err := os.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		m.choices = m.choices.Clear()
		return gameview.New(m, m.zones, game.NewSession(
			game.WithLogging(log.New(f, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
			game.WithMods(settings.LoadedSettings.Mods),
			lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
		)), cmd
	case ChoiceAbout:
		audio.Play("btn_menu")

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
