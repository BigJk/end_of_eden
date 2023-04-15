package game

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestSessionLua(t *testing.T) {
	session := NewSession()

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
player_name = get_player().Name
player_gold = get_player().Gold
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
    "DOUBLE_DAMAGE",
    {
        Name = "Stone Of Gigantic Strength",
        Description = "Double all damage dealt.",
        Price = 1000,
        Order = -10,
        Callbacks = {
            OnDamageCalc = function(type, guid, source, target, damage)
                return damage * 2
            end,
        }
    }
);

register_artifact(
    "MINUS",
    {
        Name = "Minus",
        Description = "",
        Order = 0,
        Callbacks = {
            OnDamageCalc = function(type, guid, source, target, damage)
                return damage - 5
            end,
        }
    }
);
`); err != nil {
			t.Fatal(err)
		}

		enemyActor.HP = 100
		enemyActor.MaxHP = 100

		session.GiveArtifact("DOUBLE_DAMAGE", PlayerActorID)
		session.GiveArtifact("MINUS", PlayerActorID)

		session.DealDamage(PlayerActorID, enemyGuid, 20, false)

		assert.Equal(t, 100-(20*2-5), enemyActor.HP)
	})

	//
	// Test OnCast callback
	//
	t.Run("OnCast", func(t *testing.T) {
		if err := session.luaState.DoString(`
register_card("MELEE_HIT",
    {
        Name = "Melee Hit",
        Description = "Use your bare hands to deal 10 damage",
        Color = "#cccccc",
        Callbacks = {
            OnCast = function(type, guid, caster, target)
                deal_damage(caster, target, 10)
                return nil
            end,
        }
    }
);
`); err != nil {
			t.Fatal(err)
		}

		enemyActor.HP = 100
		enemyActor.MaxHP = 100

		cardGuid := session.GiveCard("MELEE_HIT", PlayerActorID)
		session.CastCard(cardGuid, enemyGuid)

		assert.Equal(t, 100-(10*2-5), enemyActor.HP)
	})

	//
	// Test Enemy
	//

	t.Run("EnemyCast", func(t *testing.T) {
		enemyType := "MUTATED_HAMSTER"

		if err := session.luaState.DoString(`
register_enemy(
    "MUTATED_HAMSTER",
    {
        Name = "Mutated Hamster",
        Description = "Small but furious...",
        InitialHP = 20,
        Callbacks = {
            OnInit = function(type, guid)
                give_card("MELEE_HIT", guid)
            end,
            OnTurn = function(type, guid)
                local cards = get_cards(guid)
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

		enemyGuid := session.AddActorFromEnemy(enemyType)
		_, err := session.resources.Enemies[enemyType].Callbacks["OnInit"](enemyType, enemyGuid)
		assert.NoError(t, err)

		_, err = session.resources.Enemies[enemyType].Callbacks["OnTurn"](enemyType, enemyGuid)
		assert.NoError(t, err)

		assert.Equal(t, 50-10, session.GetPlayer().HP)

		_, err = session.resources.Enemies[enemyType].Callbacks["OnTurn"](enemyType, enemyGuid)
		assert.NoError(t, err)

		assert.Equal(t, 50-20, session.GetPlayer().HP)
	})
}
