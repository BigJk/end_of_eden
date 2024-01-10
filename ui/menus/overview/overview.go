package overview

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
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
	ZoneChoices   = "choices_"
	ZoneCards     = "cards_"
	ZoneArtifacts = "artifacts_"

	ChoiceCharacter = Choice("CHARACTER")
	ChoiceLogs      = Choice("LOGS")
	ChoiceArtifacts = Choice("ARTIFACTS")
	ChoiceCards     = Choice("CARDS")
	ChoiceQuit      = Choice("QUIT")
)

type choiceItem struct {
	zones       *zone.Manager
	title, desc string
	key         Choice
}

func (i choiceItem) Title() string       { return i.zones.Mark(ZoneChoices+string(i.key), i.title) }
func (i choiceItem) Description() string { return i.zones.Mark(ZoneChoices+string(i.key), i.desc) }
func (i choiceItem) FilterValue() string { return i.zones.Mark(ZoneChoices+string(i.key), i.title) }

// MenuModel is responsible for showing information about the player, cards,
// artifacts and the option to leave.
type MenuModel struct {
	ui.MenuBase

	zones     *zone.Manager
	parent    tea.Model
	choices   []list.Item
	list      list.Model
	selected  Choice
	listFocus bool

	artifactTable table.Model
	cardTable     table.Model

	logsViewport    viewport.Model
	logsInitialized bool

	Session *game.Session
}

