package game

import mapset "github.com/deckarep/golang-set/v2"

const PlayerActorID = "PLAYER"

// Actor represents a player or enemy.
type Actor struct {
	ID            string
	TypeID        string
	Name          string
	Description   string
	HP            int
	MaxHP         int
	Gold          int
	Artifacts     mapset.Set[string]
	Cards         mapset.Set[string]
	StatusEffects mapset.Set[string]
}

func (a Actor) IsNone() bool {
	return len(a.ID) == 0
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
		ID:            ID,
		Artifacts:     mapset.NewThreadUnsafeSet[string](),
		Cards:         mapset.NewThreadUnsafeSet[string](),
		StatusEffects: mapset.NewThreadUnsafeSet[string](),
	}
}
