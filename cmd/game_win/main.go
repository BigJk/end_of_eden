package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/termgl"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/ebiten/v2"
	zone "github.com/lrstanley/bubblezone"
	"os"
)

var loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)

func initSystems(hasAudio bool) {
	// Init settings
	fmt.Println(loadStyle.Render("Initializing Settings. Please wait..."))
	{
		if err := settings.LoadSettings(); err != nil {
			panic(err)
		}
	}
	fmt.Println(loadStyle.Render("Done!"))

	// Init audio
	if hasAudio {
		fmt.Println(loadStyle.Render("Initializing Audio. Please wait..."))
		audio.InitAudio()
		audio.PlayMusic("planet_mining")
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
}

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	fontSize := flag.Float64("font_size", 16, "font size")
	dpiScaling := flag.Float64("dpi", 1, "scales the dpi up")
	width := flag.Int("width", 120, "window width in cells")
	height := flag.Int("height", 40, "window height in cells")
	help := flag.Bool("help", false, "show help")
	flag.Parse()

	if *help {
		fmt.Println("End Of Eden :: Game")
		fmt.Println()

		flag.PrintDefaults()
		return
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))
	initSystems(*audioFlag)

	gameInput := termgl.NewConcurrentRW()
	gameOutput := termgl.NewConcurrentRW()

	go gameInput.Run()
	go gameOutput.Run()

	// Start game backend
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones))

	prog := tea.NewProgram(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion(), tea.WithInput(gameInput), tea.WithOutput(gameOutput), tea.WithANSICompressor())

	go func() {
		if _, err := prog.Run(); err != nil {
			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}
	}()

	prog.Send(tea.WindowSizeMsg{
		Width:  *width - 1,
		Height: *height,
	})

	// Start game frontend
	dpi := *dpiScaling
	normal := termgl.LoadFace("./assets/fonts/IosevkaTermNerdFontMono-Regular.ttf", 72*dpi, *fontSize/dpi)
	bold := termgl.LoadFace("./assets/fonts/IosevkaTermNerdFontMono-Italic.ttf", 72*dpi, *fontSize/dpi)
	italic := termgl.LoadFace("./assets/fonts/IosevkaTermNerdFontMono-Bold.ttf", 72*dpi, *fontSize/dpi)

	game := termgl.NewGame(*width, *height, normal, bold, italic, gameOutput, prog)
	sw, sh := game.Layout(0, 0)

	ebiten.SetScreenFilterEnabled(false)
	ebiten.SetWindowSize(sw, sh)
	ebiten.SetWindowTitle("End Of Eden")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
