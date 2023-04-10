package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"os"
)

const TriggerBug = true

type TestModel struct {
	size         tea.WindowSizeMsg
	selectedCard int
	numberCards  int
}

func (m TestModel) Init() tea.Cmd {
	return nil
}

func (m TestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		m.size = msg
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft || msg.Type == tea.MouseMotion {
			// Check zones to set selected card
			for i := 0; i < m.numberCards; i++ {
				if zone.Get(fmt.Sprintf("%s%d", "card_", i)).InBounds(msg) {
					m.selectedCard = i
				}
			}
		}
	}

	return m, nil
}

func (m TestModel) View() string {
	cardStyle := lipgloss.NewStyle().Width(30) //.Padding(1, 2).Margin(0, 2)

	var cardBoxes []string
	for i := 0; i < m.numberCards; i++ {
		selected := i == m.selectedCard

		style := cardStyle.
			//Border(lipgloss.NormalBorder(), selected, false, false, false).
			//BorderBackground(lipgloss.Color("#cccccc")).
			Background(lipgloss.Color("#cccccc"))
		//BorderForeground(lipgloss.Color("#ffffff")).
		//Foreground(lipgloss.Color("#ffffff"))

		// If the card is selected we give it a bit more height
		if selected {
			cardBoxes = append(cardBoxes,
				style.
					Height(5).
					Render(""),
			)
			continue
		}

		// Non-selected card style
		cardBoxes = append(cardBoxes,
			style.
				Height(5).
				Render(""),
		)
	}

	for i := range cardBoxes {
		cardBoxes[i] = zone.Mark(fmt.Sprintf("%s%d", "card_", i), cardBoxes[i])
	}

	// This works:
	if !TriggerBug {
		return zone.Scan(lipgloss.JoinHorizontal(lipgloss.Top, cardBoxes...))
	}

	// Freeze:
	return zone.Scan(lipgloss.Place(m.size.Width, m.size.Height, lipgloss.Center, lipgloss.Bottom, lipgloss.JoinHorizontal(lipgloss.Bottom, cardBoxes...)))
}

func main() {
	zone.NewGlobal()

	p := tea.NewProgram(TestModel{
		numberCards: 3,
	}, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}
