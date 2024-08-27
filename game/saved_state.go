package game

import (
	"encoding/gob"
	"math/rand/v2"
)

func init() {
	gob.Register(SavedState{})
}

// SavedState represents a save file that don't contain any pointer so the lua
// runtime or other pointer.
type SavedState struct {
	State            GameState
	Seed             uint64
	Rand             *rand.PCG
	Actors           map[string]Actor
	Instances        map[string]any
	StagesCleared    int
	CurrentEvent     string
	CurrentFight     FightState
	PointsPerRound   int
	Merchant         MerchantState
	EventHistory     []string
	StateCheckpoints []StateCheckpoint
	CtxData          map[string]any
	LoadedMods       []string
}
