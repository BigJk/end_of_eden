//go:build windows
// +build windows

package main

import (
	"os"

	"golang.org/x/sys/windows"
)

// Enable ANSI color support on windows in default terminal.
func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}
