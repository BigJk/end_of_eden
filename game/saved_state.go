package game

import "encoding/gob"

func init() {
	gob.Register(SavedState{})
}

// SavedState represents a save file that don't contain any pointer so the lua
// runtime or other pointer.
type SavedState struct {
	State            GameState
	Actors           map[string]Actor
	Instances        map[string]any
	StagesCleared    int
	CurrentEvent     string
	CurrentFight     FightState
	Merchant         MerchantState
	EventHistory     []string
	StateCheckpoints []StateCheckpoint
	CtxData          map[string]any
}
