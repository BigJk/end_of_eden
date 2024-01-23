package merchant

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

type State string

const (
	ZoneBuyItem = "buy_item"
	ZoneLeave   = "leave"
	ZoneUpgrade = "upgrade"
	ZoneRemove  = "remove"
	ZoneBack    = "back"

	StateMain    = State("Main")
	StateUpgrade = State("Upgrade")
	StateRemove  = State("Remove")
)

type Model struct {
	ui.MenuBase

	state   State
	table   table.Model
	zones   *zone.Manager
	session *game.Session
}

func New(zones *zone.Manager, session *game.Session) Model {
	return Model{
		state:   StateMain,
		zones:   zones,
		session: session,
		table:   table.New(table.WithStyles(style.TableStyle)),
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
		m.Size = msg
	case tea.MouseMsg:
		m.LastMouse = msg

		switch m.state {
		case StateMain:
			if msg.Action == tea.MouseActionRelease && msg.Type == tea.MouseLeft {
				if m.zones.Get(ZoneBuyItem).InBounds(msg) {
					audio.Play("btn_menu")
					m = m.merchantBuy()
				} else if m.zones.Get(ZoneLeave).InBounds(msg) {
					audio.Play("btn_menu")
					m.session.LeaveMerchant()
				} else if m.zones.Get(ZoneUpgrade).InBounds(msg) {
					audio.Play("btn_menu")
					m.state = StateUpgrade
					m.table.SetCursor(0)
				} else if m.zones.Get(ZoneRemove).InBounds(msg) {
					audio.Play("btn_menu")
					m.state = StateRemove
					m.table.SetCursor(0)
				}
			}
		case StateUpgrade:
			fallthrough
		case StateRemove:
			if msg.Action == tea.MouseActionRelease && msg.Type == tea.MouseLeft {
				if m.zones.Get(ZoneBuyItem).InBounds(msg) {
					audio.Play("btn_menu")
					if m.state == StateUpgrade {
						m.session.BuyUpgradeCard(m.playerCardGetSelected())
					} else {
						m.session.BuyRemoveCard(m.playerCardGetSelected())
					}
				} else if m.zones.Get(ZoneBack).InBounds(msg) {
					audio.Play("btn_menu")
					m.state = StateMain
					m.table.SetCursor(0)
				}
			}
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			audio.Play("btn_menu")
			m = m.merchantBuy()
		}
	}

	switch m.state {
	case StateMain:
		m.table.SetColumns([]table.Column{
			{Title: "Type", Width: 10},
			{Title: "Name", Width: 10},
			{Title: "Price", Width: 10},
		})

		merchant := m.session.GetMerchant()
		m.table.SetRows(lo.Flatten([][]table.Row{
			lo.Map(merchant.Artifacts, func(guid string, index int) table.Row {
				artifact, _ := m.session.GetArtifact(guid)
				return table.Row{"Artifact", artifact.Name, fmt.Sprintf("%d$", artifact.Price)}
			}),
			lo.Map(merchant.Cards, func(guid string, index int) table.Row {
				card, _ := m.session.GetCard(guid)
				return table.Row{"Card", card.Name, fmt.Sprintf("%d$", card.Price)}
			}),
		}))
	case StateUpgrade:
		fallthrough
	case StateRemove:
		m.table.SetColumns([]table.Column{
			{Title: "Type", Width: 10},
			{Title: "Name", Width: 10},
			{Title: "Level", Width: 10},
		})

		m.table.SetRows(lo.Map(m.session.GetCards(game.PlayerActorID), func(guid string, index int) table.Row {
			card, instance := m.session.GetCard(guid)
			return table.Row{"Card", card.Name, fmt.Sprintf("%d / %d", instance.Level+1, card.MaxLevel+1)}
		}))
	}

	m.table.Focus()
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Face
	var faceSection string
	switch m.state {
	case StateMain:
		buttons := []string{
			style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneUpgrade).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2, 1, 2).
				Render(m.zones.Mark(ZoneUpgrade, fmt.Sprintf("↑  Upgrade Card (%d$)", game.DefaultUpgradeCost))),
			"",
			style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneLeave).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2).
				Render(m.zones.Mark(ZoneLeave, "Leave Merchant")),
		}

		if len(m.session.GetCards(game.PlayerActorID)) > 3 {
			buttons[1] = style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneRemove).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2, 1, 2).
				Render(m.zones.Mark(ZoneRemove, fmt.Sprintf("✕  Remove Card (%d$)", game.DefaultRemoveCost)))
		} else {
			buttons[1] = style.HeaderStyle.Copy().Background(style.BaseGrayDarker).Margin(0, 2, 1, 2).
				Render(fmt.Sprintf("✕  Remove Card (%d$)", game.DefaultRemoveCost))
		}

		faceSection = m.merchantLeft(buttons, "")
	case StateUpgrade:
		fallthrough
	case StateRemove:
		faceSection = m.merchantLeft([]string{
			style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneBack).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2).
				Render(m.zones.Mark(ZoneBack, "Back")),
		}, lo.Ternary(m.state == StateUpgrade, "What do you want to upgrade?", "What do you want to remove?"))
	}

	faceSectionWidth := lipgloss.Width(faceSection)

	help := help.New()
	help.Width = m.Size.Width - faceSectionWidth - 40 - 15 - 10
	helpText := help.ShortHelpView([]key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "move up")),
		key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "move down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "buy")),
	})

	// Middle

	var rightLook string
	switch m.state {
	case StateMain:
		m.table.SetColumns([]table.Column{
			{Title: "Type", Width: 15},
			{Title: "Name", Width: m.Size.Width - faceSectionWidth - 40 - 15 - 10},
			{Title: "Price", Width: 10},
		})
		m.table.SetWidth(m.Size.Width - faceSectionWidth - 40)
		m.table.SetHeight(ui.Min(m.Size.Height-4-10, len(m.table.Rows())+1))

		canBuy := false
		selectedItem := m.merchantGetSelected()
		var selectedItemLook string
		switch item := selectedItem.(type) {
		case *game.Artifact:
			selectedItemLook = components.ArtifactCard(m.session, item.ID, 20, 30)
			canBuy = m.session.GetPlayer().Gold >= item.Price
		case *game.Card:
			selectedItemLook = components.HalfCard(m.session, item.ID, false, 20, 20, false, 0, false)
			canBuy = m.session.GetPlayer().Gold >= item.Price
		}

		rightLook = lipgloss.JoinVertical(lipgloss.Top,
			selectedItemLook,
			style.HeaderStyle.Copy().Background(
				lo.Ternary(canBuy, lo.Ternary(m.zones.Get(ZoneBuyItem).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker), style.BaseGrayDarker),
			).Margin(1, 2).Render(m.zones.Mark(ZoneBuyItem, "Buy Item")),
		)
	case StateRemove:
		fallthrough
	case StateUpgrade:
		m.table.SetColumns([]table.Column{
			{Title: "Type", Width: 15},
			{Title: "Name", Width: m.Size.Width - faceSectionWidth - 40 - 15 - 10},
			{Title: "Level", Width: 10},
		})
		m.table.SetWidth(m.Size.Width - faceSectionWidth - 40)
		m.table.SetHeight(ui.Min(m.Size.Height-4-10, len(m.table.Rows())+1))

		selectedItem := m.playerCardGetSelected()
		var selectedItemLook string
		if len(selectedItem) > 0 {
			selectedItemLook = components.HalfCard(m.session, selectedItem, false, 20, 20, false, 0, false)
		}

		rightLook = lipgloss.JoinVertical(lipgloss.Top,
			selectedItemLook,
			style.HeaderStyle.Copy().Background(
				lo.Ternary(
					lo.Ternary(m.state == StateUpgrade, m.session.GetPlayer().Gold >= game.DefaultUpgradeCost, m.session.GetPlayer().Gold >= game.DefaultRemoveCost),
					lo.Ternary(m.zones.Get(ZoneBuyItem).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker), style.BaseGrayDarker,
				),
			).Margin(1, 2).Render(m.zones.Mark(ZoneBuyItem, lo.Ternary(m.state == StateUpgrade, fmt.Sprintf("↑  Upgrade Card (%d$)", game.DefaultUpgradeCost), fmt.Sprintf("✕  Remove Card (%d$)", game.DefaultRemoveCost)))),
		)
	}

	shopSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Top, m.table.View(), helpText), rightLook,
		),
	)

	return lipgloss.Place(m.Size.Width, m.Size.Height-5, lipgloss.Left, lipgloss.Top,
		lipgloss.JoinVertical(
			lipgloss.Top,
			style.HeaderStyle.Render("Merchant Wares"),
			lipgloss.JoinHorizontal(lipgloss.Left, faceSection, shopSection),
		),
	)
}

