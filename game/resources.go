package game

import (
	"github.com/BigJk/project_gonzo/luhelp"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"path/filepath"
	"strings"
)

// ResourcesManager can load Artifacts, Cards, Events, Enemy and StoryTeller data from lua.
// The manager will walk the ./scripts directory and evaluate all found .lua files
type ResourcesManager struct {
	LuaState *lua.LState

	Artifacts     map[string]*Artifact
	Cards         map[string]*Card
	Events        map[string]*Event
	Enemies       map[string]*Enemy
	StatusEffects map[string]*StatusEffect
	StoryTeller   map[string]*StoryTeller

	registered *lua.LTable
	mapper     *luhelp.Mapper
}

func NewResourcesManager(state *lua.LState) *ResourcesManager {
	man := &ResourcesManager{
		LuaState:      state,
		Artifacts:     map[string]*Artifact{},
		Cards:         map[string]*Card{},
		Events:        map[string]*Event{},
		Enemies:       map[string]*Enemy{},
		StatusEffects: map[string]*StatusEffect{},
		StoryTeller:   map[string]*StoryTeller{},

		registered: state.NewTable(),
		mapper:     luhelp.NewMapper(state),
	}

	// Create global variable to access registered values in lua
	lo.ForEach([]string{"artifact", "card", "enemy", "event", "status_effect", "story_teller"}, func(t string, _ int) {
		man.registered.RawSetString(t, state.NewTable())
	})
	man.LuaState.SetGlobal("registered", man.registered)

	// Attach all register methods
	man.LuaState.SetGlobal("register_artifact", man.LuaState.NewFunction(man.luaRegisterArtifact))
	man.LuaState.SetGlobal("register_card", man.LuaState.NewFunction(man.luaRegisterCard))
	man.LuaState.SetGlobal("register_enemy", man.LuaState.NewFunction(man.luaRegisterEnemy))
	man.LuaState.SetGlobal("register_event", man.LuaState.NewFunction(man.luaRegisterEvent))
	man.LuaState.SetGlobal("register_status_effect", man.LuaState.NewFunction(man.luaRegisterStatusEffect))
	man.LuaState.SetGlobal("register_story_teller", man.LuaState.NewFunction(man.luaRegisterStoryTeller))

	// Load all local scripts
	_ = filepath.Walk("./assets/scripts", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".lua") {
			if err := man.LuaState.DoFile(path); err != nil {
				// TODO: error handling
				panic(err)
			}
		}

		return nil
	})

	return man
}

func (man *ResourcesManager) luaRegisterArtifact(l *lua.LState) int {
	def := Artifact{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Artifacts[def.ID] = &def
	man.registered.RawGetString("artifact").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterCard(l *lua.LState) int {
	def := Card{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Cards[def.ID] = &def
	man.registered.RawGetString("card").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterEnemy(l *lua.LState) int {
	def := Enemy{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Enemies[def.ID] = &def
	man.registered.RawGetString("enemy").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterEvent(l *lua.LState) int {
	def := Event{}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Events[def.ID] = &def
	man.registered.RawGetString("event").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterStatusEffect(l *lua.LState) int {
	def := StatusEffect{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.StatusEffects[def.ID] = &def
	man.registered.RawGetString("status_effect").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterStoryTeller(l *lua.LState) int {
	def := StoryTeller{}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.StoryTeller[def.ID] = &def
	man.registered.RawGetString("story_teller").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}
