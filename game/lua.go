package game

import (
	"fmt"
	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/luhelp"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	"io/fs"
	"log"
	"path/filepath"
	"strings"
)

// SessionAdapter creates a lua vm that is bound to the session in the given Session.
func SessionAdapter(session *Session) *lua.LState {
	l := lua.NewState()
	mapper := luhelp.NewMapper(l)

	_ = filepath.Walk("./assets/scripts/libs", func(path string, info fs.FileInfo, err error) error {
		if info != nil && info.IsDir() || !strings.HasSuffix(path, ".lua") {
			return nil
		}

		name := strings.Split(filepath.Base(path), ".")[0]

		mod, err := l.LoadFile(path)
		if err != nil {
			log.Println("Can't LoadFile module:", path)
			return nil
		}

		log.Println("Loaded lib:", path, name)

		preload := l.GetField(l.GetField(l.Get(lua.EnvironIndex), "package"), "preload")
		l.SetField(preload, name, mod)

		return nil
	})

	// Constants

	l.SetGlobal("PLAYER_ID", lua.LString(PlayerActorID))
	l.SetGlobal("GAME_STATE_FIGHT", lua.LString(GameStateFight))
	l.SetGlobal("GAME_STATE_EVENT", lua.LString(GameStateEvent))
	l.SetGlobal("GAME_STATE_MERCHANT", lua.LString(GameStateMerchant))
	l.SetGlobal("GAME_STATE_RANDOM", lua.LString(GameStateRandom))

	// Misc

	l.SetGlobal("log_i", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeInfo, state.ToString(1))
		return 0
	}))

	l.SetGlobal("log_w", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeWarning, state.ToString(1))
		return 0
	}))

	l.SetGlobal("log_d", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeDanger, state.ToString(1))
		return 0
	}))

	l.SetGlobal("log_s", l.NewFunction(func(state *lua.LState) int {
		session.Log(LogTypeSuccess, state.ToString(1))
		return 0
	}))

	l.SetGlobal("debug_log", l.NewFunction(func(state *lua.LState) int {
		dbg, ok := state.GetStack(1)
		if ok {
			_, _ = state.GetInfo("nSl", dbg, lua.LNil)
		}

		log.Printf("[LUA :: %d %s] %s \n", dbg.CurrentLine, dbg.Source, strings.Join(lo.Map(make([]any, state.GetTop()), func(_ any, index int) string {
			val := state.Get(1 + index)

			switch val.Type() {
			case lua.LTString:
				return lua.LVAsString(val)
			case lua.LTNumber:
				return fmt.Sprint(float64(lua.LVAsNumber(val)))
			case lua.LTBool:
				return fmt.Sprint(lua.LVAsBool(val))
			case lua.LTTable:
				var data map[string]interface{}
				if err := mapper.Map(val.(*lua.LTable), &data); err != nil {
					return "Error: " + err.Error()
				}
				return fmt.Sprint(data)
			case lua.LTUserData:
				return fmt.Sprint(val.(*lua.LUserData).Value)
			case lua.LTNil:
				return "nil"
			}

			return "<" + val.Type().String() + ">"
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
		state.Push(lua.LNumber(session.GetFightRound()))
		return 1
	}))

	l.SetGlobal("get_event_history", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(session.GetEventHistory()))
		return 1
	}))

	l.SetGlobal("had_event", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(session.HadEvent(state.ToString(1))))
		return 1
	}))

	// Actor Operations

	l.SetGlobal("get_player", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(session.GetPlayer()))
		return 1
	}))

	l.SetGlobal("get_actor", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(session.GetActor(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_opponent_by_index", l.NewFunction(func(state *lua.LState) int {
		log.Println(int(state.ToNumber(2)) - 1)
		state.Push(luhelp.ToLua(session.GetOpponentByIndex(state.ToString(1), int(state.ToNumber(2))-1)))
		return 1
	}))

	l.SetGlobal("get_opponent_count", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetOpponentCount(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_opponent_guids", l.NewFunction(func(state *lua.LState) int {
		state.Push(luhelp.ToLua(session.GetOpponentGUIDs(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("remove_actor", l.NewFunction(func(state *lua.LState) int {
		session.GetActor(state.ToString(1))
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

	// Artifacts

	l.SetGlobal("give_status_effect", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LString(session.GiveStatusEffect(state.ToString(1), state.ToString(2))))
		return 1
	}))

	l.SetGlobal("remove_status_effect", l.NewFunction(func(state *lua.LState) int {
		session.RemoveStatusEffect(state.ToString(1))
		return 0
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
		state.Push(luhelp.ToLua(session.GetCards(state.ToString(1))))
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
				log.Printf("Error in deal_damage_multi: %v\n", err)
				return 0
			}
		case lua.LTUserData:
			if val, ok := state.Get(2).(*lua.LUserData).Value.([]string); ok {
				guids = val
			}
		default:
			log.Printf("Error in deal_damage_multi: wrong type %v", state.Get(2).Type().String())
			return 0
		}

		if state.GetTop() == 3 {
			state.Push(luhelp.ToLua(session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), false)))
		} else {
			state.Push(luhelp.ToLua(session.DealDamageMulti(state.ToString(1), guids, int(state.ToNumber(3)), bool(state.ToBool(4)))))
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

	return l
}
