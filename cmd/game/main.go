package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/ui/mainmenu"
	"github.com/BigJk/project_gonzo/ui/root"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

var prog *tea.Program

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

	log.Println("=================================")
	log.Println("= Started")
	log.Println("=================================")

	// Run game
	prog = tea.NewProgram(root.New(mainmenu.NewModel()), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
