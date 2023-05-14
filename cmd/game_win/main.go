package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/BigJk/crt"
	teadapter "github.com/BigJk/crt/bubbletea"
	"github.com/BigJk/crt/shader"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/fs"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/ebiten/v2"
	zone "github.com/lrstanley/bubblezone"
	"image/color"
	"math"
	"time"
)

var (
	//go:embed IosevkaTermNerdFontMono-Regular.ttf
	FontNormal []byte
	//go:embed IosevkaTermNerdFontMono-Bold.ttf
	FontBold []byte
	//go:embed IosevkaTermNerdFontMono-Italic.ttf
	FontItalic []byte
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
	/*fmt.Println(loadStyle.Render("Initializing Proc-Gen. Please wait..."))
	{
		// Init face generator
		if err := faces.InitGlobal("./assets/gen/faces"); err != nil {
			panic(err)
		}

		// Init other gens
		gen.InitGen()
	}
	fmt.Println(loadStyle.Render("Done!"))*/
}

func main() {
	audioFlag := flag.Bool("audio", true, "disable audio")
	fontSize := flag.Float64("font_size", 16, "font size")
	dpiScaling := flag.Float64("dpi", 1, "scales the dpi up")
	width := flag.Int("width", 1300, "window width")
	height := flag.Int("height", 975, "window height")
	help := flag.Bool("help", false, "show help")
	crtShader := flag.Bool("crt", false, "enable crt shader")
	maxFps := flag.Int("fps", 30, "max fps")
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
	fonts, err := crt.LoadFacesBytes(FontNormal, FontBold, FontItalic, 72*dpi, *fontSize/dpi)
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

	if *crtShader {
		res, _ := fs.ReadFile("./assets/shader/grain.go")
		grain, err := ebiten.NewShader(res)

		if err != nil {
			panic(err)
		}

		crtLotte, err := shader.NewCrtLotte()
		if err != nil {
			panic(err)
		}

		crtLotte.Uniforms["WarpX"] = float32(0)
		crtLotte.Uniforms["WarpY"] = float32(0)

		w, h := win.Layout(0, 0)
		s := &shader.BaseShader{
			Shader: grain,
			Uniforms: map[string]any{
				"ScreenSize": []float32{float32(w), float32(h)},
				"Tick":       float32(0),
				"Strength":   float32(0.05),
			},
		}

		// TODO: This is a bad hack to get the shader to change its state!
		go func() {
			cur := float32(0)
			warp := 0.0
			for {
				time.Sleep(time.Millisecond * 50)

				cur += 1
				warp += 0.005

				win.Lock()
				{
					s.Uniforms["Tick"] = cur
					crtLotte.Uniforms["WarpX"] = float32(math.Abs(math.Sin(warp)*0.01) * 0.5)
					crtLotte.Uniforms["WarpY"] = float32(math.Abs(math.Sin(warp)*0.01) * 0.5)
				}
				win.Unlock()
			}
		}()

		win.ShowTPS(true)
		win.SetShader(crtLotte, s)
	}

	ebiten.SetTPS(*maxFps)
	if err := win.Run("End Of Eden"); err != nil {
		panic(err)
	}
}
