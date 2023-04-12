package mainmenu

import (
	"github.com/BigJk/project_gonzo/menus/style"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
)

type Choice string

const (
	ChoiceWaiting  = Choice("WAITING")
	ChoiceContinue = Choice("CONTINUE")
	ChoiceNewGame  = Choice("NEW_GAME")
	ChoiceAbout    = Choice("ABOUT")
	ChoiceSettings = Choice("SETTINGS")
	ChoiceExit     = Choice("EXIT")
)

type choiceItem struct {
	title, desc string
	key         Choice
}

func (i choiceItem) Title() string       { return zone.Mark("choice_"+string(i.key), i.title) }
func (i choiceItem) Description() string { return zone.Mark("choice_desc_"+string(i.key), i.desc) }
func (i choiceItem) FilterValue() string { return i.title }

type ChoicesModel struct {
	choices  []list.Item
	list     list.Model
	selected Choice
}

func NewChoicesModel() ChoicesModel {
	choices := []list.Item{
		choiceItem{"Continue", "Ready to continue dying?", ChoiceContinue},
		choiceItem{"New Game", "Start a new try.", ChoiceNewGame},
		choiceItem{"About", "Want to know more?", ChoiceAbout},
		choiceItem{"Settings", "Other settings won't let you survive...", ChoiceSettings},
		choiceItem{"Exit", "Got enough already?", ChoiceExit},
	}

	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	model := ChoicesModel{
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
		switch keypress := msg.String(); keypress {
		case "enter":
			choice, ok := m.list.SelectedItem().(choiceItem)
			if ok {
				if choice.key == ChoiceExit {
					return m, tea.Quit
				}
				m.selected = choice.key
			}
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft || msg.Type == tea.MouseMotion {
			for i := range m.choices {
				if zone.Get("choice_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) || zone.Get("choice_desc_"+string(m.choices[i].(choiceItem).key)).InBounds(msg) {
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
