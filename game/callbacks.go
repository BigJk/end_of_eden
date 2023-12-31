package game

import (
	"github.com/samber/lo"
	"log"
)

const (
	CallbackOnDamage        = "OnDamage"
	CallbackOnDamageCalc    = "OnDamageCalc"
	CallbackOnHealCalc      = "OnHealCalc"
	CallbackOnCast          = "OnCast"
	CallbackOnActorDidCast  = "OnActorDidCast"
	CallbackOnInit          = "OnInit"
	CallbackOnPickUp        = "OnPickUp"
	CallbackOnTurn          = "OnTurn"
	CallbackOnPlayerTurn    = "OnPlayerTurn"
	CallbackOnStatusAdd     = "OnStatusAdd"
	CallbackOnStatusStack   = "OnStatusStack"
	CallbackOnStatusRemove  = "OnStatusRemove"
	CallbackOnRemove        = "OnRemove"
	CallbackOnActorDie      = "OnActorDie"
	CallbackOnMerchantEnter = "OnMerchantEnter"
)

// Context represents the context arguments for a callback.
type Context map[string]any

var EmptyContext = Context{}

// CreateContext creates a new context with the given key value pairs. The number of arguments must be even.
// Example: CreateContext("key1", "value1", "key2", 124)
func CreateContext(args ...any) Context {
	if len(args)%2 != 0 {
		log.Printf("CreateContext: %v\n", args)
		panic("Please fix create context!")
	}

	val := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		val[args[i].(string)] = args[i+1]
	}
	return val
}

// Add adds the given key value pairs to the context. The number of arguments must be even.
func (c Context) Add(args ...any) Context {
	if len(args)%2 != 0 {
		log.Printf("CreateContext: %v\n", args)
		panic("Please fix create context!")
	}

	if len(args) == 0 {
		return c
	}

	kv := lo.FlatMap(lo.Entries(c), func(item lo.Entry[string, any], i int) []any {
		return []any{item.Key, item.Value}
	})

	return CreateContext(append(kv, args...)...)
}

// AddContext adds the given context to the context.
func (c Context) AddContext(ctx Context) Context {
	if len(ctx) == 0 {
		return c
	}
	return c.Add(lo.FlatMap(lo.Entries(ctx), func(item lo.Entry[string, any], i int) []any {
		return []any{item.Key, item.Value}
	})...)
}

// Trigger represents for which objects a event should be triggered.
type Trigger int

const (
	// TriggerArtifact triggers for artifacts.
	TriggerArtifact = Trigger(1)

	// TriggerStatusEffect triggers for status effects.
	TriggerStatusEffect = Trigger(2)

	// TriggerEnemy triggers for enemies.
	TriggerEnemy = Trigger(4)

	// TriggerAll triggers for all objects.
	TriggerAll = Trigger(TriggerArtifact | TriggerStatusEffect | TriggerEnemy)
)

// TriggerCallbackState represents the internal state of a trigger callback.
type TriggerCallbackState struct {
	Callback string
	Who      Trigger
	AddedCtx Context
	Ctx      []any
	Done     bool
}

func (t TriggerCallbackState) SetAddedCtx(ctx Context) TriggerCallbackState {
	t.AddedCtx = ctx
	return t
}

func (t TriggerCallbackState) SetCtx(ctx ...any) TriggerCallbackState {
	t.Ctx = ctx
	return t
}

func (t TriggerCallbackState) SetDone() TriggerCallbackState {
	t.Done = true
	return t
}

