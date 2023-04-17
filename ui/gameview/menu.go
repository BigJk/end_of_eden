package gameview

import (
	"fmt"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui"
	"github.com/BigJk/project_gonzo/ui/style"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/samber/lo"
	"strings"
	"time"
)

var (
	styleMenuContent = lipgloss.NewStyle().Margin(0, 0, 0, 2).Padding(0, 1).Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseGrayDarker)
	styleLogs        = map[game.LogType]lipgloss.Style{
		game.LogTypeInfo:    lipgloss.NewStyle(),
		game.LogTypeWarning: lipgloss.NewStyle().Foreground(style.BaseYellow),
		game.LogTypeDanger:  lipgloss.NewStyle().Foreground(style.BaseRed),
		game.LogTypeSuccess: lipgloss.NewStyle().Foreground(style.BaseGreen),
	}
)

type Choice string

const (
	ChoiceCharacter = Choice("CHARACTER")
	ChoiceLogs      = Choice("LOGS")
	ChoiceArtifacts = Choice("ARTIFACTS")
	ChoiceCards     = Choice("CARDS")
)

type choiceItem struct {
	title, desc string
	key         Choice
}

func (i choiceItem) Title() string       { return i.title }
func (i choiceItem) Description() string { return i.desc }
func (i choiceItem) FilterValue() string { return i.title }

// MenuModel is responsible for showing information about the player, cards,
// artifacts and the option to leave.
type MenuModel struct {
	ui.MenuBase

	parent    tea.Model
	choices   []list.Item
	list      list.Model
	selected  Choice
	listFocus bool

	logsViewport    viewport.Model
	logsInitialized bool

	Session *game.Session
}

func NewMenuModel(parent tea.Model, session *game.Session) MenuModel {
	choices := []list.Item{
		choiceItem{"Character", "Check your stats.", ChoiceCharacter},
		choiceItem{"Logs", "Check what happened.", ChoiceLogs},
		choiceItem{"Artifacts", "Inspect your artifacts.", ChoiceArtifacts},
		choiceItem{"Cards", "Inspect your cards.", ChoiceCards},
	}

	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	model := MenuModel{
		parent:    parent,
		choices:   choices,
		list:      list.New(choices, delegation, 0, 0),
		selected:  ChoiceCharacter,
		listFocus: true,
		Session:   session,
	}

	model.list.Title = "Overview"
	model.list.SetFilteringEnabled(false)
	model.list.SetShowFilter(false)
	model.list.SetShowStatusBar(false)
	model.list.SetShowHelp(false)
	model.list.Styles.Title = lipgloss.NewStyle().Background(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 2, 0, 2)

	return model
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg

		h, v := style.ListStyle.GetFrameSize()
		m.list.SetSize(m.Size.Width/4-h, msg.Height-v)

		m = m.updateLogViewport()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			return m.parent, nil
		case tea.KeyTab:
			m.listFocus = !m.listFocus
		}
	}

	if m.listFocus {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	} else {
		switch m.choices[m.list.Index()].(choiceItem).key {
		case ChoiceCharacter:
		case ChoiceLogs:
			var cmd tea.Cmd
			m.logsViewport, cmd = m.logsViewport.Update(msg)
			return m, cmd
		case ChoiceArtifacts:
		case ChoiceCards:
		}
	}

	return m, nil
}

func (m MenuModel) View() string {
	var contentBox string

	contentStyle := styleMenuContent.Width(m.Size.Width - m.Size.Width/4 - style.ListStyle.GetHorizontalFrameSize()).Height(m.Size.Height)
	if !m.listFocus {
		contentStyle = contentStyle.Copy().BorderForeground(style.BaseGray)
	}

	switch m.choices[m.list.Index()].(choiceItem).key {
	case ChoiceCharacter:
		contentBox = contentStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				style.HeaderStyle.Render("Character"),
			))
	case ChoiceLogs:
		contentBox = contentStyle.Render(lipgloss.JoinVertical(
			lipgloss.Top,
			style.HeaderStyle.Render("Logs"),
			m.logsViewport.View(),
		))
	case ChoiceArtifacts:
		contentBox = contentStyle.Render(style.HeaderStyle.Render("Artifacts"))
	case ChoiceCards:
		contentBox = contentStyle.Render(style.HeaderStyle.Render("Cards"))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, "\n"+m.list.View(), contentBox)
}

func (m MenuModel) listWidth() int {
	return m.Size.Width / 4
}

func (m MenuModel) contentWidth() int {
	return m.Size.Width - m.listWidth() - style.ListStyle.GetHorizontalFrameSize()
}

func (m MenuModel) updateLogViewport() MenuModel {
	headerHeight := style.HeaderStyle.GetVerticalFrameSize() + 1
	verticalMarginHeight := headerHeight

	if !m.logsInitialized {
		m.logsViewport = viewport.New(m.contentWidth(), m.Size.Height-verticalMarginHeight-1)
		m.logsViewport.YPosition = headerHeight
		m.logsViewport.HighPerformanceRendering = false
		m.logsViewport.SetContent(strings.Join(lo.Map(lo.Reverse(m.Session.Logs), func(item game.LogEntry, index int) string {
			return wordwrap.String(
				fmt.Sprintf("  %s |- %s %s %s",
					lipgloss.NewStyle().Foreground(style.BaseGray).Render(fmt.Sprintf("#%05d", index)),
					lipgloss.NewStyle().Foreground(style.BaseGray).Render(item.Time.Format(time.RFC822)),
					styleLogs[item.Type].Render(fmt.Sprintf(" [ %-8s ]", item.Type)),
					item.Message,
				),
				m.contentWidth()-style.ListStyle.GetHorizontalFrameSize(),
			)
		}), "\n"))
		m.logsInitialized = true
	} else {
		m.logsViewport.Width = m.contentWidth()
		m.logsViewport.Height = m.Size.Height - verticalMarginHeight
	}

	return m
}
