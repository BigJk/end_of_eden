package main

import (
	"fmt"
	"github.com/BigJk/crt"
	teadapter "github.com/BigJk/crt/bubbletea"
	"github.com/BigJk/crt/shader"
	"github.com/BigJk/end_of_eden/audio"
	"github.com/BigJk/end_of_eden/gen"
	"github.com/BigJk/end_of_eden/gen/faces"
	"github.com/BigJk/end_of_eden/settings"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/hajimehoshi/ebiten/v2"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/viper"
	"image/color"
	"math"
	"os"
	"time"
)

var (
	loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)
)

type settingsSaver struct{}

func (s settingsSaver) Save(values []uiset.Value) error {
	for i := range values {
		switch values[i].Key {
		// This is a special case for the volume setting, because it is not stored in the settings_win.toml file.
		case "volume":
			settings.LoadedSettings.Volume = values[i].Val.(float64)
			_ = settings.SaveSettings()
		default:
			viper.Set(values[i].Key, values[i].Val)
		}
	}

	return viper.WriteConfigAs("./settings_win.toml")
}

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
	viper.SetConfigName("settings_win")
	viper.AddConfigPath(".")

	viper.SetDefault("audio", true)
	viper.SetDefault("font_size", 12)
	viper.SetDefault("font_normal", "BigBlueTermPlusNerdFont-Regular.ttf")
	viper.SetDefault("font_italic", "BigBlueTermPlusNerdFont-Regular.ttf")
	viper.SetDefault("font_bold", "BigBlueTermPlusNerdFont-Regular.ttf")
	viper.SetDefault("dpi", 1)
	viper.SetDefault("width", 1300)
	viper.SetDefault("height", 975)
	viper.SetDefault("crt", true)
	viper.SetDefault("show_fps", false)
	viper.SetDefault("fps", 30)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			_ = viper.SafeWriteConfigAs("./settings_win.toml")
		} else {
			panic(err)
		}
	}

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))
	initSystems(viper.GetBool("audio"))

	dpi := viper.GetFloat64("dpi")
	fonts, err := crt.LoadFaces("./assets/fonts/"+viper.GetString("font_normal"), "./assets/fonts/"+viper.GetString("font_bold"), "./assets/fonts/"+viper.GetString("font_italic"), 72*dpi, viper.GetFloat64("font_size")/dpi)
	if err != nil {
		panic(err)
	}

	uiSettings := []uiset.Value{
		{Key: "audio", Name: "Audio", Description: "Enable or disable audio", Type: uiset.Bool, Val: viper.GetBool("audio"), Min: nil, Max: nil},
		{Key: "volume", Name: "Volume", Description: "Change the volume", Type: uiset.Float, Val: settings.LoadedSettings.Volume, Min: 0.0, Max: 1.0},
		{Key: "font_size", Name: "Font Size", Description: "Change the font size", Type: uiset.Float, Val: viper.GetFloat64("font_size"), Min: 6.0, Max: 64.0},
		{Key: "dpi", Name: "DPI", Description: "Change the DPI", Type: uiset.Float, Val: viper.GetFloat64("dpi"), Min: 1.0, Max: 5.0},
		{Key: "width", Name: "Width", Description: "Change the window width", Type: uiset.Float, Val: viper.GetFloat64("width"), Min: 450.0, Max: 5000.0},
		{Key: "height", Name: "Height", Description: "Change the window height", Type: uiset.Float, Val: viper.GetFloat64("height"), Min: 450.0, Max: 5000.0},
		{Key: "crt", Name: "CRT", Description: "Enable or disable CRT shader", Type: uiset.Bool, Val: viper.GetBool("crt"), Min: nil, Max: nil},
		{Key: "show_fps", Name: "Show FPS", Description: "Show the current FPS", Type: uiset.Bool, Val: viper.GetBool("show_fps"), Min: nil, Max: nil},
		{Key: "fps", Name: "FPS", Description: "Change the FPS", Type: uiset.Float, Val: viper.GetFloat64("fps"), Min: 10.0, Max: 144.0},
		{Key: "font_normal", Name: "Normal Font", Description: "Change the normal font", Type: uiset.String, Val: viper.GetString("font_normal"), Min: nil, Max: nil},
		{Key: "font_bold", Name: "Bold Font", Description: "Change the bold font", Type: uiset.String, Val: viper.GetString("font_bold"), Min: nil, Max: nil},
		{Key: "font_italic", Name: "Italic Font", Description: "Change the italic font", Type: uiset.String, Val: viper.GetString("font_italic"), Min: nil, Max: nil},
	}

	var baseModel tea.Model
	zones := zone.New()
	baseModel = root.New(zones, mainmenu.NewModel(zones, uiSettings, settingsSaver{}))
	win, _, err := teadapter.Window(viper.GetInt("width"), viper.GetInt("height"), fonts, baseModel, color.RGBA{
		R: 34,
		G: 36,
		B: 41,
		A: 255,
	}, tea.WithAltScreen())
	if err != nil {
		panic(err)
	}

	if viper.GetBool("crt") {
		res, _ := os.ReadFile("./assets/shader/grain.go")
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

		win.SetShader(crtLotte, s)
	}

	win.ShowTPS(viper.GetBool("show_fps"))
	ebiten.SetTPS(viper.GetInt("fps"))
	if err := win.Run("End Of Eden"); err != nil {
		panic(err)
	}
}
