package eventview

import (
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/BigJk/end_of_eden/util"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/wordwrap"
	"github.com/samber/lo"
	"strconv"
	"strings"
)

const (
	ZoneChoice = "choice_"
)

var (
	titleStyle          = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 1)
	infoStyle           = lipgloss.NewStyle().BorderStyle(lipgloss.ThickBorder()).BorderForeground(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 1)
	choiceStyle         = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.ThickBorder(), true).BorderForeground(style.BaseGrayDarker).Foreground(style.BaseWhite)
	choiceSelectedStyle = choiceStyle.Copy().BorderForeground(style.BaseRed).Foreground(style.BaseWhite)
)

type Model struct {
	ui.MenuBase

	zones          *zone.Manager
	session        *game.Session
	viewport       viewport.Model
	curEvent       string
	selectedChoice int
}

func New(zones *zone.Manager, session *game.Session) Model {
	return Model{
		zones:   zones,
		session: session,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if cmd := root.CheckLuaErrors(m.zones, m.session); cmd != nil {
		return m, cmd
	}

	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.eventUpdateSize(msg.Width, msg.Height, !m.HasSize())
		m.Size = msg
	case tea.KeyMsg:
		if val, err := strconv.Atoi(msg.String()); err == nil {
			m.selectedChoice = val - 1
		}

		switch msg.Type {
		case tea.KeyEnter:
			m = m.tryFinishEvent()
		case tea.KeyTab:
			m.selectedChoice = (m.selectedChoice + 1) % len(m.session.GetEvent().Choices)
		}
	case tea.MouseMsg:
		m.LastMouse = msg

		if msg.Type == tea.MouseLeft || msg.Type == tea.MouseMotion {
			if m.session.GetEvent() != nil {
				for i := 0; i < len(m.session.GetEvent().Choices); i++ {
					if choiceZone := m.zones.Get(fmt.Sprintf("%s%d", ZoneChoice, i)); choiceZone.InBounds(msg) {
						if msg.Type == tea.MouseLeft && m.selectedChoice == i {
							audio.Play("button")
							m = m.tryFinishEvent()
							break
						} else {
							m.selectedChoice = i
						}
					}
				}
			}
		}
	}

	if m.HasSize() {
		m = m.eventUpdateSize(m.Size.Width, m.Size.Height, false)
		m = m.eventUpdateContent()
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	return fmt.Sprintf("%s\n%s\n%s\n%s", m.eventHeaderView(), m.viewport.View(), m.eventFooterView(), strings.Join(m.eventChoices(), "\n"))
}

func (m Model) tryFinishEvent() Model {
	if len(m.session.GetEvent().Choices) == 0 || m.selectedChoice < len(m.session.GetEvent().Choices) {
		m.session.FinishEvent(m.selectedChoice)
		return m.eventUpdateContent()
	}
	return m
}

func (m Model) eventUpdateSize(width, height int, init bool) Model {
	headerHeight := lipgloss.Height(m.eventHeaderView())
	footerHeight := lipgloss.Height(m.eventFooterView())
	verticalMarginHeight := headerHeight + footerHeight + m.eventChoiceHeight()

	if init {
		m.viewport = viewport.New(util.Min(width, 100), height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = false

		m = m.eventUpdateContent()
	} else {
		m.viewport.Width = util.Min(width, 100)
		m.viewport.Height = height - verticalMarginHeight
	}

	return m
}

func (m Model) eventUpdateContent() Model {
	if m.session.GetEvent() == nil {
		m.viewport.SetContent("")
		return m
	}

	// Don't update if we are still in the same event.
	eventId := m.session.GetEvent().ID
	if m.curEvent == eventId {
		return m
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(glamour.DarkStyleConfig),
		glamour.WithWordWrap(m.viewport.Width),
	)
	res, _ := r.Render(m.session.GetEvent().Description)

	m.viewport.SetContent(res)
	m.curEvent = eventId
	return m
}

func (m Model) eventHeaderView() string {
	if m.session.GetEvent() == nil {
		return ""
	}

	title := titleStyle.Render(m.session.GetEvent().Name)
	line := style.GrayTextDarker.Render(strings.Repeat("━", util.Max(0, m.viewport.Width-lipgloss.Width(title))))
	return "\n" + lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) eventFooterView() string {
	if m.session.GetEvent() == nil {
		return ""
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := style.GrayTextDarker.Render(strings.Repeat("━", util.Max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

func (m Model) eventChoices() []string {
	if m.session.GetEvent() == nil {
		return nil
	}

	choices := lo.Map(m.session.GetEvent().Choices, func(item game.EventChoice, index int) string {
		if m.selectedChoice == index {
			return choiceSelectedStyle.Width(util.Min(m.Size.Width, 100)).Render(wordwrap.String(fmt.Sprintf("%d. %s", index+1, item.Description), util.Min(m.Size.Width, 100-choiceStyle.GetHorizontalFrameSize())))
		}
		return choiceStyle.Width(util.Min(m.Size.Width, 100)).Render(wordwrap.String(fmt.Sprintf("%d. %s", index+1, item.Description), util.Min(m.Size.Width, 100-choiceStyle.GetHorizontalFrameSize())))
	})

	return lo.Map(choices, func(item string, index int) string {
		return m.zones.Mark(fmt.Sprintf("%s%d", ZoneChoice, index), item)
	})
}

func (m Model) eventChoiceHeight() int {
	return lo.SumBy(m.eventChoices(), func(item string) int {
		return lipgloss.Height(item) + 1
	})
}