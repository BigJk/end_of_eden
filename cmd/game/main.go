package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
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
var loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	testCards := flag.String("cards", "", "test cards")
	testEnemies := flag.String("enemies", "", "test enemies")
	testArtifacts := flag.String("artifacts", "", "test artifacts")
	testGameState := flag.String("game_state", "", "test game state")
	flag.Parse()

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))

	fmt.Println(loadStyle.Render("Initializing Settings. Please wait..."))
	{
		if err := settings.LoadSettings(); err != nil {
			panic(err)
		}
	}
	fmt.Println(loadStyle.Render("Done!"))

	// Init clipboard
	fmt.Println(loadStyle.Render("Initializing Clipboard. Please wait..."))
	{
		if err := clipboard.Init(); err != nil {
			panic(err)
		}
	}
	fmt.Println(loadStyle.Render("Done!"))

	// Init audio
	if *audioFlag {
		fmt.Println(loadStyle.Render("Initializing audio. Please wait..."))
		audio.InitAudio()
		audio.PlayMusic("planet_mining")
		fmt.Println(loadStyle.Render("Done!"))
	}

	fmt.Println(loadStyle.Render("Initializing Proc-Gen. Please wait..."))
	{
		// Init face generator
		if err := faces.InitGlobal("./assets/gen/faces"); err != nil {
			panic(err)
		}

		// Init other gens
		gen.InitGen()
	}
	fmt.Println(loadStyle.Render("Done!"))

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
		session := game.NewSession(game.WithLogging(log.Default()), game.WithMods(settings.LoadedSettings.Mods))
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
