package lua_error

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/BigJk/end_of_eden/clipboard"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/maruel/panicparse/v2/stack"
	"github.com/samber/lo"
	"io"
	"path/filepath"
	"strings"
)

const (
	ZoneBack = "back"
	ZoneCopy = "copy"

	ErrorFormat = `%s%s

Go Call  = %s in line %d
Callback = %s
Type     = %s

%s

%s`
)

type Model struct {
	ui.MenuBase

	clipClicked bool
	zones       *zone.Manager
	err         game.LuaError
}

func New(zones *zone.Manager, luaErr game.LuaError) Model {
	errStr := luaErr.Err.Error()
	index := strings.Index(errStr, "stack traceback:")
	goErr := errStr[:index]
	rest := errStr[index:]
	if strings.Count(goErr, "\n") > 2 {
		res := &bytes.Buffer{}
		s, suffix, err := stack.ScanSnapshot(strings.NewReader(goErr), res, stack.DefaultOpts())
		if err == nil || err == io.EOF {
			// Find out similar goroutine traces and group them into buckets.
			buckets := s.Aggregate(stack.AnyValue).Buckets

			// Calculate alignment.
			srcLen := 0
			pkgLen := 0
			for _, bucket := range buckets {
				for _, line := range bucket.Signature.Stack.Calls {
					if l := len(fmt.Sprintf("%s:%d", line.SrcName, line.Line)); l > srcLen {
						srcLen = l
					}
					if l := len(filepath.Base(line.Func.ImportPath)); l > pkgLen {
						pkgLen = l
					}
				}
			}

			for _, bucket := range buckets {
				// Print the goroutine header.
				extra := ""
				if s := bucket.SleepString(); s != "" {
					extra += " [" + s + "]"
				}
				if bucket.Locked {
					extra += " [locked]"
				}

				if len(bucket.CreatedBy.Calls) != 0 {
					extra += fmt.Sprintf(" [Created by %s.%s @ %s:%d]", bucket.CreatedBy.Calls[0].Func.DirName, bucket.CreatedBy.Calls[0].Func.Name, bucket.CreatedBy.Calls[0].SrcName, bucket.CreatedBy.Calls[0].Line)
				}
				res.WriteString(fmt.Sprintf("%d: %s%s\n", len(bucket.IDs), bucket.State, extra))

				// Print the stack lines.
				for _, line := range bucket.Stack.Calls {
					text := fmt.Sprintf(
						"    %-*s %-*s %s(%s)",
						pkgLen, line.Func.DirName, srcLen,
						fmt.Sprintf("%s:%d", line.SrcName, line.Line),
						line.Func.Name, &line.Args)

					if line.Func.DirName == "game" {
						text = lipgloss.NewStyle().Bold(true).Foreground(style.BaseYellow).Render(text)
					} else {
						text = lipgloss.NewStyle().Foreground(style.BaseGray).Render(text)
					}

					res.WriteString(text + "\n")
				}
				if bucket.Stack.Elided {
					_, _ = io.WriteString(res, "    (...)\n")
				}
			}

			// If there was any remaining data in the pipe, dump it now.
			if len(suffix) != 0 {
				res.Write(suffix)
			}

			luaErr.Err = errors.New(fmt.Sprintf("%s\n\n%s", res.String(), rest))
		}
	}

	luaErr.File, _ = lo.Last(strings.Split(luaErr.File, "/"))

	return Model{
		MenuBase: ui.MenuBase{},
		zones:    zones,
		err:      luaErr,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Size = msg
	case tea.KeyMsg:
		if msg.Type == tea.KeyEsc {
			return nil, nil
		}
	case tea.MouseMsg:
		if msg.Type == tea.MouseLeft {
			if m.zones.Get(ZoneCopy).InBounds(msg) {
				clipboard.Set(fmt.Sprintf(ErrorFormat, "", "Lua Error!", m.err.File, m.err.Line, m.err.Callback, m.err.Type, "Error:", strings.Replace(m.err.Err.Error(), "\t", " ", -1)))
				m.clipClicked = true
			} else if m.zones.Get(ZoneBack).InBounds(msg) {
				return nil, nil
			}
		}
		m.LastMouse = msg
	}

	return m, nil
}

func (m Model) View() string {
	err := lipgloss.NewStyle().Width(m.Size.Width-30).Border(lipgloss.ThickBorder(), true).Padding(1, 2, 0, 1).BorderForeground(style.BaseGray).Foreground(style.BaseWhite).Render(
		fmt.Sprintf(ErrorFormat, style.RedText.Copy().Bold(true).Render("Lua Error!"), `

If you want to report this error please use "Copy Clipboard"
and provide the result together with information of what you
were doing at the moment of error.`, m.err.File, m.err.Line, m.err.Callback, m.err.Type, style.RedText.Copy().Bold(true).Render("Error:"), strings.Replace(m.err.Err.Error(), "\t", " ", -1)) +
			"\n" +
			lipgloss.JoinHorizontal(
				lipgloss.Left,
				style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneBack).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Render(m.zones.Mark(ZoneBack, "Back")),
				style.HeaderStyle.Copy().Background(lo.Ternary(m.zones.Get(ZoneCopy).InBounds(m.LastMouse), style.BaseRed, style.BaseRedDarker)).Render(m.zones.Mark(ZoneCopy, lo.Ternary(m.clipClicked, "Copied!", "Copy Clipboard"))),
			),
	)

	return lipgloss.Place(m.Size.Width, m.Size.Height, lipgloss.Center, lipgloss.Center, err, lipgloss.WithWhitespaceChars("!"), lipgloss.WithWhitespaceForeground(style.BaseGrayDarker))
}
