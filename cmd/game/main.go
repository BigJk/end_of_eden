package main

import (
	"fmt"
	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/menus/mainmenu"
	"github.com/BigJk/project_gonzo/menus/root"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"os"
)

func main() {
	// Init audio
	audio.InitAudio()

	// Redirect log to file
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Init mouse zones
	zone.NewGlobal()
	zone.SetEnabled(true)

	// Run game
	p := tea.NewProgram(root.New(mainmenu.NewModel()), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