func (m Model) merchantLeft(buttons []string, textOverwrite string) string {
	merchant := m.session.GetMerchant()
	merchantWidth := ui.Max(lipgloss.Width(merchant.Face), 30)

	faceSection := lipgloss.JoinVertical(
		lipgloss.Top,
		append([]string{
			lipgloss.NewStyle().Margin(0, 2, 0, 2).Padding(1).Border(lipgloss.InnerHalfBlockBorder()).BorderForeground(style.BaseGray).Render(
				lipgloss.Place(merchantWidth, lipgloss.Height(merchant.Face), lipgloss.Center, lipgloss.Center, lipgloss.NewStyle().Bold(true).Foreground(style.BaseGray).Render(merchant.Face)),
			),
			lipgloss.NewStyle().
				Margin(1, 2, 2, 2).
				Padding(0, 2).
				Bold(true).Italic(true).
				Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseGray).
				Width(merchantWidth).Render(lo.Ternary(len(textOverwrite) > 0, textOverwrite, merchant.Text)),
		}, buttons...)...,
	)

	return faceSection
}

func (m Model) playerCardGetSelected() string {
	cards := m.session.GetCards(game.PlayerActorID)

	if m.table.Cursor() >= len(cards) || m.table.Cursor() < 0 {
		return ""
	}

	return cards[m.table.Cursor()]
}

func (m Model) merchantGetSelected() any {
	merchant := m.session.GetMerchant()
	items := lo.Flatten([][]any{
		lo.Map(merchant.Artifacts, func(guid string, index int) any {
			artifact, _ := m.session.GetArtifact(guid)
			return artifact
		}),
		lo.Map(merchant.Cards, func(guid string, index int) any {
			card, _ := m.session.GetCard(guid)
			return card
		}),
	})

	if m.table.Cursor() >= len(items) || m.table.Cursor() < 0 {
		return nil
	}

	return items[m.table.Cursor()]
}

func (m Model) merchantBuy() Model {
	item := m.merchantGetSelected()

	switch item := item.(type) {
	case *game.Artifact:
		if m.session.PlayerBuyArtifact(item.ID) {
			m.table.SetCursor(m.table.Cursor() - 1)
		}
	case *game.Card:
		if m.session.PlayerBuyCard(item.ID) {
			m.table.SetCursor(m.table.Cursor() - 1)
		}
	}

	return m
}
