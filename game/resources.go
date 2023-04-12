package game

import (
	"github.com/BigJk/project_gonzo/gluamapper"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"path/filepath"
	"strings"
)

// ResourcesManager can load Artifacts, Cards, Events and Enemy data from lua.
// The manager will walk the ./scripts directory and evaluate all found .lua files
type ResourcesManager struct {
	LuaState  *lua.LState
	Artifacts map[string]*Artifact
	Cards     map[string]*Card
	Events    map[string]*Event
	Enemies   map[string]*Enemy

	mapper *gluamapper.Mapper
}

func NewResourcesManager(state *lua.LState) *ResourcesManager {
	man := &ResourcesManager{
		LuaState:  state,
		Artifacts: map[string]*Artifact{},
		Cards:     map[string]*Card{},
		Events:    map[string]*Event{},
		Enemies:   map[string]*Enemy{},

		mapper: NewMapper(state),
	}

	// Attach all register methods
	man.LuaState.SetGlobal("register_artifact", man.LuaState.NewFunction(man.luaRegisterArtifact))
	man.LuaState.SetGlobal("register_card", man.LuaState.NewFunction(man.luaRegisterCard))
	man.LuaState.SetGlobal("register_enemy", man.LuaState.NewFunction(man.luaRegisterEnemy))
	man.LuaState.SetGlobal("register_event", man.LuaState.NewFunction(man.luaRegisterEvent))

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
		Callbacks: map[string]OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Artifacts[def.ID] = &def
	return 0
}

func (man *ResourcesManager) luaRegisterCard(l *lua.LState) int {
	def := Card{
		Callbacks: map[string]OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Cards[def.ID] = &def
	return 0
}

func (man *ResourcesManager) luaRegisterEnemy(l *lua.LState) int {
	def := Enemy{
		Callbacks: map[string]OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		// TODO: error handling
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)

	man.Enemies[def.ID] = &def
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
	return 0
}
