package mainmenu

import (
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

type Choice string

const (
	ChoiceWaiting  = Choice("WAITING")
	ChoiceContinue = Choice("CONTINUE")
	ChoiceNewGame  = Choice("NEW_GAME")
	ChoiceAbout    = Choice("ABOUT")
	ChoiceSettings = Choice("SETTINGS")
	ChoiceMods     = Choice("MODS")
	ChoiceExit     = Choice("EXIT")
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
		choiceItem{zones, "Continue", "Ready to continue dying?", ChoiceContinue},
		choiceItem{zones, "New Game", "Start a new try.", ChoiceNewGame},
		choiceItem{zones, "About", "Want to know more?", ChoiceAbout},
		choiceItem{zones, "Settings", "Other settings won't let you survive...", ChoiceSettings},
		choiceItem{zones, "Mods", "Make the game even more fun!", ChoiceMods},
		choiceItem{zones, "Exit", "Got enough already?", ChoiceExit},
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

	model.list.Title = "Main Menu"
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
		if msg.Type == tea.MouseLeft || msg.Type == tea.MouseMotion {
			for i := range m.choices {
				if m.zones.Get("choice_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) || m.zones.Get("choice_desc_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) {
					if m.list.Index() != i {
						audio.Play("interface_move", -1.5)
					}

					m.list.Select(i)
					if msg.Type == tea.MouseLeft {
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
