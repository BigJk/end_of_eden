package game

import (
	"github.com/BigJk/end_of_eden/internal/lua/ludoc"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"io"
	"log"
	"testing"
)

var TestCardLua = `

register_card("MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 10 damage",
        color = "#cccccc",
        callbacks = {
            on_cast = function(ctx)
                return "hello_world"
            end,
        }
    }
);

`

func TestCards(t *testing.T) {
	s := lua.NewState()
	man := NewResourcesManager(s, ludoc.New(), log.New(io.Discard, "", 0))

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
	res, err := card.Callbacks[CallbackOnCast]("1", "2", "3", "4")
	assert.NoError(t, err)
	assert.Equal(t, "hello_world", res)
}