// TriggerCallback traverses all artifacts, status effects and enemies and calls the given callback for each object.
//   - callback: The name of the callback to call.
//   - who: For which objects the callback should be called.
//   - fn: The function to call for each object. The function gets the current value, the guid of the object, the type of the object and the current state as arguments. The function should return the new state.
//   - addedCtx: Additional context that should be added to the context.
//   - ctx: Additional context that should be added to the context.
func TriggerCallback(s *Session, callback string, who Trigger, fn func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState, addedCtx Context, ctx ...any) {
	state := TriggerCallbackState{
		Callback: callback,
		Who:      who,
		AddedCtx: addedCtx,
		Ctx:      ctx,
		Done:     false,
	}

	s.TraverseArtifactsStatus(lo.Keys(s.instances),
		func(instance ArtifactInstance, artifact *Artifact) {
			if state.Done || state.Who&TriggerArtifact == 0 {
				return
			}

			baseCtx := CreateContext("type_id", artifact.ID, "guid", instance.GUID, "owner", instance.Owner, "round", s.GetFightRound()).AddContext(state.AddedCtx)
			val, err := artifact.Callbacks[callback].Call(append([]any{baseCtx}, state.Ctx...)...)
			if err != nil {
				s.logLuaError(callback, instance.TypeID, err)
			}

			state = fn(val, instance.GUID, TriggerArtifact, state)
		},
		func(instance StatusEffectInstance, statusEffect *StatusEffect) {
			if state.Done || state.Who&TriggerStatusEffect == 0 {
				return
			}

			baseCtx := CreateContext("type_id", statusEffect.ID, "guid", instance.GUID, "owner", instance.Owner, "round", s.GetFightRound(), "stacks", instance.Stacks).AddContext(state.AddedCtx)
			val, err := statusEffect.Callbacks[callback].Call(append([]any{baseCtx}, state.Ctx...)...)
			if err != nil {
				s.logLuaError(callback, instance.TypeID, err)
			}

			state = fn(val, instance.GUID, TriggerStatusEffect, state)
		},
	)

	if state.Done || state.Who&TriggerEnemy == 0 {
		return
	}

	lo.ForEach(s.GetOpponents(PlayerActorID), func(actor Actor, index int) {
		if state.Done {
			return
		}

		if enemy := s.GetEnemy(actor.TypeID); enemy != nil {
			baseCtx := CreateContext("type_id", enemy.ID, "guid", actor.GUID, "round", s.GetFightRound()).AddContext(state.AddedCtx)
			val, err := enemy.Callbacks[callback].Call(append([]any{baseCtx}, state.Ctx...)...)
			if err != nil {
				s.logLuaError(callback, enemy.ID, err)
			}

			state = fn(val, actor.GUID, TriggerEnemy, state)
		}
	})

	return
}

// TriggerCallbackSimple is a helper function for TriggerCallback that returns all values of the callback as a slice.
func TriggerCallbackSimple(s *Session, callback string, who Trigger, addedCtx Context, ctx ...any) []any {
	var returns []any

	TriggerCallback(s, callback, who, func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState {
		returns = append(returns, val)
		return state
	}, addedCtx, ctx...)

	return returns
}

// TriggerCallbackReduce is a helper function for TriggerCallback that reduces all values of the callback to a single value.
func TriggerCallbackReduce[T any](s *Session, callback string, who Trigger, reducer func(cur T, val T) T, initial T, propagatedCtxKey string, addedCtx Context, ctx ...any) T {
	cur := initial

	TriggerCallback(s, callback, who, func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState {
		if val, ok := val.(T); ok {
			cur = reducer(cur, val)
		}

		if propagatedCtxKey != "" {
			return state.SetAddedCtx(state.AddedCtx.Add(propagatedCtxKey, cur))
		}

		return state
	}, addedCtx, ctx...)

	return cur
}

// TriggerCallbackFirst is a helper function for TriggerCallback that returns the first value of the callback.
func TriggerCallbackFirst[T any](s *Session, callback string, who Trigger, addedCtx Context, ctx ...any) T {
	var returns T

	TriggerCallback(s, callback, who, func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState {
		if val, ok := val.(T); ok {
			returns = val
		}

		return state
	}, addedCtx, ctx...)

	return returns
}

// TriggerCallbackFirstStop is a helper function for TriggerCallback that returns the first value of the callback and stops the callback.
func TriggerCallbackFirstStop[T any](s *Session, callback string, who Trigger, addedCtx Context, ctx ...any) T {
	var returns T

	TriggerCallback(s, callback, who, func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState {
		if val, ok := val.(T); ok {
			returns = val
			return state.SetDone()
		}

		return state
	}, addedCtx, ctx...)

	return returns
}

// TriggerCallbackFirstOfStop is a helper function for TriggerCallback that returns the first value of the callback that matches the given check function.
// After the first match the callback is stopped.
func TriggerCallbackFirstOfStop[T any](s *Session, callback string, who Trigger, check func(val T) bool, addedCtx Context, ctx ...any) T {
	var returns T

	TriggerCallback(s, callback, who, func(val any, guid string, t Trigger, state TriggerCallbackState) TriggerCallbackState {
		if val, ok := val.(T); ok {
			if check(val) {
				returns = val
				return state.SetDone()
			}
		}

		return state
	}, addedCtx, ctx...)

	return returns
}
