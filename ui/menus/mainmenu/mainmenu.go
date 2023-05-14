package mainmenu

import (
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/fs"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/image"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/menus/about"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/menus/mods"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/BigJk/imeji"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
)

type Model struct {
	ui.MenuBase

	image   string
	zones   *zone.Manager
	choices ChoicesModel
}

func NewModel(zones *zone.Manager) Model {
	img, _ := image.Fetch("title.png", imeji.WithResize(180, 9))

	model := Model{
		image:   img,
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
			Height: msg.Height - lipgloss.Height(m.image) - 1,
		})
		m.choices = updated.(ChoicesModel)
		return m, cmd
	}

	updated, cmd := m.choices.Update(msg)
	m.choices = updated.(ChoicesModel)

	switch m.choices.selected {
	case ChoiceContinue:
		audio.Play("btn_menu")

		if saved, err := fs.ReadFile("./session.save"); err == nil {
			/*f, err := os.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}*/

			session := game.NewSession(
				game.WithLogging(log.New(io.Discard, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
				lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
			)
			image.ResetSearchPaths()
			image.AddSearchPaths(lo.Map(session.GetLoadedMods(), func(item string, index int) string {
				return fmt.Sprintf("./mods/%s/images/", item)
			})...)

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

		/*_ = os.Mkdir("./logs", 0777)
		f, err := os.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}*/

		image.ResetSearchPaths()
		image.AddSearchPaths(lo.Map(settings.LoadedSettings.Mods, func(item string, index int) string {
			return fmt.Sprintf("./mods/%s/images/", item)
		})...)

		m.choices = m.choices.Clear()
		return m, root.Push(gameview.New(m, m.zones, game.NewSession(
			game.WithLogging(log.New(io.Discard, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
			game.WithMods(settings.LoadedSettings.Mods),
			lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
		)))
	case ChoiceAbout:
		audio.Play("btn_menu")

		m.choices = m.choices.Clear()
		return about.New(m, m.zones), cmd
	case ChoiceMods:
		audio.Play("btn_menu")

		m.choices = m.choices.Clear()
		return m, root.Push(mods.NewModel(m.zones))
	case ChoiceSettings:
	case ChoiceExit:
		return m, tea.Quit
	}

	return m, cmd
}

func (m Model) View() string {
	titleImage := lipgloss.NewStyle().
		Border(lipgloss.InnerHalfBlockBorder(), false, false, true, false).
		BorderForeground(style.BaseRedDarker).
		Margin(0, 0, 0, 2).
		Render("\n" + lipgloss.NewStyle().MaxWidth(m.Size.Width-3).MaxHeight(lipgloss.Height(m.image)-1).Render(m.image))

	return lipgloss.JoinVertical(lipgloss.Top,
		titleImage,
		m.choices.View())
}
