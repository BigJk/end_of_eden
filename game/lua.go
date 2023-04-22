package game

import (
	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/gen/faces"
	"github.com/BigJk/project_gonzo/luhelp"
	"github.com/BigJk/project_gonzo/util"
	"github.com/charmbracelet/lipgloss"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"path/filepath"
	"strings"
)

// SessionAdapter creates a lua vm that is bound to the session in the given Session.
func SessionAdapter(session *Session) *lua.LState {
	l := lua.NewState()
	mapper := luhelp.NewMapper(l)

	_ = filepath.Walk("./assets/scripts/libs", func(path string, info fs.FileInfo, _ error) error {
		if info != nil && info.IsDir() || !strings.HasSuffix(path, ".lua") {
			return nil
		}

		name := strings.Split(filepath.Base(path), ".")[0]

		mod, err := l.LoadFile(path)
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

	l.SetGlobal("PLAYER_ID", lua.LString(PlayerActorID))

	l.SetGlobal("GAME_STATE_FIGHT", lua.LString(GameStateFight))
	l.SetGlobal("GAME_STATE_EVENT", lua.LString(GameStateEvent))
	l.SetGlobal("GAME_STATE_MERCHANT", lua.LString(GameStateMerchant))
	l.SetGlobal("GAME_STATE_RANDOM", lua.LString(GameStateRandom))

	l.SetGlobal("DECAY_ONE", lua.LString(DecayOne))
	l.SetGlobal("DECAY_ALL", lua.LString(DecayAll))
	l.SetGlobal("DECAY_NONE", lua.LString(DecayNone))

	// Utility

	l.SetGlobal("guid", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(NewGuid("LUA")))
		return 1
	}))

	// Style

	l.SetGlobal("text_bold", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[1m" + luhelp.ToString(state.Get(1), mapper) + "\033[22m"))
		return 1
	}))

	l.SetGlobal("text_italic", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[3m" + luhelp.ToString(state.Get(1), mapper) + "\033[23m"))
		return 1
	}))

	l.SetGlobal("text_underline", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString("\033[4m" + luhelp.ToString(state.Get(1), mapper) + "\033[24m"))
		return 1
	}))

	l.SetGlobal("text_color", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(util.RemoveAnsiReset(lipgloss.NewStyle().Foreground(lipgloss.Color(luhelp.ToString(state.Get(1), mapper))).Render(luhelp.ToString(state.Get(2), mapper)))))
		return 1
	}))

	l.SetGlobal("text_bg", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(util.RemoveAnsiReset(lipgloss.NewStyle().Background(lipgloss.Color(luhelp.ToString(state.Get(1), mapper))).Render(luhelp.ToString(state.Get(2), mapper)))))
		return 1
	}))

	// Misc

	l.SetGlobal("log_i", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeInfo, luhelp.ToString(state.Get(1), mapper))
		return 0
	}))

	l.SetGlobal("log_w", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeWarning, luhelp.ToString(state.Get(1), mapper))
		return 0
	}))

	l.SetGlobal("log_d", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeDanger, luhelp.ToString(state.Get(1), mapper))
		return 0
	}))

	l.SetGlobal("log_s", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeSuccess, luhelp.ToString(state.Get(1), mapper))
		return 0
	}))

	l.SetGlobal("debug_value", l.NewFunction(func(state *lua.LState) int {
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
			return luhelp.ToString(val, mapper)
		}), " "))

		return 0
	}))

	if err := l.DoString("print = debug_log"); err != nil {
		panic("Can't overwrite print with debug_log")
	}

	// Audio

	l.SetGlobal("play_audio", l.NewFunction(func(state *lua.LState) int {
		audio.Play(state.ToString(1))
		return 0
	}))

	// Game State

	l.SetGlobal("set_event", l.NewFunction(func(state *lua.LState) int {
		session.SetEvent(state.ToString(1))
		return 0
	}))

	l.SetGlobal("set_game_state", l.NewFunction(func(state *lua.LState) int {
		session.SetGameState(GameState(state.ToString(1)))
		return 0
	}))

	l.SetGlobal("set_fight_description", l.NewFunction(func(state *lua.LState) int {
		session.SetFightDescription(state.ToString(1))
		return 0
	}))

	l.SetGlobal("get_fight_round", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetFightRound()))
		return 1
	}))

	l.SetGlobal("get_stages_cleared", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetStagesCleared()))
		return 1
	}))

	l.SetGlobal("get_fight", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetFight()))
		return 1
	}))

	l.SetGlobal("get_event_history", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetEventHistory()))
		return 1
	}))

	l.SetGlobal("had_event", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.HadEvent(state.ToString(1))))
		return 1
	}))

	// Actor Operations

	l.SetGlobal("get_player", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetPlayer()))
		return 1
	}))

	l.SetGlobal("get_actor", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetActor(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_opponent_by_index", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetOpponentByIndex(state.ToString(1), int(state.ToNumber(2))-1)))
		return 1
	}))

	l.SetGlobal("get_opponent_count", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetOpponentCount(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_opponent_guids", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetOpponentGUIDs(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("remove_actor", l.NewFunction(func(state *lua.LState) int {
		session.GetActor(state.ToString(1))
		return 0
	}))

	l.SetGlobal("actor_add_max_hp", l.NewFunction(func(state *lua.LState) int {
		session.ActorAddMaxHP(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	l.SetGlobal("add_actor_by_enemy", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.AddActorFromEnemy(state.ToString(1))))
		return 1
	}))

	// Artifacts

	l.SetGlobal("give_artifact", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GiveArtifact(state.ToString(1), state.ToString(2))))
		return 1
	}))

	l.SetGlobal("remove_artifact", l.NewFunction(func(state *lua.LState) int {
		session.RemoveArtifact(state.ToString(1))
		return 0
	}))

	l.SetGlobal("get_random_artifact_type", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GetRandomArtifactType(int(state.ToNumber(1)))))
		return 1
	}))

	l.SetGlobal("get_artifact", l.NewFunction(func(state *lua.LState) int {
		art, _ := session.GetArtifact(state.ToString(1))
		state.Push(luhelp.ToLua(state, art))
		return 1
	}))

	l.SetGlobal("get_artifact_instance", l.NewFunction(func(state *lua.LState) int {
		_, instance := session.GetArtifact(state.ToString(1))
		state.Push(luhelp.ToLua(state, instance))
		return 1
	}))

	// Status Effects

	l.SetGlobal("give_status_effect", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 2 {
			state.Push(lua.LString(session.GiveStatusEffect(state.ToString(1), state.ToString(2), 1)))
		} else {
			state.Push(lua.LString(session.GiveStatusEffect(state.ToString(1), state.ToString(2), int(state.ToNumber(3)))))
		}
		return 1
	}))

	l.SetGlobal("remove_status_effect", l.NewFunction(func(state *lua.LState) int {
		session.RemoveStatusEffect(state.ToString(1))
		return 0
	}))

	l.SetGlobal("add_status_effect_stacks", l.NewFunction(func(state *lua.LState) int {
		session.AddStatusEffectStacks(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	l.SetGlobal("set_status_effect_stacks", l.NewFunction(func(state *lua.LState) int {
		session.SetStatusEffectStacks(state.ToString(1), int(state.ToNumber(2)))
		return 0
	}))

	l.SetGlobal("get_actor_status_effects", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetActorStatusEffects(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_status_effect", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetStatusEffect(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_status_effect_instance", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetStatusEffectInstance(state.ToString(1))))
		return 1
	}))

	// Cards

	l.SetGlobal("give_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GiveCard(state.ToString(1), state.ToString(2))))
		return 1
	}))

	l.SetGlobal("remove_card", l.NewFunction(func(state *lua.LState) int {
		session.RemoveCard(state.ToString(1))
		return 0
	}))

	l.SetGlobal("cast_card", l.NewFunction(func(state *lua.LState) int {
		session.CastCard(state.ToString(1), state.ToString(2))
		return 0
	}))

	l.SetGlobal("get_cards", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetCards(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_card", l.NewFunction(func(state *lua.LState) int {
		card, _ := session.GetCard(state.ToString(1))
		state.Push(luhelp.ToLua(state, card))
		return 1
	}))

	l.SetGlobal("get_card_instance", l.NewFunction(func(state *lua.LState) int {
		_, instance := session.GetCard(state.ToString(1))
		state.Push(luhelp.ToLua(state, instance))
		return 1
	}))

	// Damage & Heal

	l.SetGlobal("deal_damage", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 3 {
			state.Push(lua.LNumber(session.DealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), false)))
		} else {
			state.Push(lua.LNumber(session.DealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

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
			state.Push(luhelp.ToLua(state, session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), false)))
		} else {
			state.Push(luhelp.ToLua(state, session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	l.SetGlobal("heal", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 3 {
			state.Push(lua.LNumber(session.Heal(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), false)))
		} else {
			state.Push(lua.LNumber(session.Heal(state.ToString(1), state.ToString(2), int(state.ToNumber(3)), bool(state.ToBool(4)))))
		}
		return 1
	}))

	// Player

	l.SetGlobal("player_draw_card", l.NewFunction(func(state *lua.LState) int {
		session.PlayerDrawCard(int(state.ToNumber(1)))
		return 0
	}))

	l.SetGlobal("give_player_gold", l.NewFunction(func(state *lua.LState) int {
		session.GivePlayerGold(int(state.ToNumber(1)))
		return 0
	}))

	// Merchant

	l.SetGlobal("get_merchant", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(state, session.GetMerchant()))
		return 1
	}))

	l.SetGlobal("add_merchant_card", l.NewFunction(func(state *lua.LState) int {
		session.AddMerchantCard()
		return 0
	}))

	l.SetGlobal("add_merchant_artifact", l.NewFunction(func(state *lua.LState) int {
		session.AddMerchantArtifact()
		return 0
	}))

	l.SetGlobal("get_merchant_gold_max", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetMerchantGoldMax()))
		return 1
	}))

	// Random

	l.SetGlobal("gen_face", l.NewFunction(func(state *lua.LState) int {
		if state.GetTop() == 1 {
			state.Push(lua.LString(faces.Global.Gen(int(state.ToNumber(1)))))
		} else {
			state.Push(lua.LString(faces.Global.GenRand()))
		}
		return 1
	}))

	l.SetGlobal("random_card", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GetRandomCard(int(state.ToNumber(1)))))
		return 1
	}))

	l.SetGlobal("random_artifact", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GetRandomArtifact(int(state.ToNumber(1)))))
		return 1
	}))

	return l
}
