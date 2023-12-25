package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/cmd/internal/testargs"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/system/gen"
	"github.com/BigJk/end_of_eden/system/gen/faces"
	"github.com/BigJk/end_of_eden/system/localization"
	"github.com/BigJk/end_of_eden/system/settings"
	"github.com/BigJk/end_of_eden/system/settings/viper"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
)

var prog *tea.Program
var loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	help := flag.Bool("help", false, "show help")
	testArgs := testargs.New()
	flag.Parse()

	// Start profiling server
	if env := os.Getenv("EOE_PROFILE"); env != "" {
		go func() {
			http.ListenAndServe(":8080", nil)
		}()
	}

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
		vi.SetDefault("language", "en")
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

	// Init Localization
	fmt.Println(loadStyle.Render("Initializing Localization. Please wait..."))
	{
		if err := localization.Global.AddFolder("./assets/locals"); err != nil {
			panic(err)
		}
		localization.SetCurrent(settings.GetString("language"))
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
		{Key: "language", Name: "Language", Description: fmt.Sprintf("Change the language (supported: %s)", strings.Join(localization.Global.GetLocales(), ", ")), Type: uiset.String, Val: settings.GetString("language")},
	}

	// Setup game
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones, settings.GetGlobal(), uiSettings, func(values []uiset.Value) error {
		for i := range values {
			settings.Set(values[i].Key, values[i].Val)
		}
		localization.SetCurrent(settings.GetString("language"))
		return settings.SaveSettings()
	}))

	// Apply test args if there are any
	baseModel = testArgs.ApplyArgs(baseModel, zones)

	// Run game
	prog = tea.NewProgram(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion(), tea.WithANSICompressor())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
