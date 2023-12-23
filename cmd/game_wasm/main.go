//go:build js
// +build js

package main

import (
	"bytes"
	"fmt"
	"github.com/BigJk/end_of_eden/system/gen"
	"github.com/BigJk/end_of_eden/system/gen/faces"
	"github.com/BigJk/end_of_eden/system/localization"
	"github.com/BigJk/end_of_eden/system/settings"
	"github.com/BigJk/end_of_eden/system/settings/browser"
	"github.com/BigJk/end_of_eden/ui/menus/mainmenu"
	uiset "github.com/BigJk/end_of_eden/ui/menus/settings"
	"github.com/BigJk/end_of_eden/ui/menus/warning"
	"github.com/BigJk/end_of_eden/ui/root"
	"github.com/BigJk/end_of_eden/ui/style"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/termenv"
	"log"
	_ "net/http/pprof"
	"os"
	"strings"
	"syscall/js"
	"time"
)

type MinReadBuffer struct {
	buf *bytes.Buffer
}

// For some reason bubbletea doesn't like a Reader that will return 0 bytes instead of blocking,
// so we use this hacky workaround for now. As javascript is single threaded this should be fine
// with regard to concurrency.
func (b *MinReadBuffer) Read(p []byte) (n int, err error) {
	for b.buf.Len() == 0 {
		time.Sleep(100 * time.Millisecond)
	}
	return b.buf.Read(p)
}

func (b *MinReadBuffer) Write(p []byte) (n int, err error) {
	return b.buf.Write(p)
}

// Creates the bubbletea program and registers the necessary functions in javascript
func createTeaForJS(model tea.Model, option ...tea.ProgramOption) *tea.Program {
	// Create buffers for input and output
	fromJs := &MinReadBuffer{buf: bytes.NewBuffer(nil)}
	fromGo := bytes.NewBuffer(nil)

	prog := tea.NewProgram(model, append([]tea.ProgramOption{tea.WithInput(fromJs), tea.WithOutput(fromGo)}, option...)...)

	// Register write function in WASM
	js.Global().Set("bubbletea_write", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		fromJs.Write([]byte(args[0].String()))
		return nil
	}))

	// Register read function in WASM
	js.Global().Set("bubbletea_read", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		b := make([]byte, fromGo.Len())
		_, _ = fromGo.Read(b)
		fromGo.Reset()
		return string(b)
	}))

	// Register resize function in WASM
	js.Global().Set("bubbletea_resize", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		width := args[0].Int()
		height := args[1].Int()
		prog.Send(tea.WindowSizeMsg{Width: width, Height: height})
		return nil
	}))

	return prog
}

var prog *tea.Program
var loadStyle = lipgloss.NewStyle().Bold(true).Italic(true).Foreground(style.BaseGray)

func main() {
	lipgloss.SetColorProfile(termenv.TrueColor)

	fmt.Println(lipgloss.NewStyle().Bold(true).Foreground(style.BaseRed).Render("End Of Eden"))

	// Init settings
	fmt.Println(loadStyle.Render("Initializing Settings. Please wait..."))
	{
		set := browser.Browser{}
		set.SetDefault("audio", true)
		set.SetDefault("volume", 1)
		set.SetDefault("language", "en")
		settings.SetSettings(set)

		if err := settings.LoadSettings(); err != nil {
			panic(err)
		}
	}
	fmt.Println(loadStyle.Render("Done!"))

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

	log.Println("=================================")
	log.Println("= Started")
	log.Println("=================================")

	// Set window title
	fmt.Println("\033]2;End of Eden\007")

	uiSettings := []uiset.Value{
		{Key: "audio", Name: "Audio", Description: "Enable or disable audio", Type: uiset.Bool, Val: settings.GetBool("audio"), Min: nil, Max: nil},
		{Key: "volume", Name: "Volume", Description: "Change the volume", Type: uiset.Float, Val: settings.GetFloat("volume"), Min: 0.0, Max: 2.0},
		{Key: "language", Name: "Language", Description: fmt.Sprintf("Change the language (supported: %s)", strings.Join(localization.Global.GetLocales(), ", ")), Type: uiset.String, Val: settings.GetString("language")},
		{Key: "font_size", Name: "Font Size", Description: "Change the font size (page reload required)", Type: uiset.Float, Val: settings.GetFloat("font_size"), Min: 6.0, Max: 64.0},
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

	baseModel = baseModel.(root.Model).PushModel(warning.New(nil, style.RedText.Render("Warning!")+"\n\nThe Browser version is still very experimental. Loading times can be long. Mouse support is clunky. For the best experience, please use the Desktop version!\n\n"+style.GrayTextDarker.Render("Press ESC to continue")))

	// Run game
	prog = createTeaForJS(baseModel, tea.WithAltScreen(), tea.WithMouseAllMotion(), tea.WithANSICompressor())
	if _, err := prog.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
