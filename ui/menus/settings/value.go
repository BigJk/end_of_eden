package settings

type Type int

const (
	Bool   Type = iota
	Int    Type = iota
	String Type = iota
	Float  Type = iota
)

type Value struct {
	Key         string
	Name        string
	Description string
	Type        Type
	Val         any
	Min         any
	Max         any
}

type Values []Value

type Saver func([]Value) error
