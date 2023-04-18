package main

import (
	"fmt"
	zone "github.com/lrstanley/bubblezone"
	"log"
	"os"

	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/ui/mainmenu"
	"github.com/BigJk/project_gonzo/ui/root"
	tea "github.com/charmbracelet/bubbletea"
)

var prog *tea.Program

func main() {
	// Init audio
	audio.InitAudio()

	// Redirect log to file
	_ = os.Mkdir("./logs", 0777)
	f, err := tea.LogToFile("./logs/global.log", "")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	log.Println("=================================")
	log.Println("= Started")
	log.Println("=================================")

	// Run game
	zones := zone.New()
	prog = tea.NewProgram(root.New(zones, mainmenu.NewModel(zones)), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
