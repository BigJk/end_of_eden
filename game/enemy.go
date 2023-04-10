package game

// Enemy represents a definition of a enemy that can be linked from a Actor.
type Enemy struct {
	ID          string
	Name        string
	Description string
	InitialHP   int
	MaxHP       int
	Callbacks   map[string]OwnedCallback
}
