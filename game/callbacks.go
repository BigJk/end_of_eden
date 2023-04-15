package game

import "log"

const (
	CallbackOnDamage       = "OnDamage"
	CallbackOnDamageCalc   = "OnDamageCalc"
	CallbackOnHealCalc     = "OnHealCalc"
	CallbackOnCast         = "OnCast"
	CallbackOnInit         = "OnInit"
	CallbackOnPickUp       = "OnPickUp"
	CallbackOnTurn         = "OnTurn"
	CallbackOnStatusAdd    = "OnStatusAdd"
	CallbackOnStatusStack  = "OnStatusStack"
	CallbackOnStatusRemove = "OnStatusRemove"
	CallbackOnRemove       = "OnRemove"
)

type Context map[string]any

func CreateContext(args ...any) Context {
	if len(args)%2 != 0 {
		log.Printf("CreateContext: %v\n", args)
		panic("Please fix create context!")
	}

	val := map[string]any{}
	for i := 0; i < len(args); i += 2 {
		val[args[i].(string)] = args[i+1]
	}
	return val
}
