//go:build windows
// +build windows

package main

import (
	"os"
	"time"

	"github.com/Azure/go-ansiterm/winterm"
	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/sys/windows"
)

// Enable ANSI color support on windows in default terminal
// and change the console to a fixed size at first.
func init() {
	initialSize := winterm.COORD{
		X: 150,
		Y: 50,
	}
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)

	/* Disable for now as microsoft/terminal works well with EoE

	winterm.SetConsoleWindowInfo(uintptr(stdout), true, winterm.SMALL_RECT{
		Left:   0,
		Top:    0,
		Right:  initialSize.X - 1,
		Bottom: initialSize.Y - 1,
	})
	winterm.SetConsoleScreenBufferSize(uintptr(stdout), initialSize)

	*/

	// Workaround to enable re-size behaviour.
	go func() {
		for prog == nil {
			time.Sleep(time.Second)
		}

		for {
			if info, err := winterm.GetConsoleScreenBufferInfo(uintptr(stdout)); err == nil {
				if info.Size.X != initialSize.X || info.Size.Y != initialSize.Y {
					initialSize = info.Size
					prog.Send(tea.WindowSizeMsg{
						Width:  int(info.Size.X),
						Height: int(info.Size.Y),
					})
				}
			}

			time.Sleep(time.Second)
		}
	}()
}
