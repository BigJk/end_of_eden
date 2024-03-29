package game

import (
	"fmt"
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/internal/lua/ludoc"
	luhelp2 "github.com/BigJk/end_of_eden/internal/lua/luhelp"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"log"
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
	mapper     *luhelp2.Mapper
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
		mapper:     luhelp2.NewMapper(state),
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
	man.luaState.SetGlobal("delete_card", man.luaState.NewFunction(man.luaDeleteCard))
	man.luaState.SetGlobal("delete_enemy", man.luaState.NewFunction(man.luaDeleteEnemy))
	man.luaState.SetGlobal("delete_event", man.luaState.NewFunction(man.luaDeleteEvent))
	man.luaState.SetGlobal("delete_status_effect", man.luaState.NewFunction(man.luaDeleteStatusEffect))
	man.luaState.SetGlobal("delete_story_teller", man.luaState.NewFunction(man.luaDeleteStoryTeller))
	man.luaState.SetGlobal("delete_base_game", man.luaState.NewFunction(man.luaDeleteBaseGame))
	man.defineDocs(docs)

	// Load all local scripts
	_ = fs.Walk("./assets/scripts", func(path string, isDir bool) error {
		// Don't load libs, definitions and paths containing two __
		if strings.Contains(path, "scripts/libs") || strings.Contains(path, "scripts/definitions") || strings.Contains(path, "__") {
			return nil
		}

		if !isDir && strings.HasSuffix(path, ".lua") {
			luaBytes, err := fs.ReadFile(path)
			if err != nil {
				// TODO: error handling
				panic(err)
			}

			if strings.HasPrefix(string(luaBytes), "---@meta") {
				return nil
			}

			if err := man.luaState.DoString(string(luaBytes)); err != nil {
				// TODO: error handling
				panic(err)
			}
		}

		return nil
	})

	return man
}

// MarkBaseGame marks all currently registered resources as base game resources.
func (man *ResourcesManager) MarkBaseGame() {
	for _, v := range man.Artifacts {
		v.BaseGame = true
	}
	for _, v := range man.Cards {
		v.BaseGame = true
	}
	for _, v := range man.Events {
		v.BaseGame = true
	}
	for _, v := range man.Enemies {
		v.BaseGame = true
	}
	for _, v := range man.StatusEffects {
		v.BaseGame = true
	}
	for _, v := range man.StoryTeller {
		v.BaseGame = true
	}
}

