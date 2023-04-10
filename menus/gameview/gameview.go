package gameview

import (
	"fmt"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/menus"
	"github.com/BigJk/project_gonzo/util"
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

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).BorderForeground(menus.BaseWhite).Foreground(menus.BaseWhite).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

const (
	ZoneCard        = "card_"
	ZoneEnemy       = "enemy_"
	ZoneEventChoice = "event_choice_"
)

type Model struct {
	menus.MenuBase

	parent              tea.Model
	viewport            viewport.Model
	selectedChoice      int
	selectedCard        int
	selectedOpponent    int
	inOpponentSelection bool

	Session *game.Session
}

func New(parent tea.Model, session *game.Session) Model {
	session.Log(game.LogTypeSuccess, "Game started! Good luck...")

	return Model{
		parent:  parent,
		Session: session,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	//
	// Keyboard
	//
	case tea.KeyMsg:
		if val, err := strconv.Atoi(msg.String()); err == nil {
			m.selectedChoice = val - 1
		}

		switch msg.Type {
		case tea.KeyEnter:
			switch m.Session.GetGameState() {
			// If we are in an event commit the choice. Only commit if choice is in range.
			case game.GameStateEvent:
				m = m.tryFinishEvent()
			// Cast a card
			case game.GameStateFight:
				if m.selectedCard >= len(m.Session.GetFight().Hand) {
					m.selectedCard = 0
				}

				m = m.tryCast()
			}
		case tea.KeyEscape:
			if m.inOpponentSelection {
				m.inOpponentSelection = false
			} else {
				return NewMenuModel(m, m.Session), nil
			}
		case tea.KeyTab:
			switch m.Session.GetGameState() {
			// Select a choice
			case game.GameStateEvent:
				if len(m.Session.GetEvent().Choices) > 0 {
					m.selectedChoice = (m.selectedChoice + 1) % len(m.Session.GetEvent().Choices)
				}
			// Select a card or opponent
			case game.GameStateFight:
				if len(m.Session.GetFight().Hand) > 0 {
					if m.inOpponentSelection {
						m.selectedOpponent = (m.selectedOpponent + 1) % m.Session.GetOpponentCount(game.PlayerActorID)
					} else {
						m.selectedCard = (m.selectedCard + 1) % len(m.Session.GetFight().Hand)
					}
				}
			}
		case tea.KeyLeft:
		case tea.KeyRight:
			// TODO: right / left movement
		}
	//
	// Mouse
	//
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft || msg.Type == tea.MouseMotion {
			switch m.Session.GetGameState() {
			case game.GameStateEvent:
				for i := 0; i < len(m.Session.GetEvent().Choices); i++ {
					if choiceZone := zone.Get(fmt.Sprintf("%s%d", ZoneEventChoice, i)); choiceZone.InBounds(msg) {
						if msg.Type == tea.MouseLeft && m.selectedChoice == i {
							m = m.tryFinishEvent()
						} else {
							m.selectedChoice = i
						}
					}
				}
			case game.GameStateFight:

				if m.inOpponentSelection {
					for i := 0; i < m.Session.GetOpponentCount(game.PlayerActorID); i++ {
						if cardZone := zone.Get(fmt.Sprintf("%s%d", ZoneEnemy, i)); cardZone.InBounds(msg) {
							if msg.Type == tea.MouseLeft && m.selectedOpponent == i {
								m = m.tryCast()
							} else {
								m.selectedOpponent = i
							}
						}
					}
				} else {
					for i := 0; i < len(m.Session.GetFight().Hand); i++ {
						if cardZone := zone.Get(fmt.Sprintf("%s%d", ZoneCard, i)); cardZone.InBounds(msg) {
							if msg.Type == tea.MouseLeft && m.selectedCard == i {
								m = m.tryCast()
							} else {
								m.selectedCard = i
							}
						}
					}
				}
			}
		}
	//
	// Window Size
	//
	case tea.WindowSizeMsg:
		headerHeight := lipgloss.Height(m.eventHeaderView())
		footerHeight := lipgloss.Height(m.eventFooterView())
		verticalMarginHeight := headerHeight + footerHeight + m.eventChoiceHeight()

		if !m.HasSize() {
			m.viewport = viewport.New(util.Min(msg.Width, 100), msg.Height-verticalMarginHeight)
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = false

			m = m.eventUpdateContent()
		} else {
			m.viewport.Width = util.Min(msg.Width, 100)
			m.viewport.Height = msg.Height - verticalMarginHeight
		}

		m.Size = msg
	}

	//
	// Updating
	//

	switch m.Session.GetGameState() {
	case game.GameStateFight:
	case game.GameStateMerchant:
	case game.GameStateEvent:
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if !m.HasSize() {
		return "..."
	}

	switch m.Session.GetGameState() {
	case game.GameStateFight:
		return lipgloss.JoinVertical(
			lipgloss.Top,
			m.fightStatusTop(),
			lipgloss.NewStyle().Width(m.Size.Width).Height(m.fightViewHeight()).Render(m.fightEnemyView()),
			lipgloss.NewStyle().Foreground(menus.BaseGrayDarker).Render(strings.Repeat("─", m.Size.Width)),
			lipgloss.NewStyle().Width(m.Size.Width).Height(m.fightViewHeight()).Render(m.fightCardView()),
			m.fightStatusBottom(),
		)
	case game.GameStateMerchant:
	case game.GameStateEvent:
		return lipgloss.Place(
			m.Size.Width,
			m.Size.Height,
			lipgloss.Center,
			lipgloss.Center,
			fmt.Sprintf("%s\n%s\n%s\n%s", m.eventHeaderView(), m.viewport.View(), m.eventFooterView(), strings.Join(m.eventChoices(), "\n")),
		)
	}

	return ""
}

//
// Actions
//

func (m Model) tryCast() Model {
	if len(m.Session.GetFight().Hand) > 0 {
		card, _ := m.Session.GetCard(m.Session.GetFight().Hand[m.selectedCard])
		if card.NeedTarget {
			if m.inOpponentSelection {
				m.inOpponentSelection = false
				m.Session.PlayerCastHand(m.selectedCard, m.Session.GetOpponentByIndex(game.PlayerActorID, m.selectedOpponent).ID)
			} else {
				m.inOpponentSelection = true
			}
		} else {
			m.Session.PlayerCastHand(m.selectedCard, "")
		}
	}

	return m
}

func (m Model) tryFinishEvent() Model {
	if len(m.Session.GetEvent().Choices) == 0 || m.selectedChoice < len(m.Session.GetEvent().Choices) {
		m.Session.FinishEvent(m.selectedChoice)
		return m.eventUpdateContent()
	}
	return m
}

//
// Fight View
//

func (m Model) fightStatusTop() string {
	style := lipgloss.NewStyle().Width(m.Size.Width).Height(1).Background(menus.BaseGrayDarker).Foreground(menus.BaseWhite).AlignHorizontal(lipgloss.Center)

	if m.inOpponentSelection {
		return style.Render("Select a target for your card...")
	}

	return style.Render(m.Session.GetFight().Description)
}

func (m Model) fightStatusBottom() string {
	lastLog, _ := lo.Last(m.Session.Logs)

	return lipgloss.NewStyle().Height(1).Foreground(menus.BaseWhite).Background(menus.BaseGrayDarker).Width(m.Size.Width).Render(
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.NewStyle().Width(m.Size.Width/2).Render(fmt.Sprintf(" Deck: %d / Used: %d / Exhausted: %d", len(m.Session.GetFight().Deck), len(m.Session.GetFight().Used), len(m.Session.GetFight().Exhausted))),
			lipgloss.NewStyle().Width(m.Size.Width/2).Align(lipgloss.Right).Render(styleLogs[lastLog.Type].Render(lastLog.Message+" ")),
		),
	)
}

