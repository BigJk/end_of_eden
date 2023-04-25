package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"golang.design/x/clipboard"
	"log"
	"os"
	"strings"

	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/ui/root"
	tea "github.com/charmbracelet/bubbletea"
)

var prog *tea.Program

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	testCards := flag.String("cards", "", "test cards")
	testEnemies := flag.String("enemies", "", "test enemies")
	testArtifacts := flag.String("artifacts", "", "test artifacts")
	testGameState := flag.String("game_state", "", "test game state")
	flag.Parse()

	// Init clipboard
	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	// Init audio
	if *audioFlag {
		audio.InitAudio()
		audio.PlayMusic("theme")
	}

	// Init face generator
	if err := faces.InitGlobal("./assets/gen/faces"); err != nil {
		panic(err)
	}

	// Init other gens
	gen.InitGen()

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

	// Set window title
	fmt.Println("\033]2;End of Eden\007")

	// Setup game
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones))

	// If test flags are present we load up a session with the given cards, enemies and artifacts.
	if len(*testCards) > 0 || len(*testEnemies) > 0 || len(*testArtifacts) > 0 || len(*testGameState) > 0 {
		session := game.NewSession(game.WithLogging(log.Default()))
		session.SetGameState(game.GameStateFight)
		session.GetPlayer().Cards.Clear()

		if len(*testEnemies) == 0 {
			*testEnemies = "DUMMY,DUMMY,DUMMY"
		}

		lo.ForEach(strings.Split(*testCards, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.GiveCard(item, game.PlayerActorID)
		})

		lo.ForEach(strings.Split(*testEnemies, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.AddActorFromEnemy(item)
		})

		lo.ForEach(strings.Split(*testArtifacts, ","), func(item string, index int) {
			if len(item) == 0 {
				return
			}
			session.GiveArtifact(item, game.PlayerActorID)
		})

		session.SetupFight()

		if len(*testGameState) > 0 {
			session.SetGameState(game.GameState(*testGameState))
		}

		baseModel = baseModel.(root.Model).PushModel(gameview.New(baseModel, zones, session))
	}

	// Run game
	prog = tea.NewProgram(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
