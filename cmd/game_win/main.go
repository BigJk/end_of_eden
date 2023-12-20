package main

import (
	"flag"
	"fmt"
	"github.com/BigJk/crt"
	teadapter "github.com/BigJk/crt/bubbletea"
	"github.com/BigJk/crt/shader"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/cmd/testargs"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/localization"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/settings/viper"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/ebiten/v2"
	zone "github.com/lrstanley/bubblezone"
	"image/color"
	"math"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"
)

var (
	loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)
)

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
}

func main() {
	testArgs := testargs.New()
	flag.Parse()

	// Start profiling server
	if env := os.Getenv("EOE_PROFILE"); env != "" {
		go func() {
			http.ListenAndServe(":8080", nil)
		}()
	}

	vi := viper.Viper{
		SettingsName: "settings_gl",
	}

	vi.SetDefault("audio", true)
	vi.SetDefault("volume", 1)
	vi.SetDefault("language", "en")
	vi.SetDefault("font_size", 12)
	vi.SetDefault("font_normal", "IosevkaTermNerdFontMono-Regular.ttf")
	vi.SetDefault("font_italic", "IosevkaTermNerdFontMono-Italic.ttf")
	vi.SetDefault("font_bold", "IosevkaTermNerdFontMono-Bold.ttf")
	vi.SetDefault("width", 1100)
	vi.SetDefault("height", 900)
	vi.SetDefault("crt", false)
	vi.SetDefault("grain", false)
	vi.SetDefault("show_fps", false)
	vi.SetDefault("fps", 30)

	settings.SetSettings(vi)

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))
	initSystems(settings.GetBool("audio"))

	fonts, err := crt.LoadFaces("./assets/fonts/"+settings.GetString("font_normal"), "./assets/fonts/"+settings.GetString("font_bold"), "./assets/fonts/"+settings.GetString("font_italic"), crt.GetFontDPI(), settings.GetFloat("font_size"))
	if err != nil {
		panic(err)
	}

	uiSettings := []uiset.Value{
		{Key: "audio", Name: "Audio", Description: "Enable or disable audio", Type: uiset.Bool, Val: settings.GetBool("audio"), Min: nil, Max: nil},
		{Key: "volume", Name: "Volume", Description: "Change the volume", Type: uiset.Float, Val: settings.GetFloat("volume"), Min: 0.0, Max: 2.0},
		{Key: "language", Name: "Language", Description: fmt.Sprintf("Change the language (supported: %s)", strings.Join(localization.Global.GetLocales(), ", ")), Type: uiset.String, Val: settings.GetString("language")},
		{Key: "font_size", Name: "Font Size", Description: "Change the font size", Type: uiset.Float, Val: settings.GetFloat("font_size"), Min: 6.0, Max: 64.0},
		{Key: "width", Name: "Width", Description: "Change the window width", Type: uiset.Float, Val: settings.GetFloat("width"), Min: 450.0, Max: 5000.0},
		{Key: "height", Name: "Height", Description: "Change the window height", Type: uiset.Float, Val: settings.GetFloat("height"), Min: 450.0, Max: 5000.0},
		{Key: "crt", Name: "CRT", Description: "Enable or disable CRT shader. Increases GPU usage.", Type: uiset.Bool, Val: settings.GetBool("crt"), Min: nil, Max: nil},
		{Key: "grain", Name: "Grain", Description: "Enable or disable grain shader. Matches well with the CRT shader. Increases GPU usage.", Type: uiset.Bool, Val: settings.GetBool("grain"), Min: nil, Max: nil},
		{Key: "show_fps", Name: "Show FPS", Description: "Show the current FPS", Type: uiset.Bool, Val: settings.GetBool("show_fps"), Min: nil, Max: nil},
		{Key: "fps", Name: "FPS", Description: "Change the FPS", Type: uiset.Float, Val: settings.GetFloat("fps"), Min: 10.0, Max: 144.0},
		{Key: "font_normal", Name: "Normal Font", Description: "Change the normal font", Type: uiset.String, Val: settings.GetString("font_normal"), Min: nil, Max: nil},
		{Key: "font_bold", Name: "Bold Font", Description: "Change the bold font", Type: uiset.String, Val: settings.GetString("font_bold"), Min: nil, Max: nil},
		{Key: "font_italic", Name: "Italic Font", Description: "Change the italic font", Type: uiset.String, Val: settings.GetString("font_italic"), Min: nil, Max: nil},
	}

	// Create base model
	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones, settings.GetGlobal(), uiSettings, func(values []uiset.Value) error {
		for i := range values {
			settings.Set(values[i].Key, values[i].Val)
		}
		localization.SetCurrent(settings.GetString("language"))
		return settings.SaveSettings()
	}))

	// Apply test args if given
	baseModel = testArgs.ApplyArgs(baseModel, zones)

	// Create window
	win, _, err := teadapter.Window(settings.GetInt("width"), settings.GetInt("height"), fonts, baseModel, color.RGBA{R: 34, G: 36, B: 41, A: 255}, tea.WithAltScreen())
	if err != nil {
		panic(err)
	}

	var loadedShader []shader.Shader

	// Setup CRT shader
	if settings.GetBool("crt") {
		crtLotte, err := shader.NewCrtLotte()
		if err != nil {
			panic(err)
		}

		crtLotte.Uniforms["WarpX"] = float32(0)
		crtLotte.Uniforms["WarpY"] = float32(0)

		// TODO: This is a bad hack to get the shader to change its state!
		go func() {
			warp := 0.0
			for {
				time.Sleep(time.Millisecond * 50)
				warp += 0.005
				win.Lock()
				{
					crtLotte.Uniforms["WarpX"] = float32(math.Abs(math.Sin(warp)*0.01) * 0.5)
					crtLotte.Uniforms["WarpY"] = float32(math.Abs(math.Sin(warp)*0.01) * 0.5)
				}
				win.Unlock()
			}
		}()

		loadedShader = append(loadedShader, crtLotte)
	}

	// Setup grain shader
	if settings.GetBool("grain") {
		res, _ := os.ReadFile("./assets/shader/grain.go")
		grain, err := ebiten.NewShader(res)

		if err != nil {
			panic(err)
		}

		w, h := win.Layout(ebiten.WindowSize())
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
			for {
				time.Sleep(time.Millisecond * 50)
				cur += 1
				win.Lock()
				{
					s.Uniforms["Tick"] = cur
				}
				win.Unlock()
			}
		}()

		loadedShader = append(loadedShader, s)
	}

	if len(loadedShader) > 0 {
		win.SetShader(loadedShader...)
	}

	// Run game
	win.ShowTPS(settings.GetBool("show_fps"))
	ebiten.SetTPS(settings.GetInt("fps"))
	if err := win.Run("End Of Eden"); err != nil {
		panic(err)
	}
}