func New(parent tea.Model, zones *zone.Manager, session *game.Session) MenuModel {
	choices := []list.Item{
		choiceItem{zones, "Character", "Check your stats.", ChoiceCharacter},
		choiceItem{zones, "Logs", "Check what happened.", ChoiceLogs},
		choiceItem{zones, "Artifacts", "Inspect your artifacts.", ChoiceArtifacts},
		choiceItem{zones, "Cards", "Inspect your cards.", ChoiceCards},
		choiceItem{zones, "Quit", "Return to menu.", ChoiceQuit},
	}

	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	model := MenuModel{
		zones:     zones,
		parent:    parent,
		choices:   choices,
		list:      list.New(choices, delegation, 0, 0),
		selected:  ChoiceCharacter,
		listFocus: true,
		cardTable: table.New(
			table.WithStyles(style.TableStyle),
			table.WithColumns([]table.Column{
				{Title: "Name", Width: 25},
				{Title: "Tags", Width: 20},
				{Title: "Level", Width: 5},
			}),
		),
		artifactTable: table.New(
			table.WithStyles(style.TableStyle),
			table.WithColumns([]table.Column{
				{Title: "Name", Width: 25},
				{Title: "Tags", Width: 20},
				{Title: "Price", Width: 5},
			}),
		),
		Session: session,
	}

	model.list.Title = "Overview"
	model.list.SetFilteringEnabled(false)
	model.list.SetShowFilter(false)
	model.list.SetShowStatusBar(false)
	model.list.DisableQuitKeybindings()
	model.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys("tab"),
				key.WithHelp("tab", "switch view"),
			),
		}
	}
	model.list.Styles.Title = lipgloss.NewStyle().Background(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 2, 0, 2)

	return model
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cards := m.Session.GetCards(game.PlayerActorID)
	artifacts := m.Session.GetArtifacts(game.PlayerActorID)

	// Update card table
	m.cardTable.SetRows(lo.Map(cards, func(guid string, index int) table.Row {
		card, instance := m.Session.GetCard(guid)
		return table.Row{m.zones.Mark(ZoneCards+fmt.Sprint(index), card.Name), strings.Join(card.Tags, ", "), fmt.Sprint(instance.Level)}
	}))
	m.cardTable.SetHeight(m.Size.Height - style.HeaderStyle.GetVerticalFrameSize() - 1 - 2)

	// Update artifact table
	m.artifactTable.SetRows(lo.Map(artifacts, func(guid string, index int) table.Row {
		art, _ := m.Session.GetArtifact(guid)
		return table.Row{m.zones.Mark(ZoneArtifacts+fmt.Sprint(index), art.Name), strings.Join(art.Tags, ", "), fmt.Sprintf("%d$", art.Price)}
	}))
	m.artifactTable.SetHeight(m.Size.Height - style.HeaderStyle.GetVerticalFrameSize() - 1 - 2)

	// Messages
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg

		h, v := style.ListStyle.GetFrameSize()
		m.list.SetSize(m.listWidth()-h, msg.Height-v)

		m = m.updateLogViewport()

		m.cardTable.SetWidth(m.contentWidth() - 55 - 4)
		m.cardTable.SetColumns([]table.Column{
			{Title: "Name", Width: m.contentWidth() - 55 - 4 - 10 - 20 - 4},
			{Title: "Tags", Width: 20},
			{Title: "Level", Width: 10},
		})

		m.artifactTable.SetWidth(m.contentWidth() - 55 - 4)
		m.artifactTable.SetColumns([]table.Column{
			{Title: "Name", Width: m.contentWidth() - 55 - 4 - 10 - 20 - 4},
			{Title: "Tags", Width: 20},
			{Title: "Price", Width: 10},
		})
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEscape:
			if !m.listFocus {
				m.listFocus = true
			} else {
				return m.parent, nil
			}
		case tea.KeyTab:
			m.listFocus = !m.listFocus
		case tea.KeyEnter:
			if m.list.SelectedItem().(choiceItem).key == ChoiceQuit {
				m.Session.Close()
				return nil, nil
			}
		case tea.KeyDown:
			fallthrough
		case tea.KeyUp:
			audio.Play("interface_move", -1.5)
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			for i := range m.choices {
				if m.zones.Get(ZoneChoices + string(m.choices[i].(choiceItem).key)).InBounds(msg) {
					m.list.Select(i)
					m.listFocus = false
					break
				}
			}

			switch m.choices[m.list.Index()].(choiceItem).key {
			case ChoiceCharacter:
			case ChoiceLogs:
			case ChoiceArtifacts:
			case ChoiceCards:
				for i := range cards {
					if m.zones.Get(ZoneCards + fmt.Sprint(cards[i])).InBounds(msg) {
						m.cardTable.SetCursor(i)
						break
					}
				}
			}
		}
	}

	// List focused
	if m.listFocus {
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)

		return m, cmd
	}

	// Content focused
	switch m.choices[m.list.Index()].(choiceItem).key {
	case ChoiceCharacter:
	case ChoiceLogs:
		var cmd tea.Cmd
		m.logsViewport, cmd = m.logsViewport.Update(msg)
		return m, cmd
	case ChoiceArtifacts:
		var cmd tea.Cmd
		m.artifactTable.Focus()
		m.artifactTable, cmd = m.artifactTable.Update(msg)
		return m, cmd
	case ChoiceCards:
		var cmd tea.Cmd
		m.cardTable.Focus()
		m.cardTable, cmd = m.cardTable.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m MenuModel) View() string {
	var contentBox string

	contentStyle := styleMenuContent.Width(m.contentWidth()).Height(m.Size.Height)
	if !m.listFocus {
		contentStyle = contentStyle.Copy().BorderForeground(style.BaseRedDarker)
	}

	switch m.choices[m.list.Index()].(choiceItem).key {
	case ChoiceCharacter:
		player := m.Session.GetPlayer()

		status := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(style.BaseWhite).Render("Status Effects:") + "\n\n" + strings.Join(lo.Map(player.StatusEffects.ToSlice(), func(guid string, index int) string {
			return components.StatusEffect(m.Session, guid) + ": " + m.Session.GetStatusEffectState(guid)
		}), "\n\n")

		contentBox = contentStyle.Render(
			lipgloss.JoinVertical(
				lipgloss.Top,
				style.HeaderStyle.Render("Character"),
				lipgloss.NewStyle().Margin(0, 0, 0, 2).Render(status),
			))
	case ChoiceLogs:
		contentBox = contentStyle.Render(lipgloss.JoinVertical(
			lipgloss.Top,
			style.HeaderStyle.Render("Logs"),
			m.logsViewport.View(),
		))
	case ChoiceArtifacts:
		var selected string
		if m.artifactTable.Cursor() < len(m.Session.GetArtifacts(game.PlayerActorID)) {
			selected = components.ArtifactCard(m.Session, m.Session.GetArtifacts(game.PlayerActorID)[m.artifactTable.Cursor()], 20, 40, 45)
		}

		contentBox = contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			style.HeaderStyle.Render("Artifacts"),
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Margin(0, 2).Render(m.artifactTable.View()),
				selected,
			),
		))
	case ChoiceCards:
		var selected string
		if m.artifactTable.Cursor() < len(m.Session.GetCards(game.PlayerActorID)) {
			selected = components.HalfCard(m.Session, m.Session.GetCards(game.PlayerActorID)[m.cardTable.Cursor()], false, 20, 40, false, 45)
		}

		contentBox = contentStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
			style.HeaderStyle.Render("Cards"),
			lipgloss.JoinHorizontal(lipgloss.Left,
				lipgloss.NewStyle().Margin(0, 2).Render(m.cardTable.View()),
				selected,
			),
		))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, "\n"+m.list.View(), contentBox)
}

func (m MenuModel) listWidth() int {
	return lipgloss.Width(m.list.View())
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
