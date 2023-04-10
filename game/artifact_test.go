package game

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

var TestArtifactLua = `

register_artifact(
    "TEST_ARTIFACT",
    {
        Name = "Test Artifact",
        Description = "This is a cool description",
        Price = 1337,
        Order = -10,
        Callbacks = {
            OnPickUp = function(type, guid, owner)
                -- hello world
                return type
            end,
            OnCombatEnd = function(type, guid, owner)
                -- hello world
                return type
            end
        }
    }
);


`

func TestArtifact(t *testing.T) {
	s := lua.NewState()
	man := NewResourcesManager(s)

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
	res, err := art.Callbacks["OnCombatEnd"]("1", "2", "3")
	assert.NoError(t, err)
	assert.Equal(t, "1", res)
}