func (m Model) fightViewHeight() int {
	return (m.Size.Height / 2) - 2
}

func (m Model) fightEnemyView() string {
	var enemyBoxes []string

	c := m.Session.GetOpponentCount(game.PlayerActorID)
	for i := 0; i < c; i++ {
		enemy := m.Session.GetOpponentByIndex(game.PlayerActorID, i)
		if enemy == nil {
			continue
		}

		faceStyle := lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(menus.BaseGrayDarker).Foreground(menus.BaseRed)
		if m.inOpponentSelection && i == m.selectedOpponent {
			faceStyle.BorderForeground(menus.BaseWhite)
		}
		face := faceStyle.Render("@")
		enemyBoxes = append(enemyBoxes, zone.Mark(fmt.Sprintf("%s%d", ZoneEnemy, i), lipgloss.NewStyle().Foreground(menus.BaseWhite).Margin(0, 2).Render(lipgloss.JoinVertical(lipgloss.Center, face, enemy.Name, fmt.Sprintf("%d / %d", enemy.HP, enemy.MaxHP)))))
	}

	return lipgloss.Place(m.Size.Width, m.fightViewHeight(), lipgloss.Center, lipgloss.Center, lipgloss.JoinHorizontal(lipgloss.Center, enemyBoxes...))
}

