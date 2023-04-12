package game

import (
	"github.com/BigJk/project_gonzo/audio"
	"github.com/BigJk/project_gonzo/gluamapper"
	"github.com/samber/lo"
	lua "github.com/yuin/gopher-lua"
	luar "layeh.com/gopher-luar"
)

// OwnedCallback represents a callback that will execute inside a lua vm.
type OwnedCallback func(args ...any) (any, error)

func (cb OwnedCallback) Call(args ...any) (any, error) {
	if cb == nil {
		return nil, nil
	}

	return cb(args...)
}

// NewMapper creates a new lua -> go mapper that is able to convert lua functions to OwnedCallback.
func NewMapper(state *lua.LState) *gluamapper.Mapper {
	return gluamapper.NewMapper(gluamapper.Option{
		TagName: "lua",
		FnHook: func(value lua.LValue) any {
			return BindToLua(state, value)
		},
	})
}

// BindToLua will create a OwnedCallback from a lua function and state.
func BindToLua(state *lua.LState, value lua.LValue) OwnedCallback {
	return func(args ...any) (any, error) {
		// Call our lua function
		if err := state.CallByParam(lua.P{
			Fn:      value,
			NRet:    1,
			Protect: true,
		}, lo.Map(args, func(item any, index int) lua.LValue {
			return luar.New(state, item)
		})...); err != nil {
			return nil, err
		}

		// Fetch return value
		ret := state.Get(-1)
		state.Pop(1)

		// Parse to accepted return values
		switch ret.Type() {
		case lua.LTString:
			return lua.LVAsString(ret), nil
		case lua.LTNumber:
			return float64(lua.LVAsNumber(ret)), nil
		case lua.LTBool:
			return lua.LVAsBool(ret), nil
		case lua.LTTable:
			mapper := NewMapper(state)
			var data map[string]interface{}
			if err := mapper.Map(ret.(*lua.LTable), &data); err != nil {
				return nil, err
			}
			return data, nil
		}

		// Don't error for now
		return nil, nil
	}
}

// SessionAdapter creates a lua vm that is bound to the session in the given Session.
func SessionAdapter(session *Session) *lua.LState {
	l := lua.NewState()

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
		session.DebugLog(state.ToString(1))
		return 0
	}))

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

	l.SetGlobal("set_fight_description", l.NewFunction(func(state *lua.LState) int {
		session.SetFightDescription(state.ToString(1))
		return 0
	}))

	// Actor Operations

	l.SetGlobal("get_player", l.NewFunction(func(state *lua.LState) int {
		state.Push(luar.New(state, session.GetPlayer()))
		return 1
	}))

	l.SetGlobal("get_actor", l.NewFunction(func(state *lua.LState) int {
		state.Push(luar.New(state, session.GetActor(state.ToString(1))))
		return 1
	}))

	l.SetGlobal("get_opponent_by_index", l.NewFunction(func(state *lua.LState) int {
		state.Push(luar.New(state, session.GetOpponentByIndex(state.ToString(1), int(state.ToNumber(2)))))
		return 1
	}))

	l.SetGlobal("remove_actor", l.NewFunction(func(state *lua.LState) int {
		session.GetActor(state.ToString(1))
		return 0
	}))

	l.SetGlobal("get_opponent_count", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.GetOpponentCount(state.ToString(1))))
		return 1
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
		state.Push(luar.New(state, session.GetCards(state.ToString(1))))
		return 1
	}))

	// Damage & Heal

	l.SetGlobal("deal_damage", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.DealDamage(state.ToString(1), state.ToString(2), int(state.ToNumber(3)))))
		return 1
	}))

	l.SetGlobal("heal", l.NewFunction(func(state *lua.LState) int {
		state.Push(lua.LNumber(session.Heal(state.ToString(1), state.ToString(2), int(state.ToNumber(3)))))
		return 1
	}))

	return l
}
