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
	audio.InitAudio()

	zone.NewGlobal()
	zone.SetEnabled(true)

	p := tea.NewProgram(root.New(mainmenu.NewModel()), tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
