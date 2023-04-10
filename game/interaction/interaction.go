package interaction

type Event = string

// Provider makes it possible to stop the game to request interaction from a consumer like the UI.
type Provider interface {
	Choice([]string) chan int
	YesNo(string) chan bool
	Notify(string) chan bool
	On(Event, any) chan struct{}
}

// EmptyProvider response is always 0, true and struct{}.
type EmptyProvider struct{}

func (e *EmptyProvider) Choice(strings []string) chan int {
	return Instant(0)
}

func (e *EmptyProvider) YesNo(s string) chan bool {
	return Instant(true)
}

func (e *EmptyProvider) Notify(s string) chan bool {
	return Instant(true)
}

func (e *EmptyProvider) On(event Event, a any) chan struct{} {
	return InstantEmpty()
}

// Instant creates a channel containing the given value.
func Instant[T any](val T) chan T {
	c := make(chan T, 1)
	c <- val
	return c
}

// InstantEmpty creates a channel containing a empty struct.
func InstantEmpty() chan struct{} {
	c := make(chan struct{}, 1)
	c <- struct{}{}
	return c
}
