package mainmenu

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/system/audio"
	image2 "github.com/BigJk/end_of_eden/system/image"
	"github.com/BigJk/end_of_eden/system/settings"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components/loader"
	"github.com/BigJk/end_of_eden/ui/menus/about"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/menus/intro"
	"github.com/BigJk/end_of_eden/ui/menus/mods"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/root"
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

	image    string
	settings settings.Settings
	zones    *zone.Manager
	choices  ChoicesModel

	settingValues []uiset.Value
	settingSaver  uiset.Saver
	didLoad       bool
}

func NewModel(zones *zone.Manager, settings settings.Settings, values []uiset.Value, saver uiset.Saver) Model {
	img, _ := image2.Fetch("title.jpg", image2.WithResize(180, 9))

	audio.PlayMusic("planet_mining")

	model := Model{
		image:         img,
		zones:         zones,
		settings:      settings,
		choices:       NewChoicesModel(zones, len(values) == 0 || saver == nil),
		settingSaver:  saver,
		settingValues: values,
	}

	return model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if !m.didLoad {
		m.didLoad = true

		if m.settings.GetBool("experimental") {
			l, done, _ := loader.New(intro.New(m), "Initial loading")
			go func() {
				_, _ = image2.FetchAnimation("intro.gif", image2.WithMaxWidth(100))
				done <- true
			}()
			return l, nil
		}
	}

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
			f, err := fs.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
			if err != nil {
				panic(err)
			}

			session := game.NewSession(
				game.WithLogging(log.New(f, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
				lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
			)
			image2.ResetSearchPaths()
			image2.AddSearchPaths(lo.Map(session.GetLoadedMods(), func(item string, index int) string {
				return fmt.Sprintf("./mods/%s/images/", item)
			})...)

			err = session.GobDecode(saved)
			if err != nil {
				log.Println("Error loading save:", err)
			} else {
				m.choices = m.choices.Clear()
				return m, tea.Sequence(
					cmd,
					root.Push(gameview.New(m, m.zones, session)),
				)
			}
		}

	case ChoiceNewGame:
		audio.Play("btn_menu")

		_ = os.Mkdir("./logs", 0777)
		f, err := fs.OpenFile("./logs/S "+strings.ReplaceAll(time.Now().Format(time.DateTime), ":", "-")+".txt", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		image2.ResetSearchPaths()
		image2.AddSearchPaths(lo.Map(m.settings.GetStrings("mods"), func(item string, index int) string {
			return fmt.Sprintf("./mods/%s/images/", item)
		})...)

		m.choices = m.choices.Clear()
		return m, tea.Sequence(
			cmd,
			root.Push(gameview.New(m, m.zones, game.NewSession(
				game.WithLogging(log.New(f, "SESSION ", log.Ldate|log.Ltime|log.Lshortfile)),
				game.WithMods(m.settings.GetStrings("mods")),
				lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil),
			))),
		)
	case ChoiceAbout:
		audio.Play("btn_menu")

		m.choices = m.choices.Clear()
		return about.New(m, m.zones), cmd
	case ChoiceMods:
		audio.Play("btn_menu")

		m.choices = m.choices.Clear()
		return m, root.Push(mods.NewModel(m.zones, m.settings))
	case ChoiceSettings:
		if m.settingSaver != nil {
			audio.Play("btn_menu")

			m.choices = m.choices.Clear()
			return m, root.Push(uiset.NewModel(m.zones, m.settingValues, m.settingSaver))
		}

		// TODO: don't show settings item if no settings saver is set
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
