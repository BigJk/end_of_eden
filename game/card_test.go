package game

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

var TestCardLua = `

register_card("MELEE_HIT",
    {
        Name = "Melee Hit",
        Description = "Use your bare hands to deal 10 damage",
        Color = "#cccccc",
        Callbacks = {
            OnCast = function(type, guid, caster, target)
                return "hello_world"
            end,
        }
    }
);

`

func TestCards(t *testing.T) {
	s := lua.NewState()
	man := NewResourcesManager(s)

	// Evaluate lua
	if !assert.NoError(t, s.DoString(TestCardLua)) {
		return
	}

	card := man.Cards["MELEE_HIT"]
	if !assert.NotNil(t, card) {
		return
	}

	// Check values
	assert.Equal(t, "#cccccc", card.Color)
	assert.Equal(t, "Melee Hit", card.Name)

	// Check if callback can be called without error
	res, err := card.Callbacks["OnCast"]("1", "2", "3", "4")
	assert.NoError(t, err)
	assert.Equal(t, "hello_world", res)
}
