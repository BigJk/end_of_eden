package gameview

import (
	"fmt"
	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui"
	"github.com/BigJk/project_gonzo/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/harmonica"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type DamageAnimationFrame string

type DamageAnimationModel struct {
	id     string
	width  int
	height int

	sources        []game.Actor
	sourcesEnemies []*game.Enemy
	damages        []game.StateEventDamageData

	started  bool
	elapsed  float64
	finished float64

	startHp int
	reached []bool
	springs []harmonica.Spring
	px      []float64
	velx    []float64
}

func NewDamageAnimationModel(width int, height int, startHp int, sources []game.Actor, sourcesEnemies []*game.Enemy, damages []game.StateEventDamageData) DamageAnimationModel {
	return DamageAnimationModel{
		id:             fmt.Sprint(rand.Intn(100000)),
		width:          width,
		height:         height,
		sources:        sources,
		sourcesEnemies: sourcesEnemies,
		damages:        damages,
		startHp:        startHp,
		springs: lo.Map(make([]any, len(sources)), func(_ any, _ int) harmonica.Spring {
			return harmonica.NewSpring(harmonica.FPS(30), 6.0, 0.3)
		}),
		reached: make([]bool, len(sources)),
		px:      make([]float64, len(sources)),
		velx:    make([]float64, len(sources)),
	}
}

func (m DamageAnimationModel) SetSize(width int, height int) DamageAnimationModel {
	m.width = width
	m.height = height
	return m
}

func (m DamageAnimationModel) Init() tea.Cmd {
	return nil
}

func (m DamageAnimationModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m = m.SetSize(msg.Width, msg.Height)
	case tea.Key:
		if m.elapsed > 0.2 && (msg.Type == tea.KeyEnter || msg.Type == tea.KeySpace) {
			return nil, nil
		}
	case tea.MouseMsg:
		if m.elapsed > 0.2 && msg.Type == tea.MouseLeft {
			return nil, nil
		}
	case DamageAnimationFrame:
		if string(msg) != m.id {
			return m, nil
		}

		for i := 0; i < len(m.sources); i++ {
			if m.reached[i] {
				continue
			}

			m.px[i], m.velx[i] = m.springs[i].Update(m.px[i], m.velx[i], 85.0)
			if m.px[i] >= 95 {
				audio.Play("dmg1")
				m.reached[i] = true
			}

			break
		}

		elapsed := 1.0 / 30.0
		m.elapsed += elapsed

		if lo.EveryBy(m.reached, func(reached bool) bool {
			return reached
		}) && m.finished == 0 {
			m.finished = m.elapsed
		}

		if m.finished > 0 && (m.elapsed-m.finished) > 1.0 {
			return nil, nil
		}

		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return DamageAnimationFrame(m.id)
		})
	}

	// Send first tick
	if !m.started {
		m.started = true
		return m, tea.Tick(time.Second/time.Duration(30), func(t time.Time) tea.Msg {
			return DamageAnimationFrame(m.id)
		})
	}

	return m, nil
}

func (m DamageAnimationModel) View() string {
	hp := m.startHp - lo.SumBy(
		lo.Filter(m.damages, func(_ game.StateEventDamageData, i int) bool {
			return m.reached[i]
		}),
		func(item game.StateEventDamageData) int {
			return item.Damage
		})

	facePlayer := lipgloss.JoinVertical(
		lipgloss.Center,
		faceStyle.Copy().BorderForeground(style.BaseGray).Foreground(style.BaseRed).Render(
			lipgloss.JoinHorizontal(lipgloss.Center, lo.Map([]rune(fmt.Sprint(hp)), func(item rune, index int) string {
				digit, _ := strconv.ParseInt(string(item), 10, 64)
				return ui.Numbers[lo.Clamp(digit, 0, 9)]
			})...),
		),
	)

	playerWidth := lipgloss.Width(facePlayer)
	playerHeight := lipgloss.Width(facePlayer)
	side := m.width/2 - playerWidth/2 - 5

	middle := lipgloss.JoinVertical(lipgloss.Center,
		lipgloss.NewStyle().Foreground(style.BaseGrayDarker).Render(strings.Repeat("│\n", (m.height-playerHeight)/2)),
		facePlayer,
		lipgloss.NewStyle().Foreground(style.BaseGrayDarker).Render(strings.Repeat("│\n", (m.height-playerHeight)/2)),
	)

	faceEnemies := lo.Map(m.sources, func(item game.Actor, i int) string {
		face := lipgloss.JoinVertical(
			lipgloss.Right,
			"\n\n",
			faceStyle.Copy().BorderForeground(style.BaseRed).Foreground(lipgloss.Color(m.sourcesEnemies[i].Color)).Render(m.sourcesEnemies[i].Look),
			style.BaseText.Render(m.sourcesEnemies[i].Name),
		)

		width := lipgloss.Width(face)

		return lipgloss.NewStyle().Width(side).Render(lipgloss.NewStyle().Padding(0, 0, 0, lo.Clamp(int((float64(side)-float64(width))*(m.px[i]/100.0)), 0, side-width-2)).Render(face))
	})

	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center,
		lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.JoinVertical(lipgloss.Left, faceEnemies...),
			middle,
			lipgloss.NewStyle().Width(side-2).MaxWidth(side-2).Render(lipgloss.JoinHorizontal(lipgloss.Left, lo.Map(m.damages, func(dmg game.StateEventDamageData, i int) string {
				if !m.reached[i] {
					return ""
				}

				digits := lipgloss.JoinHorizontal(lipgloss.Center, lo.Map([]rune(fmt.Sprint(dmg.Damage)), func(item rune, index int) string {
					digit, _ := strconv.ParseInt(string(item), 10, 64)
					return ui.Numbers[lo.Clamp(digit, 0, 9)]
				})...) + "\n"

				return lipgloss.NewStyle().Margin(0, 2, 0, 4).Foreground(style.BaseRedDarker).Render(digits)
			})...)),
		),
	)
}
