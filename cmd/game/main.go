package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/settings/viper"
	"github.com/BigJk/end_of_eden/ui/menus/gameview"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
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
	testEvent := flag.String("event", "", "test event")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		fmt.Println("End Of Eden :: Game")
		fmt.Println()

		flag.PrintDefaults()
		return
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))

	// Init settings
	fmt.Println(loadStyle.Render("Initializing Settings. Please wait..."))
	{
		vi := viper.Viper{
			SettingsName: "settings_term",
		}
		vi.SetDefault("audio", true)
		vi.SetDefault("volume", 1)
		settings.SetSettings(vi)

		if err := settings.LoadSettings(); err != nil {
			panic(err)
		}
	}
	fmt.Println(loadStyle.Render("Done!"))

	// Init audio
	if *audioFlag {
		fmt.Println(loadStyle.Render("Initializing Audio. Please wait..."))
		audio.InitAudio()
		fmt.Println(loadStyle.Render("Done!"))
	}

	// Init generators
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

	uiSettings := []uiset.Value{
		{Key: "audio", Name: "Audio", Description: "Enable or disable audio", Type: uiset.Bool, Val: settings.GetBool("audio"), Min: nil, Max: nil},
		{Key: "volume", Name: "Volume", Description: "Change the volume", Type: uiset.Float, Val: settings.GetFloat("volume"), Min: 0.0, Max: 2.0},
	}

	// Setup game
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones, settings.GetGlobal(), uiSettings, func(values []uiset.Value) error {
		for i := range values {
			settings.Set(values[i].Key, values[i].Val)
		}
		return settings.SaveSettings()
	}))

	// If test flags are present we load up a session with the given cards, enemies and artifacts.
	if len(*testCards) > 0 || len(*testEnemies) > 0 || len(*testArtifacts) > 0 || len(*testGameState) > 0 || len(*testEvent) > 0 {
		session := game.NewSession(game.WithLogging(log.Default()), game.WithMods(settings.GetStrings("mods")), lo.Ternary(os.Getenv("EOE_DEBUG") == "1", game.WithDebugEnabled(8272), nil))
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

		if len(*testEvent) > 0 {
			session.SetGameState(game.GameStateEvent)
			session.SetEvent(*testEvent)
		}

		if len(*testGameState) > 0 {
			session.SetGameState(game.GameState(*testGameState))
		}

		baseModel = baseModel.(root.Model).PushModel(gameview.New(baseModel, zones, session))
	}

	// Run game
	prog = tea.NewProgram(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion(), tea.WithANSICompressor())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
