package mainmenu

import (
	"github.com/BigJk/project_gonzo/menus"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

func (i choiceItem) Title() string       { return i.title }
func (i choiceItem) Description() string { return i.desc }
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
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(menus.BaseRed).BorderForeground(menus.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(menus.BaseRedDarker).BorderForeground(menus.BaseRed)

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
	model.list.Styles.Title = lipgloss.NewStyle().Background(menus.BaseRedDarker).Foreground(menus.BaseWhite).Padding(0, 2, 0, 2)

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
		h, v := menus.ListStyle.GetFrameSize()
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
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m ChoicesModel) View() string {
	return menus.ListStyle.Render(m.list.View())
}
