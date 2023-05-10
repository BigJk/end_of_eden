package game

import "encoding/gob"

func init() {
	gob.Register(Actor{})
}

const PlayerActorID = "PLAYER"

// Actor represents a player or enemy.
type Actor struct {
	GUID          string `lua:"guid"`
	TypeID        string
	Name          string
	Description   string
	HP            int
	MaxHP         int
	Gold          int
	Artifacts     *StringSet
	Cards         *StringSet
	StatusEffects *StringSet
}

func (a Actor) IsNone() bool {
	return len(a.GUID) == 0
}

// Sanitize ensures that the actor has all the required fields.
func (a Actor) Sanitize() Actor {
	if a.Artifacts == nil {
		a.Artifacts = NewStringSet()
	}
	if a.Cards == nil {
		a.Cards = NewStringSet()
	}
	if a.StatusEffects == nil {
		a.StatusEffects = NewStringSet()
	}
	return a
}

func (a Actor) Clone() Actor {
	// The sets are backed by maps, so we need to clone them to create new pointer instances.
	a.Artifacts = a.Artifacts.Clone()
	a.Cards = a.Cards.Clone()
	a.StatusEffects = a.StatusEffects.Clone()
	return a
}

func NewActor(ID string) Actor {
	return Actor{
		GUID:          ID,
		Artifacts:     NewStringSet(),
		Cards:         NewStringSet(),
		StatusEffects: NewStringSet(),
	}
}
