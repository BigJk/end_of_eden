package game

import (
	"fmt"
	"github.com/BigJk/end_of_eden/lua/ludoc"
	luhelp "github.com/BigJk/end_of_eden/lua/luhelp"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// ResourcesManager can load Artifacts, Cards, Events, Enemy and StoryTeller data from lua.
// The manager will walk the ./scripts directory and evaluate all found .lua files.
type ResourcesManager struct {
	Artifacts     map[string]*Artifact
	Cards         map[string]*Card
	Events        map[string]*Event
	Enemies       map[string]*Enemy
	StatusEffects map[string]*StatusEffect
	StoryTeller   map[string]*StoryTeller

	luaState   *lua.LState
	luaDocs    *ludoc.Docs
	log        *log.Logger
	registered *lua.LTable
	mapper     *luhelp.Mapper
}

func NewResourcesManager(state *lua.LState, docs *ludoc.Docs, logger *log.Logger) *ResourcesManager {
	man := &ResourcesManager{
		log:           logger,
		luaState:      state,
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
	man.luaState.SetGlobal("registered", man.registered)

	// Attach all register methods
	man.luaState.SetGlobal("register_artifact", man.luaState.NewFunction(man.luaRegisterArtifact))
	man.luaState.SetGlobal("register_card", man.luaState.NewFunction(man.luaRegisterCard))
	man.luaState.SetGlobal("register_enemy", man.luaState.NewFunction(man.luaRegisterEnemy))
	man.luaState.SetGlobal("register_event", man.luaState.NewFunction(man.luaRegisterEvent))
	man.luaState.SetGlobal("register_status_effect", man.luaState.NewFunction(man.luaRegisterStatusEffect))
	man.luaState.SetGlobal("register_story_teller", man.luaState.NewFunction(man.luaRegisterStoryTeller))
	man.luaState.SetGlobal("delete_event", man.luaState.NewFunction(man.luaDeleteEvent))
	man.defineDocs(docs)

	// Load all local scripts
	_ = filepath.Walk("./assets/scripts", func(path string, info fs.FileInfo, err error) error {
		// Don't load libs
		if strings.Contains(path, "scripts/libs") {
			return nil
		}

		if err != nil {
			return nil
		}

		if !info.IsDir() && strings.HasSuffix(path, ".lua") {
			if err := man.luaState.DoFile(path); err != nil {
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
		man.log.Println("Error while luaRegisterArtifact:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered artifact:", def.ID, def.Name)

	man.Artifacts[def.ID] = &def
	man.registered.RawGetString("artifact").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterCard(l *lua.LState) int {
	def := Card{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterCard:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered card:", def.ID, def.Name)

	man.Cards[def.ID] = &def
	man.registered.RawGetString("card").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterEnemy(l *lua.LState) int {
	def := Enemy{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterEnemy:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered enemy:", def.ID, def.Name)

	man.Enemies[def.ID] = &def
	man.registered.RawGetString("enemy").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterEvent(l *lua.LState) int {
	def := Event{}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterEvent:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered event:", def.ID, def.Name)

	man.Events[def.ID] = &def
	man.registered.RawGetString("event").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterStatusEffect(l *lua.LState) int {
	def := StatusEffect{
		Callbacks: map[string]luhelp.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterStatusEffect:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered status_effect:", def.ID, def.Name)

	man.StatusEffects[def.ID] = &def
	man.registered.RawGetString("status_effect").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaRegisterStoryTeller(l *lua.LState) int {
	def := StoryTeller{}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterStoryTeller:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered story_teller:", def.ID)

	man.StoryTeller[def.ID] = &def
	man.registered.RawGetString("story_teller").(*lua.LTable).RawSetString(def.ID, l.ToTable(2))
	return 0
}

func (man *ResourcesManager) luaDeleteEvent(l *lua.LState) int {
	delete(man.Events, l.ToString(1))
	man.registered.RawGetString("event").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) defineDocs(docs *ludoc.Docs) {
	if docs == nil {
		return
	}

	docs.Category("Content Registry", "These functions are used to define new content in the base game and in mods.", 100)

	docs.Function("register_artifact", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)`), "", "id : String", "definition : Table")

	docs.Function("register_card", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_card("MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 5 (+3 for each upgrade) damage.",
        state = function(ctx)
            return "Use your bare hands to deal " .. highlight(5 + ctx.level * 3) .. " damage."
        end,
        max_level = 1,
        color = "#2f3e46",
        need_target = true,
        point_cost = 1,
        price = 30,
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 5 + ctx.level * 3)
                return nil
            end,
        }
    }
)`), "", "id : String", "definition : Table")

	docs.Function("register_enemy", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_enemy("RUST_MITE",
    {
        name = "Rust Mite",
        description = "Loves to eat metal.",
        look = "/v\\",
        color = "#e6e65a",
        initial_hp = 22,
        max_hp = 22,
        gold = 10,
        callbacks = {
            on_turn = function(ctx)
                if ctx.round % 4 == 0 then
                    give_status_effect("RITUAL", ctx.guid)
                else
                    deal_damage(ctx.guid, PLAYER_ID, 6)
                end

                return nil
            end
        }
    }
)`), "", "id : String", "definition : Table")

	docs.Function("register_event", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_event("SOME_EVENT",
	{
		name = "Event Name",
		description = [[Flavor Text... Can include **Markdown** Syntax!]],
		choices = {
			{
				description = "Go...",
				callback = function()
					-- If you return nil on_end will decide the next game state
					return nil 
				end
			},
			{
				description = "Other Option",
				callback = function() return GAME_STATE_FIGHT end
			}
		},
		on_enter = function()
			play_music("energetic_orthogonal_expansions")
	
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("RUPTURE", PLAYER_ID)
			give_card("BLOCK", PLAYER_ID)
			give_artifact(get_random_artifact_type(150), PLAYER_ID)
		end,
		on_end = function(choice)
			-- Choice will be nil or the index of the choice taken
			return GAME_STATE_RANDOM
		end,
	}
)`), "", "id : String", "definition : Table")

	docs.Function("register_status_effect", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)`), "", "id : String", "definition : Table")

	docs.Function("register_story_teller", fmt.Sprintf("Registers a new artifact.\n\n```lua\n%s\n```", `register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)`), "", "id : String", "definition : Table")
}
