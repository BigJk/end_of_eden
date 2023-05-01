package game

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"io"
	"log"
	"testing"
)

var TestArtifactLua = `

register_artifact(
    "TEST_ARTIFACT",
    {
        name = "Test Artifact",
        description = "This is a cool description",
        price = 1337,
        order = -10,
        callbacks = {
            on_pick_up = function(ctx)
                -- hello world
                return ctx
            end,
            on_player_turn = function(ctx)
                -- hello world
                return ctx
            end
        }
    }
);


`

func TestArtifact(t *testing.T) {
	s := lua.NewState()
	man := NewResourcesManager(s, nil, log.New(io.Discard, "", 0))

	// Evaluate lua
	if !assert.NoError(t, s.DoString(TestArtifactLua)) {
		return
	}

	art := man.Artifacts["TEST_ARTIFACT"]
	if !assert.NotNil(t, art) {
		return
	}

	// Check values
	assert.Equal(t, 1337, art.Price)
	assert.Equal(t, -10, art.Order)

	// Check if callback can be called without error
	res, err := art.Callbacks["OnPlayerTurn"]("1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, "1", res)
}
