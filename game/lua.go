package game

import (
	"github.com/BigJk/end_of_eden/internal/fs"
	"github.com/BigJk/end_of_eden/internal/lua/ludoc"
	luhelp2 "github.com/BigJk/end_of_eden/internal/lua/luhelp"
	"github.com/BigJk/end_of_eden/system/audio"
	"github.com/BigJk/end_of_eden/system/gen/faces"
	"github.com/BigJk/end_of_eden/system/localization"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
)

// SessionAdapter creates a lua vm that is bound to the session in the given Session.
func SessionAdapter(session *Session) (*lua.LState, *ludoc.Docs) {
	l := lua.NewState(lua.Options{
		IncludeGoStackTrace: true,
	})
	d := ludoc.New()

	mapper := luhelp2.NewMapper(l)

	_ = fs.Walk("./assets/scripts/libs", func(path string, isDir bool) error {
		if isDir || !strings.HasSuffix(path, ".lua") {
			return nil
		}

		name := strings.Split(filepath.Base(path), ".")[0]

		luaBytes, err := fs.ReadFile(path)
		if err != nil {
			return err
		}

		mod, err := l.LoadString(string(luaBytes))
		if err != nil {
			session.log.Println("Can't LoadFile module:", path)
			return nil
		}

		session.log.Println("Loaded lib:", path, name)

		preload := l.GetField(l.GetField(l.Get(lua.EnvironIndex), "package"), "preload")
		l.SetField(preload, name, mod)

		return nil
	})

	// Require fun by default

	_ = l.DoString(`
require("fun")()
fun = require "fun"
`)

	// Constants

	d.Category("Game Constants", "General game constants.", 0)

	d.Global("PLAYER_ID", "Player actor id for use in functions where the guid is needed, for example: ``deal_damage(PLAYER_ID, enemy_guid, 10)``.")

	l.SetGlobal("PLAYER_ID", lua.LString(PlayerActorID))

	d.Global("GAME_STATE_FIGHT", "Represents the fight game state.")
	d.Global("GAME_STATE_EVENT", "Represents the event game state.")
	d.Global("GAME_STATE_MERCHANT", "Represents the merchant game state.")
	d.Global("GAME_STATE_RANDOM", "Represents the random game state in which the active story teller will decide what happens next.")

	l.SetGlobal("GAME_STATE_FIGHT", lua.LString(GameStateFight))
	l.SetGlobal("GAME_STATE_EVENT", lua.LString(GameStateEvent))
	l.SetGlobal("GAME_STATE_MERCHANT", lua.LString(GameStateMerchant))
	l.SetGlobal("GAME_STATE_RANDOM", lua.LString(GameStateRandom))

	d.Global("DECAY_ONE", "Status effect decays by 1 stack per turn.")
	d.Global("DECAY_ALL", "Status effect decays by all stacks per turn.")
	d.Global("DECAY_NONE", "Status effect never decays.")

	l.SetGlobal("DECAY_ONE", lua.LString(DecayOne))
	l.SetGlobal("DECAY_ALL", lua.LString(DecayAll))
	l.SetGlobal("DECAY_NONE", lua.LString(DecayNone))

	// Utility

	d.Category("Utility", "General game constants.", 1)

	d.Function("guid", "returns a new random guid.", "guid")
	l.SetGlobal("guid", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(NewGuid("LUA")))
		return 1
	}))

	d.Function("store", "Stores a persistent value for this run that will be restored after a save load. Can store any lua basic value or table.", "", "key : string", "value : any")
	l.SetGlobal("store", l.NewFunction(func(state *lua.LState) int {
		session.Store(state.ToString(1), mapper.ToGoValue(state.Get(2)))
		return 0
	}))

	d.Function("fetch", "Fetches a value from the persistent store", "any", "key : string")
	l.SetGlobal("fetch", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.Fetch(state.ToString(1))))
		return 1
	}))

	// Style

	d.Category("Styling", "Helper functions for text styling.", 2)

	d.Function("text_bold", "Makes the text bold.", "string", "value : any")
	l.SetGlobal("text_bold", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[1m" + luhelp2.ToString(state.Get(1), mapper) + "\033[22m"))
		return 1
	}))

	d.Function("text_italic", "Makes the text italic.", "string", "value : any")
	l.SetGlobal("text_italic", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[3m" + luhelp2.ToString(state.Get(1), mapper) + "\033[23m"))
		return 1
	}))

	d.Function("text_underline", "Makes the text underlined.", "string", "value : any")
	l.SetGlobal("text_underline", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[4m" + luhelp2.ToString(state.Get(1), mapper) + "\033[24m"))
		return 1
	}))

	d.Function("text_red", "Makes the text colored red.", "string", "value : any")
	l.SetGlobal("text_red", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\x1b[38;5;9m" + luhelp2.ToString(state.Get(1), mapper)))
		return 1
	}))

	d.Function("text_bg", "Makes the text background colored. Takes hex values like #ff0000.", "string", "color : string", "value : any")
	l.SetGlobal("text_bg", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(removeAnsiReset(lipgloss.NewStyle().Background(lipgloss.Color(luhelp2.ToString(state.Get(1), mapper))).Render(luhelp2.ToString(state.Get(2), mapper)))))
		return 1
	}))

	// Misc

	d.Category("Logging", "Various logging functions.", 3)

	d.Function("log_i", "Log at **information** level to player log.", "", "value : any")
	l.SetGlobal("log_i", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeInfo, luhelp2.ToString(state.Get(1), mapper))
		return 0
	}))

	d.Function("log_w", "Log at **warning** level to player log.", "", "value : any")
	l.SetGlobal("log_w", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeWarning, luhelp2.ToString(state.Get(1), mapper))
		return 0
	}))

	d.Function("log_d", "Log at **danger** level to player log.", "", "value : any")
	l.SetGlobal("log_d", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeDanger, luhelp2.ToString(state.Get(1), mapper))
		return 0
	}))

	d.Function("log_s", "Log at **success** level to player log.", "", "value : any")
	l.SetGlobal("log_s", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeSuccess, luhelp2.ToString(state.Get(1), mapper))
		return 0
	}))

	l.SetGlobal("breakpoint", l.NewFunction(func(state *lua.LState) int {
		val := lo.Map(make([]lua.LValue, state.GetTop()), func(_ lua.LValue, index int) lua.LValue {
			val := state.Get(1 + index)
			return val
		})
		session.log.Println(val)
		return 0
	}))

	l.SetGlobal("debug_log", l.NewFunction(func(state *lua.LState) int {
		dbg, ok := state.GetStack(1)
		if ok {
			_, _ = state.GetInfo("nSl", dbg, lua.LNil)
		}

		session.log.Printf("[LUA :: %d %s] %s \n", dbg.CurrentLine, dbg.Source, strings.Join(lo.Map(make([]any, state.GetTop()), func(_ any, index int) string {
			val := state.Get(1 + index)
			return luhelp2.ToString(val, mapper)
		}), " "))

		return 0
	}))

	d.Function("print", "Log to session log.", "", "...")
	if err := l.DoString("print = debug_log"); err != nil {
		panic("Can't overwrite print with debug_log")
	}

	// Audio

	d.Category("Audio", "Audio helper functions.", 4)

	d.Function("play_audio", "Plays a sound effect. If you want to play ``button.mp3`` you call ``play_audio(\"button\")``.", "", "sound : string")
	l.SetGlobal("play_audio", l.NewFunction(func(state *lua.LState) int {
		audio.Play(state.ToString(1))
		return 0
	}))

	d.Function("play_music", "Start a song for the background loop. If you want to play ``song.mp3`` you call ``play_music(\"song\")``.", "", "sound : string")
	l.SetGlobal("play_music", l.NewFunction(func(state *lua.LState) int {
		audio.PlayMusic(state.ToString(1))
		return 0
	}))

	// Game State

	d.Category("Game State", "Functions that modify the general game state.", 5)

	d.Function("set_event", "Set event by id.", "", "event_id : type_id")
	l.SetGlobal("set_event", l.NewFunction(func(state *lua.LState) int {
		session.SetEvent(state.ToString(1))
		return 0
	}))

	d.Function("set_game_state", "Set the current game state. See globals.", "", "state : next_game_state")
	l.SetGlobal("set_game_state", l.NewFunction(func(state *lua.LState) int {
		session.SetGameState(GameState(state.ToString(1)))
		return 0
	}))

	d.Function("set_fight_description", "Set the current fight description. This will be shown on the top right in the game.", "", "desc : string")
	l.SetGlobal("set_fight_description", l.NewFunction(func(state *lua.LState) int {
		session.SetFightDescription(state.ToString(1))
		return 0
	}))

	d.Function("get_fight_round", "Gets the fight round.", "number")
	l.SetGlobal("get_fight_round", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetFightRound()))
		return 1
	}))

	d.Function("get_stages_cleared", "Gets the number of stages cleared.", "number")
	l.SetGlobal("get_stages_cleared", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetStagesCleared()))
		return 1
	}))

	d.Function("get_fight", "Gets the fight state. This contains the player hand, used, exhausted and round information.", "fight_state")
	l.SetGlobal("get_fight", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetFight()))
		return 1
	}))

	d.Function("get_event_history", "Gets the ids of all the encountered events in the order of occurrence.", "string[]")
	l.SetGlobal("get_event_history", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetEventHistory()))
		return 1
	}))

	d.Function("had_event", "Checks if the event happened at least once.", "boolean", "event_id : type_id")
	l.SetGlobal("had_event", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.HadEvent(state.ToString(1))))
		return 1
	}))

	d.Function("had_events", "Checks if all the events happened at least once.", "boolean", "event_ids : type_id[]")
	l.SetGlobal("had_events", l.NewFunction(func(state *lua.LState) int {
		var ids []string
		if err := mapper.Map(state.Get(1).(*lua.LTable), &ids); err != nil {
			session.logLuaError("had_event", "", err)
			return 0
		} else {
			state.Push(luhelp2.ToLua(state, session.HadEvents(ids)))
		}
		return 1
	}))

	d.Function("had_events_any", "Checks if any of the events happened at least once.", "boolean", "eventIds : string[]")
	l.SetGlobal("had_events_any", l.NewFunction(func(state *lua.LState) int {
		var ids []string
		if err := mapper.Map(state.Get(1).(*lua.LTable), &ids); err != nil {
			session.logLuaError("had_events_any", "", err)
			return 0
		} else {
			state.Push(luhelp2.ToLua(state, session.HadEventsAny(ids)))
		}
		return 1
	}))

	// Actor Operations

	d.Category("Actor Operations", "Functions that modify or access the actors. Actors are either the player or enemies.", 6)

	d.Function("get_player", "Get the player actor. Equivalent to ``get_actor(PLAYER_ID)``", "actor")
	l.SetGlobal("get_player", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetPlayer()))
		return 1
	}))

	d.Function("get_actor", "Get a actor by guid.", "actor", "guid : guid")
	l.SetGlobal("get_actor", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetActor(state.ToString(1))))
		return 1
	}))

	d.Function("get_opponent_by_index", "Get opponent (actor) by index of a certain actor. ``get_opponent_by_index(PLAYER_ID, 2)`` would return the second alive opponent of the player.", "actor", "guid : guid", "index : number")
	l.SetGlobal("get_opponent_by_index", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetOpponentByIndex(state.ToString(1), int(state.ToNumber(2))-1)))
		return 1
	}))

	d.Function("get_opponent_count", "Get the number of opponents (actors) of a certain actor. ``get_opponent_count(PLAYER_ID)`` would return 2 if the player had 2 alive enemies.", "number", "guid : guid")
	l.SetGlobal("get_opponent_count", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetOpponentCount(state.ToString(1))))
		return 1
	}))

	d.Function("get_opponent_guids", "Get the guids of opponents (actors) of a certain actor. If the player had 2 enemies, ``get_opponent_guids(PLAYER_ID)`` would return a table with 2 strings containing the guids of these actors.", "guid[]", "guid : guid")
	l.SetGlobal("get_opponent_guids", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetOpponentGUIDs(state.ToString(1))))
		return 1
	}))

	d.Function("remove_actor", "Deletes a actor by id.", "", "guid : guid")
	l.SetGlobal("remove_actor", l.NewFunction(func(state *lua.LState) int {
		session.GetActor(state.ToString(1))
		return 0
	}))

	d.Function("actor_add_max_hp", "Increases the max hp value of a actor by a number. Can be negative value to decrease it.", "", "guid : guid", "amount : number")
	l.SetGlobal("actor_add_max_hp", l.NewFunction(func(state *lua.LState) int {
		session.ActorAddMaxHP(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	d.Function("actor_set_max_hp", "Sets the max hp value of a actor to a number.", "", "guid : guid", "amount : number")
	l.SetGlobal("actor_set_max_hp", l.NewFunction(func(state *lua.LState) int {
		session.UpdateActor(state.ToString(1), func(actor *Actor) bool {
			actor.MaxHP = int(state.ToNumber(2))
			return true
		})
		return 0
	}))

	d.Function("actor_add_hp", "Increases the hp value of a actor by a number. Can be negative value to decrease it. This won't trigger any on_damage callbacks", "", "guid : guid", "amount : number")
	l.SetGlobal("actor_add_hp", l.NewFunction(func(state *lua.LState) int {
		session.ActorAddHP(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	d.Function("actor_set_hp", "Sets the hp value of a actor to a number. This won't trigger any on_damage callbacks", "", "guid : guid", "amount : number")
	l.SetGlobal("actor_set_hp", l.NewFunction(func(state *lua.LState) int {
		session.UpdateActor(state.ToString(1), func(actor *Actor) bool {
			actor.HP = int(state.ToNumber(2))
			return true
		})
		return 0
	}))

	d.Function("add_actor_by_enemy", "Creates a new enemy fighting against the player. Example ``add_actor_by_enemy(\"RUST_MITE\")``.", "string", "enemy_guid : type_id")
	l.SetGlobal("add_actor_by_enemy", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.AddActorFromEnemy(state.ToString(1))))
		return 1
	}))

	// Artifacts

	d.Category("Artifact Operations", "Functions that modify or access the artifacts.", 7)

	d.Function("give_artifact", "Gives a actor a artifact. Returns the guid of the newly created artifact.", "string", "type_id : type_id", "actor : guid")
	l.SetGlobal("give_artifact", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GiveArtifact(state.ToString(1), state.ToString(2))))
		return 1
	}))

	d.Function("remove_artifact", "Removes a artifact.", "", "guid : guid")
	l.SetGlobal("remove_artifact", l.NewFunction(func(state *lua.LState) int {
		session.RemoveArtifact(state.ToString(1))
		return 0
	}))

	d.Function("get_artifacts", "Returns all the artifacts guids from the given actor.", "guid[]", "actor_guid : string")
	l.SetGlobal("get_artifacts", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetArtifacts(state.ToString(1))))
		return 1
	}))

	d.Function("get_artifact", "Returns the artifact definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.", "artifact", "id : string")
	l.SetGlobal("get_artifact", l.NewFunction(func(state *lua.LState) int {
		art, _ := session.GetArtifact(state.ToString(1))
		state.Push(luhelp2.ToLua(state, art))
		return 1
	}))

	d.Function("get_artifact_instance", "Returns the artifact instance by guid.", "artifact_instance", "guid : guid")
	l.SetGlobal("get_artifact_instance", l.NewFunction(func(state *lua.LState) int {
		_, instance := session.GetArtifact(state.ToString(1))
		state.Push(luhelp2.ToLua(state, instance))
		return 1
	}))

	// Status Effects

	d.Category("Status Effect Operations", "Functions that modify or access the status effects.", 8)

	d.Function("give_status_effect", "Gives a status effect to a actor. If count is not specified a stack of 1 is applied.", "", "type_id : string", "actor_guid : string", "(optional) count : number")
	l.SetGlobal("give_status_effect", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 2 {
			state.Push(lua.LString(session.GiveStatusEffect(state.ToString(1), state.ToString(2), 1)))
		} else {
			state.Push(lua.LString(session.GiveStatusEffect(state.ToString(1), state.ToString(2), int(state.ToNumber(3)))))
		}
		return 1
	}))

	d.Function("remove_status_effect", "Removes a status effect.", "", "guid : guid")
	l.SetGlobal("remove_status_effect", l.NewFunction(func(state *lua.LState) int {
		session.RemoveStatusEffect(state.ToString(1))
		return 0
	}))

	d.Function("add_status_effect_stacks", "Adds to the stack count of a status effect. Negative values are also allowed.", "", "guid : guid", "count : number")
	l.SetGlobal("add_status_effect_stacks", l.NewFunction(func(state *lua.LState) int {
		session.AddStatusEffectStacks(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	d.Function("set_status_effect_stacks", "Sets the stack count of a status effect by guid.", "", "guid : guid", "count : number")
	l.SetGlobal("set_status_effect_stacks", l.NewFunction(func(state *lua.LState) int {
		session.SetStatusEffectStacks(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	d.Function("get_actor_status_effects", "Returns the guids of all status effects that belong to a actor.", "guid[]", "actor_guid : string")
	l.SetGlobal("get_actor_status_effects", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetActorStatusEffects(state.ToString(1))))
		return 1
	}))

	d.Function("get_status_effect", "Returns the status effect definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.", "status_effect", "id : string")
	l.SetGlobal("get_status_effect", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetStatusEffect(state.ToString(1))))
		return 1
	}))

	d.Function("get_status_effect_instance", "Returns the status effect instance.", "status_effect_instance", "effect_guid : guid")
	l.SetGlobal("get_status_effect_instance", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetStatusEffectInstance(state.ToString(1))))
		return 1
	}))

	// Cards

	d.Category("Card Operations", "Functions that modify or access the cards.", 9)

	d.Function("give_card", "Gives a card.", "string", "card_type_id : type_id", "owner_actor_guid : guid")
	l.SetGlobal("give_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GiveCard(state.ToString(1), state.ToString(2))))
		return 1
	}))

	d.Function("remove_card", "Removes a card.", "", "card_guid : string")
	l.SetGlobal("remove_card", l.NewFunction(func(state *lua.LState) int {
		session.RemoveCard(state.ToString(1))
		return 0
	}))

	d.Function("cast_card", "Tries to cast a card with a guid and optional target. If the cast isn't successful returns false.", "boolean", "card_guid : guid", "(optional) target_actor_guid : guid")
	l.SetGlobal("cast_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LBool(session.CastCard(state.ToString(1), state.ToString(2))))
		return 1
	}))

	d.Function("get_cards", "Returns all the card guids from the given actor.", "guid[]", "actor_guid : string")
	l.SetGlobal("get_cards", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetCards(state.ToString(1))))
		return 1
	}))

	d.Function("get_card", "Returns the card type definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.", "card", "id : type_id")
	l.SetGlobal("get_card", l.NewFunction(func(state *lua.LState) int {
		card, _ := session.GetCard(state.ToString(1))
		state.Push(luhelp2.ToLua(state, card))
		return 1
	}))

	d.Function("get_card_instance", "Returns the instance object of a card.", "card_instance", "card_guid : guid")
	l.SetGlobal("get_card_instance", l.NewFunction(func(state *lua.LState) int {
		_, instance := session.GetCard(state.ToString(1))
		state.Push(luhelp2.ToLua(state, instance))
		return 1
	}))

	d.Function("upgrade_card", "Upgrade a card without paying for it.", "boolean", "card_guid : guid")
	l.SetGlobal("upgrade_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LBool(session.UpgradeCard(state.ToString(1))))
		return 1
	}))

	d.Function("upgrade_random_card", "Upgrade a random card without paying for it.", "boolean", "actor_guid : guid")
	l.SetGlobal("upgrade_random_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LBool(session.UpgradeRandomCard(state.ToString(1))))
		return 1
	}))

	// Damage & Heal

	d.Category("Damage & Heal", "Functions that deal damage or heal.", 10)

	d.Function("deal_damage", "Deal damage from one source to a target. If flat is true the damage can't be modified by status effects or artifacts. Returns the damage that was dealt.", "number", "source : guid", "target : guid", "damage : number", "(optional) flat : boolean")
	l.SetGlobal("deal_damage", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 3 {
			state.Push(lua.LNumber(session.DealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), false)))
		} else {
			state.Push(lua.LNumber(session.DealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	d.Function("simulate_deal_damage", "Simulate damage from a source to a target. If flat is true the damage can't be modified by status effects or artifacts. Returns the damage that would be dealt.", "number", "source : guid", "target : guid", "damage : number", "(optional) flat : boolean")
	l.SetGlobal("simulate_deal_damage", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 3 {
			state.Push(lua.LNumber(session.SimulateDealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), false)))
		} else {
			state.Push(lua.LNumber(session.SimulateDealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	d.Function("deal_damage_multi", "Deal damage to multiple enemies from one source. If flat is true the damage can't be modified by status effects or artifacts. Returns a array of damages for each actor hit.", "number[]", "source : guid", "targets : guid[]", "damage : number", "(optional) flat : boolean")
	l.SetGlobal("deal_damage_multi", l.NewFunction(func(state *lua.LState) int {
		var guids []string

		switch state.Get(2).Type() {
		case lua.LTTable:
			if err := mapper.Map(state.Get(2).(*lua.LTable), &guids); err != nil {
				session.log.Printf("Error in deal_damage_multi: %v\n", err)
				return 0
			}
		case lua.LTUserData:
			if val, ok := state.Get(2).(*lua.LUserData).Value.([]string); ok {
				guids = val
			}
		default:
			session.log.Printf("Error in deal_damage_multi: wrong type %v", state.Get(2).Type().String())
			return 0
		}

		if state.GetTop() == 3 {
			state.Push(luhelp2.ToLua(state, session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), false)))
		} else {
			state.Push(luhelp2.ToLua(state, session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	d.Function("heal", "Heals the target triggered by the source.", "", "source : guid", "target : guid", "amount : number")
	l.SetGlobal("heal", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 3 {
			state.Push(lua.LNumber(session.Heal(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), false)))
		} else {
			state.Push(lua.LNumber(session.Heal(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	// Player

	d.Category("Player Operations", "Functions that are related to the player.", 11)

	d.Function("player_draw_card", "Let the player draw additional cards for this turn.", "", "amount : number")
	l.SetGlobal("player_draw_card", l.NewFunction(func(state *lua.LState) int {
		session.PlayerDrawCard(int(state.ToNumber(1)))
		return 0
	}))

	d.Function("player_give_action_points", "Gives the player more action points for this turn.", "", "points : number")
	l.SetGlobal("player_give_action_points", l.NewFunction(func(state *lua.LState) int {
		session.PlayerGiveActionPoints(int(state.ToNumber(1)))
		return 0
	}))

	d.Function("give_player_gold", "Gives the player gold.", "", "amount : number")
	l.SetGlobal("give_player_gold", l.NewFunction(func(state *lua.LState) int {
		session.GivePlayerGold(int(state.ToNumber(1)))
		return 0
	}))

	d.Function("player_buy_card", "Let the player buy the card with the given id. This will deduct the price form the players gold and return true if the buy was successful.", "boolean", "card_id : type_id")
	l.SetGlobal("player_buy_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LBool(session.PlayerBuyCard(state.ToString(1))))
		return 1
	}))

	d.Function("player_buy_artifact", "Let the player buy the artifact with the given id. This will deduct the price form the players gold and return true if the buy was successful.", "boolean", "card_id : type_id")
	l.SetGlobal("player_buy_artifact", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LBool(session.PlayerBuyArtifact(state.ToString(1))))
		return 1
	}))

	d.Function("finish_player_turn", "Finishes the player turn.", "")
	l.SetGlobal("finish_player_turn", l.NewFunction(func(state *lua.LState) int {
		session.FinishPlayerTurn()
		return 0
	}))

	// Merchant

	d.Category("Merchant Operations", "Functions that are related to the merchant.", 12)

	d.Function("get_merchant", "Returns the merchant state.", "merchant_state")
	l.SetGlobal("get_merchant", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp2.ToLua(state, session.GetMerchant()))
		return 1
	}))

	d.Function("add_merchant_card", "Adds another random card to the merchant", "")
	l.SetGlobal("add_merchant_card", l.NewFunction(func(state *lua.LState) int {
		session.AddMerchantCard()
		return 0
	}))

	d.Function("add_merchant_artifact", "Adds another random artifact to the merchant", "")
	l.SetGlobal("add_merchant_artifact", l.NewFunction(func(state *lua.LState) int {
		session.AddMerchantArtifact()
		return 0
	}))

	d.Function("get_merchant_gold_max", "Returns the maximum value of artifacts and cards that the merchant will sell. Good to scale ``random_card`` and ``random_artifact``.", "number")
	l.SetGlobal("get_merchant_gold_max", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetMerchantGoldMax()))
		return 1
	}))

	// Random

	d.Category("Random Utility", "Functions that help with random generation.", 13)

	d.Function("gen_face", "Generates a random face.", "string", "(optional) category : number")
	l.SetGlobal("gen_face", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 1 {
			state.Push(lua.LString(faces.Global.Gen(int(state.ToNumber(1)))))
		} else {
			state.Push(lua.LString(faces.Global.GenRand()))
		}
		return 1
	}))

	d.Function("random_card", "Returns the type id of a random card.", "type_id", "max_price : number")
	l.SetGlobal("random_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GetRandomCard(int(state.ToNumber(1)))))
		return 1
	}))

	d.Function("random_artifact", "Returns the type id of a random artifact.", "type_id", "max_price : number")
	l.SetGlobal("random_artifact", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GetRandomArtifact(int(state.ToNumber(1)))))
		return 1
	}))

	// Localization

	d.Category("Localization", "Functions that help with localization.", 14)

	d.Function("l", "Returns the localized string for the given key. Examples on locals definition can be found in `/assets/locals`. Example: ``\nl('cards.MY_CARD.name', \"English Default Name\")``", "string", "key : string", "(optional) default : string")
	l.SetGlobal("l", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 1 {
			state.Push(lua.LString(localization.G(state.ToString(1))))
		} else {
			state.Push(lua.LString(localization.G(state.ToString(1), state.ToString(2))))
		}
		return 1
	}))

	return l, d
}

// removeAnsiReset removes the first ansi reset code from a string.
func removeAnsiReset(s string) string {
	return strings.Replace(s, "\x1b[0m", "", 1)
}
