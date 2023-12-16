package gifviewer

import (
	"fmt"
	"github.com/BigJk/end_of_eden/image"
	"github.com/BigJk/end_of_eden/ui"
	tea "github.com/charmbracelet/bubbletea"
	"math/rand"
	"time"
)

type GifAnimationFinished string

type GifAnimationFrame string

type Model struct {
	ui.MenuBase

	id     string
	parent tea.Model

	fps       int
	lastFrame int64
	elapsedMs int64
	frames    []string
}

func New(parent tea.Model, file string, fps int, width int, height int) (tea.Model, error) {
	var options []image.Option

	if width > 0 && height == 0 {
		options = append(options, image.WithMaxWidth(width))
	} else if width > 0 && height > 0 {
		options = append(options, image.WithResize(width, height))
	} else {
		options = append(options, image.WithMaxWidth(100))
	}

	frames, err := image.FetchAnimation(file, options...)
	if err != nil {
		return nil, err
	}

	return Model{id: fmt.Sprint(rand.Int()), parent: parent, fps: fps, frames: frames, lastFrame: time.Now().UnixMilli()}, nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case GifAnimationFrame:
		if string(msg) != m.id {
			return m, nil
		}

		m.elapsedMs += time.Now().UnixMilli() - m.lastFrame
		m.lastFrame = time.Now().UnixMilli()

		if m.Frame() == len(m.frames)-1 {
			return m, func() tea.Msg {
				return GifAnimationFinished(m.id)
			}
		}
	}

	return m, tea.Tick(time.Second/time.Duration(m.fps), func(t time.Time) tea.Msg {
		return GifAnimationFrame(m.id)
	})
}

func (m Model) Frame() int {
	return int(m.elapsedMs / (1000 / int64(m.fps)) % int64(len(m.frames)))
}

func (m Model) View() string {
	return m.frames[m.Frame()]
}
