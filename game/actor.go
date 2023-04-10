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

func NewActor(ID string) *Actor {
	return &Actor{
		ID:            ID,
		Artifacts:     mapset.NewSet[string](),
		Cards:         mapset.NewSet[string](),
		StatusEffects: mapset.NewSet[string](),
	}
}