func (man *ResourcesManager) luaRegisterArtifact(l *lua.LState) int {
	def := Artifact{
		Callbacks: map[string]luhelp2.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterArtifact:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered artifact:", def.ID, def.Name)

	man.Artifacts[def.ID] = &def

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("artifact").(*lua.LTable).RawSetString(def.ID, table)
	return 0
}

func (man *ResourcesManager) luaRegisterCard(l *lua.LState) int {
	def := Card{
		Callbacks: map[string]luhelp2.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterCard:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered card:", def.ID, def.Name)

	man.Cards[def.ID] = &def

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("card").(*lua.LTable).RawSetString(def.ID, table)
	return 0
}

func (man *ResourcesManager) luaRegisterEnemy(l *lua.LState) int {
	def := Enemy{
		Callbacks: map[string]luhelp2.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterEnemy:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered enemy:", def.ID, def.Name)

	man.Enemies[def.ID] = &def

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("enemy").(*lua.LTable).RawSetString(def.ID, table)
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

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("event").(*lua.LTable).RawSetString(def.ID, table)
	return 0
}

func (man *ResourcesManager) luaRegisterStatusEffect(l *lua.LState) int {
	def := StatusEffect{
		Callbacks: map[string]luhelp2.OwnedCallback{},
	}

	if err := man.mapper.Map(l.ToTable(2), &def); err != nil {
		man.log.Println("Error while luaRegisterStatusEffect:", err)
		return 0
	}

	// Set id after evaluating the table to avoid ID overwrite
	def.ID = l.ToString(1)
	man.log.Println("Registered status_effect:", def.ID, def.Name)

	man.StatusEffects[def.ID] = &def

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("status_effect").(*lua.LTable).RawSetString(def.ID, table)
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

	table := l.ToTable(2)
	l.SetTable(table, lua.LString("id"), lua.LString(def.ID))
	man.registered.RawGetString("story_teller").(*lua.LTable).RawSetString(def.ID, table)
	return 0
}

func (man *ResourcesManager) luaDeleteEvent(l *lua.LState) int {
	delete(man.Events, l.ToString(1))
	man.registered.RawGetString("event").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) luaDeleteCard(l *lua.LState) int {
	man.log.Println("Delete card:", l.ToString(1))

	delete(man.Cards, l.ToString(1))
	man.registered.RawGetString("card").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) luaDeleteEnemy(l *lua.LState) int {
	man.log.Println("Delete enemy:", l.ToString(1))

	delete(man.Enemies, l.ToString(1))
	man.registered.RawGetString("enemy").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) luaDeleteStatusEffect(l *lua.LState) int {
	man.log.Println("Delete status_effect:", l.ToString(1))

	delete(man.StatusEffects, l.ToString(1))
	man.registered.RawGetString("status_effect").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) luaDeleteStoryTeller(l *lua.LState) int {
	man.log.Println("Delete story_teller:", l.ToString(1))

	delete(man.StoryTeller, l.ToString(1))
	man.registered.RawGetString("story_teller").(*lua.LTable).RawSetString(l.ToString(1), lua.LNil)
	return 0
}

func (man *ResourcesManager) luaDeleteBaseGame(l *lua.LState) int {
	if l.GetTop() == 1 {
		t := l.ToString(1)
		switch t {
		case "artifact":
			man.Artifacts = lo.PickBy(man.Artifacts, func(k string, v *Artifact) bool { return !v.BaseGame })
		case "card":
			man.Cards = lo.PickBy(man.Cards, func(k string, v *Card) bool { return !v.BaseGame })
		case "enemy":
			man.Enemies = lo.PickBy(man.Enemies, func(k string, v *Enemy) bool { return !v.BaseGame })
		case "event":
			man.Events = lo.PickBy(man.Events, func(k string, v *Event) bool { return !v.BaseGame })
		case "status_effect":
			man.StatusEffects = lo.PickBy(man.StatusEffects, func(k string, v *StatusEffect) bool { return !v.BaseGame })
		case "story_teller":
			man.StoryTeller = lo.PickBy(man.StoryTeller, func(k string, v *StoryTeller) bool { return !v.BaseGame })
		}
		return 0
	}

	man.Artifacts = lo.PickBy(man.Artifacts, func(k string, v *Artifact) bool { return !v.BaseGame })
	man.Cards = lo.PickBy(man.Cards, func(k string, v *Card) bool { return !v.BaseGame })
	man.Enemies = lo.PickBy(man.Enemies, func(k string, v *Enemy) bool { return !v.BaseGame })
	man.Events = lo.PickBy(man.Events, func(k string, v *Event) bool { return !v.BaseGame })
	man.StatusEffects = lo.PickBy(man.StatusEffects, func(k string, v *StatusEffect) bool { return !v.BaseGame })
	man.StoryTeller = lo.PickBy(man.StoryTeller, func(k string, v *StoryTeller) bool { return !v.BaseGame })

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
)`), "", "id : type_id", "definition : artifact")

	docs.Function("register_card", fmt.Sprintf("Registers a new card.\n\n```lua\n%s\n```", `register_card("MELEE_HIT",
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
)`), "", "id : type_id", "definition : card")

	docs.Function("register_enemy", fmt.Sprintf("Registers a new enemy.\n\n```lua\n%s\n```", `register_enemy("RUST_MITE",
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
)`), "", "id : type_id", "definition : enemy")

	docs.Function("register_event", fmt.Sprintf("Registers a new event.\n\n```lua\n%s\n```", `register_event("SOME_EVENT",
	{
		name = "Event Name",
		description = "Flavor Text... Can include **Markdown** Syntax!",
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
)`), "", "id : type_id", "definition : event")

	docs.Function("register_status_effect", fmt.Sprintf("Registers a new status effect.\n\n```lua\n%s\n```", `register_status_effect("BLOCK", {
    name = "Block",
    description = "Decreases incoming damage for each stack",
    look = "Blk",
    foreground = "#219ebc",
    state = function(ctx)
        return "Takes " .. highlight(ctx.stacks) .. " less damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                add_status_effect_stacks(ctx.guid, -ctx.damage)
                return ctx.damage - ctx.stacks
            end
            return ctx.damage
        end
    }
})`), "", "id : type_id", "definition : status_effect")

	docs.Function("register_story_teller", fmt.Sprintf("Registers a new story teller.\n\n```lua\n%s\n```", `register_story_teller("STORY_TELLER_XYZ", {
    active = function(ctx)
        if not had_events_any({ "A", "B", "C" }) then
            return 1
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage >= 3 then
            set_event("SOME_EVENT")
            return GAME_STATE_EVENT
        end

        -- Fight against rust mites or clean bots
        local d = math.random(2)
        if d == 1 then
            add_actor_by_enemy("RUST_MITE")
        elseif d == 2 then
            add_actor_by_enemy("CLEAN_BOT")
        end

        return GAME_STATE_FIGHT
    end
})`), "", "id : type_id", "definition : story_teller")

	docs.Function("delete_event", fmt.Sprintf("Deletes an event.\n\n```lua\n%s\n```", `delete_event("SOME_EVENT")`), "", "id : type_id")
	docs.Function("delete_card", fmt.Sprintf("Deletes a card.\n\n```lua\n%s\n```", `delete_card("SOME_CARD")`), "", "id : type_id")
	docs.Function("delete_enemy", fmt.Sprintf("Deletes an enemy.\n\n```lua\n%s\n```", `delete_enemy("SOME_ENEMY")`), "", "id : type_id")
	docs.Function("delete_status_effect", fmt.Sprintf("Deletes a status effect.\n\n```lua\n%s\n```", `delete_status_effect("SOME_STATUS_EFFECT")`), "", "id : type_id")
	docs.Function("delete_story_teller", fmt.Sprintf("Deletes a story teller.\n\n```lua\n%s\n```", `delete_story_teller("SOME_STORY_TELLER")`), "", "id : type_id")

	docs.Function("delete_base_game", fmt.Sprintf("Deletes all base game content. Useful if you don't want to include base game content in your mod.\n\n```lua\n%s\n```", `delete_base_game() -- delete all base game content
delete_base_game("artifact") -- deletes all artifacts
delete_base_game("card") -- deletes all cards
delete_base_game("enemy") -- deletes all enemies
delete_base_game("event") -- deletes all events
delete_base_game("status_effect") -- deletes all status effects
delete_base_game("story_teller") -- deletes all story tellers
`), "", "(optional) type : string")
}
