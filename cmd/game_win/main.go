package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/crt"
	teadapter "github.com/BigJk/crt/bubbletea"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"image/color"
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
	width := flag.Int("width", 1300, "window width")
	height := flag.Int("height", 975, "window height")
	help := flag.Bool("help", false, "show help")
	crtShader := flag.Bool("crt", true, "enable crt shader")
	flag.Parse()

	if *help {
		fmt.Println("End Of Eden :: Game")
		fmt.Println()

		flag.PrintDefaults()
		return
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))
	initSystems(*audioFlag)

	dpi := *dpiScaling
	fonts, err := crt.LoadFaces("./assets/fonts/IosevkaTermNerdFontMono-Regular.ttf", "./assets/fonts/IosevkaTermNerdFontMono-Bold.ttf", "./assets/fonts/IosevkaTermNerdFontMono-Italic.ttf", 72*dpi, *fontSize/dpi)
	if err != nil {
		panic(err)
	}

	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones))
	win, err := teadapter.Window(*width, *height, fonts, baseModel, color.RGBA{
		R: 34,
		G: 36,
		B: 41,
		A: 255,
	}, tea.WithAltScreen())
	if err != nil {
		panic(err)
	}

	win.CRTShader(*crtShader)
	if err := win.Run("End Of Eden"); err != nil {
		panic(err)
	}
}
