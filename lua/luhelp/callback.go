package luhelp

// OwnedCallback represents a callback that will execute inside a lua vm.
type OwnedCallback func(args ...any) (any, error)

func (cb OwnedCallback) MarshalJSON() ([]byte, error) {
	return []byte("function"), nil
}

func (cb OwnedCallback) Call(args ...any) (any, error) {
	if cb == nil {
		return nil, nil
	}

	return cb(args...)
}
