package game

// EventChoice represents a possible choice in the Event.
type EventChoice struct {
	Description string
	Callback    OwnedCallback
}

// Event represents a encounter-able event.
type Event struct {
	ID          string
	Name        string
	Description string
	Choices     []EventChoice
	OnEnter     OwnedCallback
	OnEnd       OwnedCallback
}

func (e Event) IsNone() bool {
	return len(e.ID) == 0
}
