package loader

import (
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"math"
	"strings"
	"time"
)

type LoaderTick string

var loaderFrames = strings.Split("◐◓◑◒", "")

type Model struct {
	ui.MenuBase

	parent tea.Model
	text   chan string
	done   chan bool

	currentMessage string
	lastFrame      int64
	elapsedMs      int64
}

func New(parent tea.Model, loadingMessage string) (Model, chan bool, chan string) {
	m := Model{
		parent:         parent,
		text:           make(chan string, 1),
		done:           make(chan bool, 1),
		currentMessage: loadingMessage,
		lastFrame:      time.Now().UnixMilli(),
	}
	return m, m.done, m.text
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if len(m.done) > 0 {
		<-m.done
		return m.parent, nil
	}

	if len(m.text) > 0 {
		m.currentMessage = <-m.text
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case LoaderTick:
		m.elapsedMs += time.Now().UnixMilli() - m.lastFrame
		m.lastFrame = time.Now().UnixMilli()
	}

	return m, tea.Tick(time.Second/5, func(t time.Time) tea.Msg {
		return LoaderTick("")
	})
}

func (m Model) Frame() int {
	return int(m.elapsedMs / (1000 / 10) % int64(len(loaderFrames)))
}

func (m Model) loadingIndicator() string {
	width := 20.0
	leftPad := ((1 + math.Sin(float64(m.elapsedMs)/1000)) / 2) * width
	rightPad := width - leftPad
	return "[ " + strings.Repeat(" ", int(leftPad)) + style.RedText.Render("=") + strings.Repeat(" ", int(rightPad)) + " ]"
}

func (m Model) View() string {
	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, lipgloss.JoinVertical(lipgloss.Center, style.RedText.Render(ui.Title), "", m.currentMessage, "", m.loadingIndicator()))
}
