package luhelp

import "encoding/json"

// OwnedCallback represents a callback that will execute inside a lua vm.
type OwnedCallback func(args ...any) (any, error)

func (cb OwnedCallback) MarshalJSON() ([]byte, error) {
	return json.Marshal("lua function")
}

// Call executes the callback with the given arguments. If the callback is nil
// it will return nil, nil. If the callback returns an error it will be returned
// as the error of this function.
func (cb OwnedCallback) Call(args ...any) (any, error) {
	if cb == nil {
		return nil, nil
	}

	return cb(args...)
}

// Present returns true if the callback is not nil.
func (cb OwnedCallback) Present() bool {
	return cb != nil
}