func (m Model) fightCardView() string {
	cardStyle := lipgloss.NewStyle().Width(30).Padding(1, 2).Margin(0, 2)

	var cardBoxes = lo.Map(m.Session.GetFight().Hand, func(guid string, index int) string {
		card, _ := m.Session.GetCard(guid)
		selected := index == m.selectedCard
		selected = false

		style := cardStyle.Border(lipgloss.NormalBorder(), selected, false, false, false).BorderBackground(lipgloss.Color(card.Color)).Background(lipgloss.Color(card.Color)).BorderForeground(menus.BaseWhite).Foreground(menus.BaseWhite)
		if selected {
			return style.Height(util.Min(m.fightViewHeight()-1, m.fightViewHeight()/2+5)).Render(wordwrap.String(fmt.Sprintf("%s\n\n%s\n\n%s", strings.Repeat("•", card.PointCost), menus.BoldStyle.Render(card.Name), card.Description), 20))
		}
		return style.Height(m.fightViewHeight() / 2).Render(wordwrap.String(fmt.Sprintf("%s\n\n%s\n\n%s", strings.Repeat("•", card.PointCost), menus.BoldStyle.Render(card.Name), card.Description), 20))
	})

	cardBoxes = lo.Map(cardBoxes, func(item string, i int) string {
		return zone.Mark(fmt.Sprintf("%s%d", ZoneCard, i), item)
	})

	return lipgloss.Place(m.Size.Width, m.fightViewHeight(), lipgloss.Center, lipgloss.Bottom, lipgloss.JoinHorizontal(lipgloss.Bottom, cardBoxes...))
}

//
// Event View
//

func (m Model) eventUpdateContent() Model {
	if m.Session.GetEvent() == nil {
		m.viewport.SetContent("")
		return m
	}

	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(glamour.DarkStyleConfig),
		glamour.WithWordWrap(m.viewport.Width),
	)
	res, _ := r.Render(m.Session.GetEvent().Description)

	m.viewport.SetContent(res)
	return m
}

func (m Model) eventHeaderView() string {
	if m.Session.GetEvent() == nil {
		return ""
	}

	title := titleStyle.Render(m.Session.GetEvent().Name)
	line := menus.BaseText.Render(strings.Repeat("─", util.Max(0, m.viewport.Width-lipgloss.Width(title))))
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m Model) eventFooterView() string {
	if m.Session.GetEvent() == nil {
		return ""
	}

	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := menus.BaseText.Render(strings.Repeat("─", util.Max(0, m.viewport.Width-lipgloss.Width(info))))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}

var choiceStyle = lipgloss.NewStyle().Padding(0, 1).Border(lipgloss.RoundedBorder(), true).BorderForeground(menus.BaseWhite).Foreground(menus.BaseWhite)
var choiceSelectedStyle = choiceStyle.Copy().BorderForeground(menus.BaseRed).Foreground(menus.BaseWhite)

func (m Model) eventChoices() []string {
	if m.Session.GetEvent() == nil {
		return nil
	}

	choices := lo.Map(m.Session.GetEvent().Choices, func(item game.EventChoice, index int) string {
		if m.selectedChoice == index {
			return choiceSelectedStyle.Width(util.Min(m.Size.Width, 100)).Render(wordwrap.String(fmt.Sprintf("%d. %s", index+1, item.Description), util.Min(m.Size.Width, 100-choiceStyle.GetHorizontalFrameSize())))
		}
		return choiceStyle.Width(util.Min(m.Size.Width, 100)).Render(wordwrap.String(fmt.Sprintf("%d. %s", index+1, item.Description), util.Min(m.Size.Width, 100-choiceStyle.GetHorizontalFrameSize())))
	})

	return lo.Map(choices, func(item string, index int) string {
		return zone.Mark(fmt.Sprintf("%s%d", ZoneEventChoice, index), item)
	})
}

func (m Model) eventChoiceHeight() int {
	return lo.SumBy(m.eventChoices(), func(item string) int {
		return lipgloss.Height(item)
	})
}
