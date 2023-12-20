package settings

import (
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/localization"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"strconv"
)

type item struct {
	val Value
}

func (i item) Title() string {
	var val string
	switch i.val.Type {
	case Bool:
		if i.val.Val.(bool) {
			val = lipgloss.NewStyle().Foreground(style.BaseGreen).Render(localization.G("basics.on", "on"))
		} else {
			val = lipgloss.NewStyle().Foreground(style.BaseRed).Render(localization.G("basics.off", "off"))
		}
	case Int:
		val = fmt.Sprint(i.val.Val)
	case Float:
		val = fmt.Sprintf("%.2f", i.val.Val)
	case String:
		val = fmt.Sprint(i.val.Val)
	}

	return fmt.Sprintf("%-20s", localization.G(fmt.Sprintf("settings.%s.title", i.val.Key), i.val.Name)) + " : " + val
}

func (i item) Description() string {
	return localization.G(fmt.Sprintf("settings.%s.description", i.val.Key), i.val.Description)
}

func (i item) FilterValue() string {
	return localization.G(fmt.Sprintf("settings.%s.title", i.val.Key), i.val.Name)
}

type Model struct {
	ui.MenuBase

	editInput textinput.Model
	editError string
	editValue int
	values    []Value
	saver     Saver
	list      list.Model
	zones     *zone.Manager
}

func NewModel(zones *zone.Manager, values []Value, saver Saver) Model {
	return Model{
		zones:     zones,
		saver:     saver,
		values:    values,
		editInput: textinput.New(),
		editValue: -1,
	}.setup()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() != list.Filtering {
			switch msg.String() {
			case "q":
				fallthrough
			case "esc":
				if m.editValue >= 0 {
					m.editValue = -1
					return m, nil
				}
				return nil, nil
			}
		}

		switch msg.Type {
		case tea.KeyDown:
			fallthrough
		case tea.KeyUp:
			audio.Play("interface_move", -1.5)
		case tea.KeyEnter:
			if m.editValue >= 0 {
				switch m.values[m.editValue].Type {
				case Int:
					res, err := strconv.Atoi(m.editInput.Value())
					if err != nil {
						m.editError = "Invalid number"
						return m, nil
					}

					if m.values[m.editValue].Min != nil && res < m.values[m.editValue].Min.(int) {
						m.editError = "Number too small (min: " + fmt.Sprint(m.values[m.editValue].Min) + ")"
						return m, nil
					}

					if m.values[m.editValue].Max != nil && res > m.values[m.editValue].Max.(int) {
						m.editError = "Number too big (max: " + fmt.Sprint(m.values[m.editValue].Max) + ")"
						return m, nil
					}

					m.values[m.editValue].Val = res
					m.editValue = -1
					return m.saveValues(), nil
				case Float:
					res, err := strconv.ParseFloat(m.editInput.Value(), 64)
					if err != nil {
						m.editError = "Invalid number"
						return m, nil
					}

					if m.values[m.editValue].Min != nil && res < m.values[m.editValue].Min.(float64) {
						m.editError = "Number too small (min: " + fmt.Sprint(m.values[m.editValue].Min) + ")"
						return m, nil
					}

					if m.values[m.editValue].Max != nil && res > m.values[m.editValue].Max.(float64) {
						m.editError = "Number too big (max: " + fmt.Sprint(m.values[m.editValue].Max) + ")"
						return m, nil
					}

					m.values[m.editValue].Val = res
					m.editValue = -1
					return m.saveValues(), nil
				case String:
					m.values[m.editValue].Val = m.editInput.Value()
					m.editValue = -1
					return m.saveValues(), nil
				}
			} else if m.list.Cursor() >= 0 && m.list.Cursor() < len(m.values) {
				switch m.values[m.list.Cursor()].Type {
				case Bool:
					m.values[m.list.Cursor()].Val = !m.values[m.list.Cursor()].Val.(bool)
					return m.saveValues(), nil
				default:
					m.editValue = m.list.Cursor()

					m.editInput.CharLimit = 200
					m.editInput.Prompt = "> "
					m.editInput.Placeholder = fmt.Sprint(m.values[m.editValue].Val)
					m.editInput.PlaceholderStyle = m.editInput.PlaceholderStyle.Copy().Foreground(style.BaseGrayDarker)
					m.editInput.Cursor.TextStyle = lipgloss.NewStyle().Foreground(style.BaseGrayDarker).Background(style.BaseWhite)
					m.editInput.Cursor.Blink = false

					m.editInput.SetValue(fmt.Sprint(m.values[m.editValue].Val))
					m.editInput.CursorEnd()

					return m, m.editInput.Focus()
				}
			}
		}
	case tea.WindowSizeMsg:
		m.Size = msg
		m.list.SetSize(msg.Width-4, msg.Height-2)
	}

	if m.editValue >= 0 {
		var cmd tea.Cmd
		m.editInput, cmd = m.editInput.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.editValue >= 0 && m.editValue < len(m.values) {
		return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Padding(1, 3).Border(lipgloss.ThickBorder(), true).Render(lipgloss.JoinVertical(
				lipgloss.Top,
				fmt.Sprintf("Enter new value for %s:\n", style.BoldStyle.Render(m.values[m.editValue].Name)),
				m.editInput.View(),
				lo.Ternary(m.editError != "", style.RedText.Render("\nError: "+m.editError), ""),
				style.GrayText.Render("\nPress 'esc' to cancel, 'enter' to save"),
			)),
			lipgloss.WithWhitespaceChars("?"),
			lipgloss.WithWhitespaceForeground(style.BaseGrayDarker),
		)
	}

	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

func (m Model) items() []list.Item {
	items := make([]list.Item, 0)
	for _, v := range m.values {
		items = append(items, item{v})
	}
	return items
}

func (m Model) setup() Model {
	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	m.list = list.New(m.items(), delegation, 0, 0)
	m.list.Title = "Settings"
	m.list.SetFilteringEnabled(true)
	m.list.SetShowFilter(true)
	m.list.SetShowStatusBar(false)
	m.list.KeyMap.Filter = key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter"),
	)
	m.list.KeyMap.NextPage = key.NewBinding(
		key.WithKeys("right", "l", "pgdown", "d"),
		key.WithHelp("→/l/pgdn", "next page"),
	)
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "edit value")),
		}
	}
	m.list.AdditionalFullHelpKeys = m.list.AdditionalShortHelpKeys
	m.list.Styles.Title = lipgloss.NewStyle().Background(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 2, 0, 2)

	return m
}

func (m Model) saveValues() Model {
	_ = m.saver(m.values)
	m.list.SetItems(m.items())
	return m
}
