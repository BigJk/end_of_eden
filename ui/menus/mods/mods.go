package mods

import (
	"github.com/BigJk/end_of_eden/game"
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/system/settings"
	"github.com/BigJk/end_of_eden/ui"
	"github.com/BigJk/end_of_eden/ui/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"log"
	"path/filepath"
	"sort"
)

type item struct {
	active bool
	mod    game.Mod
	key    string
}

func (i item) Title() string {
	if !i.active {
		return i.mod.Name + style.RedDarkerText.Render(" by ") + style.GrayText.Render(i.mod.Author)
	}

	return lipgloss.NewStyle().Italic(true).Foreground(style.BaseGreen).Render("Active") + " " + style.RedText.Render(i.mod.Name) + style.RedDarkerText.Render(" by ") + style.GrayText.Render(i.mod.Author)
}
func (i item) Description() string { return i.mod.Description }
func (i item) FilterValue() string { return i.mod.Name }

type Model struct {
	ui.MenuBase

	settings settings.Settings
	list     list.Model
	mods     map[string]game.Mod
	zones    *zone.Manager
}

func NewModel(zones *zone.Manager, settings settings.Settings) Model {
	return Model{
		settings: settings,
		zones:    zones,
	}.setup()
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "d":
			m = m.modDown(m.list.SelectedItem().(item).key)
		case "u":
			m = m.modUp(m.list.SelectedItem().(item).key)
		case "q":
			fallthrough
		case "esc":
			return nil, nil
		}

		switch msg.Type {
		case tea.KeyDown:
			fallthrough
		case tea.KeyUp:
			audio.Play("interface_move", -1.5)
		case tea.KeyEnter:
			m = m.modSetActive(m.list.SelectedItem().(item).key, !m.list.SelectedItem().(item).active)
		}
	case tea.WindowSizeMsg:
		m.Size = msg
		m.list.SetSize(msg.Width-4, msg.Height-2)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	return lipgloss.NewStyle().Padding(1, 2).Render(m.list.View())
}

func (m Model) setup() Model {
	m = m.fetchMods()

	delegation := list.NewDefaultDelegate()
	delegation.Styles.SelectedTitle = delegation.Styles.SelectedTitle.Foreground(style.BaseRed).BorderForeground(style.BaseRed)
	delegation.Styles.SelectedDesc = delegation.Styles.SelectedDesc.Foreground(style.BaseRedDarker).BorderForeground(style.BaseRed)

	m.list = list.New(m.items(), delegation, 0, 0)
	m.list.Title = "Mods"
	m.list.SetFilteringEnabled(false)
	m.list.SetShowFilter(false)
	m.list.SetShowStatusBar(false)
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(key.WithKeys("u"), key.WithHelp("u", "move up")),
			key.NewBinding(key.WithKeys("d"), key.WithHelp("d", "move down")),
			key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "toggle mod")),
		}
	}
	m.list.AdditionalFullHelpKeys = m.list.AdditionalShortHelpKeys
	m.list.Styles.Title = lipgloss.NewStyle().Background(style.BaseRedDarker).Foreground(style.BaseWhite).Padding(0, 2, 0, 2)

	return m
}

func (m Model) items() []list.Item {
	baseKeys := lo.Keys(m.mods)
	sort.Strings(baseKeys)

	keys := lo.Uniq(append(m.settings.GetStrings("mods"), baseKeys...))
	items := lo.FilterMap(keys, func(modName string, _ int) (list.Item, bool) {
		mod, ok := m.mods[modName]
		if !ok {
			return item{}, false
		}
		return item{
			active: m.modActive(modName),
			key:    modName,
			mod:    mod,
		}, true
	})

	return items
}

func (m Model) modUp(mod string) Model {
	index := lo.IndexOf(m.settings.GetStrings("mods"), mod)
	if index <= 0 {
		return m
	}

	m.settings.GetStrings("mods")[index] = m.settings.GetStrings("mods")[index-1]
	m.settings.GetStrings("mods")[index-1] = mod
	_ = m.settings.SaveSettings()

	m.list.SetItems(m.items())
	return m
}

func (m Model) modDown(mod string) Model {
	index := lo.IndexOf(m.settings.GetStrings("mods"), mod)
	if index < 0 || index >= len(m.settings.GetStrings("mods"))-1 {
		return m
	}

	m.settings.GetStrings("mods")[index] = m.settings.GetStrings("mods")[index+1]
	m.settings.GetStrings("mods")[index+1] = mod
	_ = m.settings.SaveSettings()

	m.list.SetItems(m.items())
	return m
}

func (m Model) modSetActive(mod string, val bool) Model {
	if val {
		m.settings.Set("mods", append(m.settings.GetStrings("mods"), mod))
	} else {
		m.settings.Set("mods", lo.Filter(m.settings.GetStrings("mods"), func(item string, index int) bool {
			return item != mod
		}))
	}
	_ = m.settings.SaveSettings()

	m.list.SetItems(m.items())
	return m
}

func (m Model) modActive(mod string) bool {
	return lo.Contains(m.settings.GetStrings("mods"), mod)
}

func (m Model) fetchMods() Model {
	entries, err := fs.ReadDir("./mods")
	if err != nil {
		log.Println("Error while reading mods directory:", err)
		return m
	}

	mods := map[string]game.Mod{}
	for _, e := range entries {
		if e.IsDir() {
			mod, err := game.ModDescription(filepath.Join("./mods", e.Name()))
			if err != nil {
				log.Println("Error while reading mod:", e.Name(), err)
			} else {
				mods[e.Name()] = mod
			}
		}
	}

	m.mods = mods
	return m
}
