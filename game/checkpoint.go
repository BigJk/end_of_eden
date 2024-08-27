package game

import (
	"encoding/gob"

	"github.com/samber/lo"
)

func init() {
	gob.Register(StateEventDeathData{})
	gob.Register(StateEventDamageData{})
	gob.Register(StateEventHealData{})
	gob.Register(StateEventMoneyData{})
	gob.Register(StateEventArtifactAddedData{})
	gob.Register(StateEventArtifactRemovedData{})
	gob.Register(StateEventCardAddedData{})
	gob.Register(StateEventCardRemovedData{})

	gob.Register(StateCheckpoint{})
	gob.Register(StateCheckpointMarker{})
}

type StateEvent string

const (
	StateEventDeath           = StateEvent("Death")
	StateEventDamage          = StateEvent("Damage")
	StateEventHeal            = StateEvent("Heal")
	StateEventMoney           = StateEvent("Money")
	StateEventArtifactAdded   = StateEvent("ArtifactAdded")
	StateEventArtifactRemoved = StateEvent("ArtifactRemoved")
	StateEventCardAdded       = StateEvent("CardAdded")
	StateEventCardRemoved     = StateEvent("CardRemoved")
)

type StateEventDeathData struct {
	Source string
	Target string
	Damage int
}

type StateEventDamageData struct {
	Source string
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

type StateEventArtifactAddedData struct {
	Owner  string
	GUID   string
	TypeID string
}

type StateEventArtifactRemovedData struct {
	Owner  string
	GUID   string
	TypeID string
}

type StateEventCardAddedData struct {
	Owner  string
	GUID   string
	TypeID string
}

type StateEventCardRemovedData struct {
	Owner  string
	GUID   string
	TypeID string
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
