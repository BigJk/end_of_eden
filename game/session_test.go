package game

import (
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"io"
	"log"
	"testing"
)

func TestSessionLua(t *testing.T) {
	session := NewSession(WithLogging(log.New(io.Discard, "", 0)))

	// Add one enemy
	enemyGuid := NewGuid()
	enemyActor := NewActor(enemyGuid)
	session.AddActor(enemyActor)

	//
	// Test get_player
	//
	t.Run("GetPlayer", func(t *testing.T) {
		// Set player name
		testName := "Test Value"
		testGold := 4123
		session.UpdatePlayer(func(actor *Actor) bool {
			actor.Name = testName
			actor.Gold = testGold
			return true
		})

		if err := session.luaState.DoString(`
player_name = get_player().name
player_gold = get_player().gold
`); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, testName, lua.LVAsString(session.luaState.GetGlobal("player_name")))
		assert.Equal(t, testGold, int(lua.LVAsNumber(session.luaState.GetGlobal("player_gold"))))
	})

	//
	// Test OnDamageCalc callback
	//
	t.Run("OnDamageCalc", func(t *testing.T) {
		if err := session.luaState.DoString(`
register_artifact(
    "DEBUG_DOUBLE_DAMAGE",
    {
        name = "Stone Of Gigantic Strength",
        description = "Double all damage dealt.",
        price = 1000,
        order = 100,
        callbacks = {
            on_damage_calc = function(ctx)
                return ctx.damage * 2
            end,
        }
    }
);

register_artifact(
    "DEBUG_MINUS",
    {
        name = "Minus",
        description = "",
        order = 0,
        callbacks = {
            on_damage_calc = function(ctx)
                return ctx.damage - 5
            end,
        }
    }
);
`); err != nil {
			t.Fatal(err)
		}

		session.UpdateActor(enemyGuid, func(actor *Actor) bool {
			actor.HP = 100
			actor.MaxHP = 100
			return true
		})

		session.GiveArtifact("DEBUG_DOUBLE_DAMAGE", PlayerActorID)
		session.GiveArtifact("DEBUG_MINUS", PlayerActorID)
		session.DealDamage(PlayerActorID, enemyGuid, 20, false)

		assert.Equal(t, 100-(20*2-5), session.GetActor(enemyGuid).HP)
	})

	//
	// Test OnCast callback
	//
	t.Run("OnCast", func(t *testing.T) {
		if err := session.luaState.DoString(`
register_card("DEBUG_MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 10 damage",
        color = "#cccccc",
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 10)
                return nil
            end,
        }
    }
);
`); err != nil {
			t.Fatal(err)
		}

		session.UpdateActor(enemyGuid, func(actor *Actor) bool {
			actor.HP = 100
			actor.MaxHP = 100
			return true
		})

		cardGuid := session.GiveCard("DEBUG_MELEE_HIT", PlayerActorID)
		session.CastCard(cardGuid, enemyGuid)

		assert.Equal(t, 100-(10*2-5), session.GetActor(enemyGuid).HP)
	})

	//
	// Test Enemy
	//

	t.Run("EnemyCast", func(t *testing.T) {
		enemyType := "DEBUG_ENEMY"

		if err := session.luaState.DoString(`
register_enemy(
    "DEBUG_ENEMY",
    {
        name = "Mutated Hamster",
        description = "Small but furious...",
        initial_hp = 20,
        callbacks = {
            on_init = function(ctx)
                give_card("DEBUG_MELEE_HIT", ctx.guid)
            end,
            on_turn = function(ctx)
                local cards = get_cards(ctx.guid)
                cast_card(cards[1], PLAYER_ID)
            end
        }
    }
)
`); err != nil {
			t.Fatal(err)
		}

		session.UpdatePlayer(func(actor *Actor) bool {
			actor.HP = 50
			actor.MaxHP = 50
			return true
		})

		// Remove old artifacts
		lo.ForEach(session.GetPlayer().Artifacts.ToSlice(), func(item string, index int) {
			session.RemoveArtifact(item)
		})

		enemyGuid := session.AddActorFromEnemy(enemyType)
		_, err := session.resources.Enemies[enemyType].Callbacks[CallbackOnInit](CreateContext("type_id", enemyType, "guid", enemyGuid))
		assert.NoError(t, err)

		assert.Equal(t, 50, session.GetPlayer().HP)

		_, err = session.resources.Enemies[enemyType].Callbacks[CallbackOnTurn](CreateContext("type_id", enemyType, "guid", enemyGuid))
		assert.NoError(t, err)

		assert.Equal(t, 50-10, session.GetPlayer().HP)

		_, err = session.resources.Enemies[enemyType].Callbacks[CallbackOnTurn](CreateContext("type_id", enemyType, "guid", enemyGuid))
		assert.NoError(t, err)

		assert.Equal(t, 50-20, session.GetPlayer().HP)
	})
}
