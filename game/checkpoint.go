package game

import "github.com/samber/lo"

type StateEvent string

const (
	StateEventDeath  = StateEvent("Death")
	StateEventDamage = StateEvent("Damage")
	StateEventHeal   = StateEvent("Heal")
	StateEventMoney  = StateEvent("Money")
)

type StateEventDeathData struct {
	Target string
	Damage int
}

type StateEventDamageData struct {
	Target string
	Damage int
}

type StateEventHealData struct {
	Target string
	Damage int
}

type StateEventMoneyData struct {
	Target string
	Money  int
}

// StateCheckpoint saves the state of a session at a certain point. This can be used
// to retroactively check what happened between certain actions.
type StateCheckpoint struct {
	Session *Session

	// Events describe the events that
	Events map[StateEvent]any
}

// StateCheckpointMarker is a saved state of a checkpoint log.
type StateCheckpointMarker struct {
	checkpoints []StateCheckpoint
}

// Diff returns the new states that happened between the marker and a new session.
func (sm StateCheckpointMarker) Diff(session *Session) []StateCheckpoint {
	if len(sm.checkpoints) >= len(session.stateCheckpoints) {
		return nil
	}
	return session.stateCheckpoints[len(sm.checkpoints):]
}

// DiffEvent returns the new states that happened between the marker and a new session that contain a certain event.
func (sm StateCheckpointMarker) DiffEvent(session *Session, event StateEvent) []StateCheckpoint {
	return lo.Filter(sm.Diff(session), func(item StateCheckpoint, index int) bool {
		if item.Events == nil {
			return false
		}
		_, ok := item.Events[event]
		return ok
	})
}
