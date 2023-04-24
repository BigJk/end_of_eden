package merchant

import (
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/components"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/BigJk/end_of_eden/util"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
)

const (
	ZoneBuyItem = "buy_item"
	ZoneLeave   = "leave"
	ZoneUpgrade = "upgrade"
	ZoneRemove  = "remove"
)

type Model struct {
	ui.MenuBase

	table   table.Model
	zones   *zone.Manager
	session *game.Session
}

func New(zones *zone.Manager, session *game.Session) Model {
	return Model{
		zones:   zones,
		session: session,
		table: table.New(table.WithStyles(style.TableStyle), table.WithColumns([]table.Column{
			{Title: "Type", Width: 10},
			{Title: "Name", Width: 10},
			{Title: "Price", Width: 10},
		})),
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

		if msg.Type == tea.MouseLeft {
			if m.zones.Get(ZoneBuyItem).InBounds(msg) {
				m = m.merchantBuy()
			} else if m.zones.Get(ZoneLeave).InBounds(msg) {
				m.session.LeaveMerchant()
			}
		}
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m = m.merchantBuy()
		}
	}

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

	m.table.Focus()
	m.table, cmd = m.table.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	// Face
	merchant := m.session.GetMerchant()
	merchantWidth := util.Max(lipgloss.Width(merchant.Face), 30)

	faceSection := lipgloss.JoinVertical(
		lipgloss.Top,
		lipgloss.NewStyle().Margin(0, 2, 0, 2).Padding(1).Border(lipgloss.InnerHalfBlockBorder()).BorderForeground(style.BaseGray).Render(
			lipgloss.Place(merchantWidth, lipgloss.Height(merchant.Face), lipgloss.Center, lipgloss.Center, lipgloss.NewStyle().Bold(true).Foreground(style.BaseGray).Render(merchant.Face)),
		),
		lipgloss.NewStyle().
			Margin(1, 2, 2, 2).
			Padding(0, 2).
			Bold(true).Italic(true).
			Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(style.BaseGray).
			Width(merchantWidth).Render(merchant.Text),
		style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneUpgrade).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2, 1, 2).
			Render(m.zones.Mark(ZoneUpgrade, fmt.Sprintf("↑  Upgrade Card (%d$)", game.DefaultUpgradeCost))),
		style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneRemove).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2, 1, 2).
			Render(m.zones.Mark(ZoneRemove, fmt.Sprintf("✕  Remove Card (%d$)", game.DefaultRemoveCost))),
		style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneLeave).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Margin(0, 2).
			Render(m.zones.Mark(ZoneLeave, "Leave Merchant")),
	)
	faceSectionWidth := lipgloss.Width(faceSection)

	// Wares
	m.table.SetColumns([]table.Column{
		{Title: "Type", Width: 15},
		{Title: "Name", Width: m.Size.Width - faceSectionWidth - 40 - 15 - 10},
		{Title: "Price", Width: 10},
	})
	m.table.SetWidth(m.Size.Width - faceSectionWidth - 40)
	m.table.SetHeight(util.Min(m.Size.Height-4-10, len(m.table.Rows())+1))

	canBuy := false
	selectedItem := m.merchantGetSelected()
	var selectedItemLook string
	switch item := selectedItem.(type) {
	case *game.Artifact:
		selectedItemLook = components.ArtifactCard(m.session, item.ID, 20, 20)
		canBuy = m.session.GetPlayer().Gold >= item.Price
	case *game.Card:
		selectedItemLook = components.HalfCard(m.session, item.ID, false, 20, 20, false)
		canBuy = m.session.GetPlayer().Gold >= item.Price
	}

	help := help.New()
	help.Width = m.Size.Width - faceSectionWidth - 40 - 15 - 10
	helpText := help.ShortHelpView([]key.Binding{
		key.NewBinding(key.WithKeys("up"), key.WithHelp("↑", "move up")),
		key.NewBinding(key.WithKeys("down"), key.WithHelp("↓", "move down")),
		key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "buy")),
	})

	shopSection := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Left,
			lipgloss.JoinVertical(lipgloss.Top, m.table.View(), helpText),
			lipgloss.JoinVertical(lipgloss.Top,
				selectedItemLook,
				style.HeaderStyle.Copy().Background(
					lo.Ternary(canBuy, lo.Ternary(m.zones.Get(ZoneBuyItem).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker), style.BaseGrayDarker),
				).Margin(1, 2).Render(m.zones.Mark(ZoneBuyItem, "Buy Item")),
			),
		),
	)

	return lipgloss.JoinVertical(
		lipgloss.Top,
		style.HeaderStyle.Render("Merchant Wares"),
		lipgloss.JoinHorizontal(lipgloss.Left, faceSection, shopSection),
	)
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

	if m.table.Cursor() >= len(items) {
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
