package mainmenu

import (
	"runtime"

	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/system/localization"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

type Choice string

const (
	ChoiceWaiting    = Choice("WAITING")
	ChoiceContinue   = Choice("CONTINUE")
	ChoiceTutorial   = Choice("TUTORIAL")
	ChoiceNewGame    = Choice("NEW_GAME")
	ChoiceNewGameSOD = Choice("NEW_GAME_SOD")
	ChoiceAbout      = Choice("ABOUT")
	ChoiceSettings   = Choice("SETTINGS")
	ChoiceMods       = Choice("MODS")
	ChoiceExit       = Choice("EXIT")
)

type choiceItem struct {
	zones       *zone.Manager
	title, desc string
	key         Choice
}

func (i choiceItem) Title() string       { return i.zones.Mark("choice_"+string(i.key), i.title) }
func (i choiceItem) Description() string { return i.zones.Mark("choice_desc_"+string(i.key), i.desc) }
func (i choiceItem) FilterValue() string { return i.title }

type ChoicesModel struct {
	zones    *zone.Manager
	choices  []list.Item
	list     list.Model
	selected Choice
}

func NewChoicesModel(zones *zone.Manager, hideSettings bool) ChoicesModel {
	choices := []list.Item{
		choiceItem{zones, localization.G("ui.menu.continue", "Continue"), localization.G("ui.menu.continue_desc", "Ready to continue dying?"), ChoiceContinue},
		choiceItem{zones, localization.G("ui.menu.tutorial", "Tutorial"), localization.G("ui.menu.tutorial_desc", "Learn the basics."), ChoiceTutorial},
		choiceItem{zones, localization.G("ui.menu.new_game", "New Game"), localization.G("ui.menu.new_game_desc", "Start a new try."), ChoiceNewGame},
		choiceItem{zones, localization.G("ui.menu.new_game_sod", "New Game: Seed of the Day"), localization.G("ui.menu.new_game_sod_desc", "Start a new try with the daily seed."), ChoiceNewGameSOD},
		choiceItem{zones, localization.G("ui.menu.about", "About"), localization.G("ui.menu.about_desc", "Want to know more?"), ChoiceAbout},
		choiceItem{zones, localization.G("ui.menu.settings", "Settings"), localization.G("ui.menu.settings_desc", "Other settings won't let you survive..."), ChoiceSettings},
		choiceItem{zones, localization.G("ui.menu.mods", "Mods"), localization.G("ui.menu.mods_desc", "Make the game even more fun!"), ChoiceMods},
		choiceItem{zones, localization.G("ui.menu.exit", "Exit"), localization.G("ui.menu.exit_desc", "Got enough already?"), ChoiceExit},
	}

	// Hide exit on web
	if runtime.GOOS == "js" {
		choices = lo.Filter(choices, func(value list.Item, i int) bool {
			return value.(choiceItem).title != localization.G("ui.menu.exit", "Exit") || value.(choiceItem).title == localization.G("ui.menu.mods", "Mods")
		})
	}

	if hideSettings {
		choices = lo.Filter(choices, func(value list.Item, i int) bool {
			return value.(choiceItem).key != ChoiceSettings
		})
	}

	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	model := ChoicesModel{
		zones:    zones,
		choices:  choices,
		list:     list.New(choices, delegation, 0, 0),
		selected: ChoiceWaiting,
	}

	model.list.Title = localization.G("ui.menu.title", "Main Menu")
	model.list.SetFilteringEnabled(false)
	model.list.SetShowFilter(false)
	model.list.SetShowStatusBar(false)
	//model.list.SetShowHelp(false)
	model.list.Styles.Title = lipgloss.NewStyle().Background(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 2, 0, 2)

	return model
}

func (m ChoicesModel) Clear() ChoicesModel {
	m.selected = ChoiceWaiting
	return m
}

func (m ChoicesModel) Init() tea.Cmd {
	return nil
}

func (m ChoicesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		h, v := style.ListStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			choice, ok := m.list.SelectedItem().(choiceItem)
			if ok {
				if choice.key == ChoiceExit {
					return m, tea.Quit
				}
				m.selected = choice.key
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyUp:
			audio.Play("interface_move", -1.5)
		}
	case tea.MouseMsg:
		if (msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft) || msg.Action == tea.MouseActionMotion {
			for i := range m.choices {
				if m.zones.Get("choice_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) || m.zones.Get("choice_desc_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) {
					if m.list.Index() != i {
						audio.Play("interface_move", -1.5)
					}

					m.list.Select(i)
					if msg.Action == tea.MouseActionRelease && msg.Button == tea.MouseButtonLeft {
						m.selected = m.choices[i].(choiceItem).key
					}
					break
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m ChoicesModel) View() string {
	return style.ListStyle.Render(m.list.View())
}
