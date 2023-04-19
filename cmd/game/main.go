package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/project_gonzo/game"
	"github.com/BigJk/project_gonzo/ui/gameview"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"log"
	"os"
	"strings"

	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/ui/mainmenu"
	"github.com/BigJk/project_gonzo/ui/root"
	tea "github.com/charmbracelet/bubbletea"
)

var prog *tea.Program

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	testCards := flag.String("cards", "", "test cards")
	testEnemies := flag.String("enemies", "", "test enemies")
	testArtifacts := flag.String("artifacts", "", "test artifacts")
	flag.Parse()

	// Init audio
	if *audioFlag {
		audio.InitAudio()
	}

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

	// Setup game
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones))

	// If test flags are present we load up a session with the given cards, enemies and artifacts.
	if len(*testCards) > 0 || len(*testEnemies) > 0 || len(*testArtifacts) > 0 {
		session := game.NewSession(game.WithLogging(log.Default()))
		session.SetGameState(game.GameStateFight)
		session.GetPlayer().Cards.Clear()

		if len(*testEnemies) == 0 {
			*testEnemies = "DUMMY,DUMMY,DUMMY"
		}

		lo.ForEach(strings.Split(*testCards, ","), func(item string, index int) {
			session.GiveCard(item, game.PlayerActorID)
		})

		lo.ForEach(strings.Split(*testEnemies, ","), func(item string, index int) {
			session.AddActorFromEnemy(item)
		})

		lo.ForEach(strings.Split(*testArtifacts, ","), func(item string, index int) {
			session.GiveArtifact(item, game.PlayerActorID)
		})

		session.SetupFight()
		baseModel = baseModel.(root.Model).SetModel(gameview.New(baseModel, zones, session))
	}

	// Run game
	prog = tea.NewProgram(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
