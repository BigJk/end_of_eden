# End Of Eden Lua Docs
## Index

- [Game Constants](#game-constants)
- [Utility](#utility)
- [Styling](#styling)
- [Logging](#logging)
- [Audio](#audio)
- [Game State](#game-state)
- [Actor Operations](#actor-operations)
- [Artifact Operations](#artifact-operations)
- [Status Effect Operations](#status-effect-operations)
- [Card Operations](#card-operations)
- [Damage & Heal](#damage--heal)
- [Player Operations](#player-operations)
- [Merchant Operations](#merchant-operations)
- [Random Utility](#random-utility)
- [Content Registry](#content-registry)

## Game Constants

General game constants.

### Globals
<details> <summary><b><code>DECAY_ALL</code></b> </summary> <br/>

Status effect decays by all stacks per turn.

</details>

<details> <summary><b><code>DECAY_NONE</code></b> </summary> <br/>

Status effect never decays.

</details>

<details> <summary><b><code>DECAY_ONE</code></b> </summary> <br/>

Status effect decays by 1 stack per turn.

</details>

<details> <summary><b><code>GAME_STATE_EVENT</code></b> </summary> <br/>

Represents the event game state.

</details>

<details> <summary><b><code>GAME_STATE_FIGHT</code></b> </summary> <br/>

Represents the fight game state.

</details>

<details> <summary><b><code>GAME_STATE_MERCHANT</code></b> </summary> <br/>

Represents the merchant game state.

</details>

<details> <summary><b><code>GAME_STATE_RANDOM</code></b> </summary> <br/>

Represents the random game state in which the active story teller will decide what happens next.

</details>

<details> <summary><b><code>PLAYER_ID</code></b> </summary> <br/>

Player actor id for use in functions where the guid is needed, for example: ``deal_damage(PLAYER_ID, enemy_id, 10)``.

</details>

### Functions

None

## Utility

General game constants.

### Globals

None

### Functions
<details> <summary><b><code>fetch</code></b> </summary> <br/>

Fetches a value from the persistent store

**Signature:**

```
fetch(key : String) -> Any
```

</details>

<details> <summary><b><code>guid</code></b> </summary> <br/>

returns a new random guid.

**Signature:**

```
guid() -> String
```

</details>

<details> <summary><b><code>store</code></b> </summary> <br/>

Stores a persistent value for this run that will be restored after a save load. Can store any lua basic value or table.

**Signature:**

```
store(key : String) -> None
```

</details>

## Styling

Helper functions for text styling.

### Globals

None

### Functions
<details> <summary><b><code>text_bg</code></b> </summary> <br/>

Makes the text background colored. Takes hex values like #ff0000.

**Signature:**

```
text_bg(color : String, value) -> String
```

</details>

<details> <summary><b><code>text_bold</code></b> </summary> <br/>

Makes the text bold.

**Signature:**

```
text_bold(value) -> String
```

</details>

<details> <summary><b><code>text_color</code></b> </summary> <br/>

Makes the text foreground colored. Takes hex values like #ff0000.

**Signature:**

```
text_color(color : String, value) -> String
```

</details>

<details> <summary><b><code>text_italic</code></b> </summary> <br/>

Makes the text italic.

**Signature:**

```
text_italic(value) -> String
```

</details>

<details> <summary><b><code>text_underline</code></b> </summary> <br/>

Makes the text underlined.

**Signature:**

```
text_underline(value) -> String
```

</details>

## Logging

Various logging functions.

### Globals

None

### Functions
<details> <summary><b><code>log_d</code></b> </summary> <br/>

Log at **danger** level to player log.

**Signature:**

```
log_d(value) -> None
```

</details>

<details> <summary><b><code>log_i</code></b> </summary> <br/>

Log at **information** level to player log.

**Signature:**

```
log_i(value) -> None
```

</details>

<details> <summary><b><code>log_s</code></b> </summary> <br/>

Log at **success** level to player log.

**Signature:**

```
log_s(value) -> None
```

</details>

<details> <summary><b><code>log_w</code></b> </summary> <br/>

Log at **warning** level to player log.

**Signature:**

```
log_w(value) -> None
```

</details>

<details> <summary><b><code>print</code></b> </summary> <br/>

Log to session log.

**Signature:**

```
print(value, value, value...) -> None
```

</details>

## Audio

Audio helper functions.

### Globals

None

### Functions
<details> <summary><b><code>play_audio</code></b> </summary> <br/>

Plays a sound effect. If you want to play ``button.mp3`` you call ``play_audio("button")``.

**Signature:**

```
play_audio(sound : String) -> None
```

</details>

<details> <summary><b><code>play_music</code></b> </summary> <br/>

Start a song for the background loop. If you want to play ``song.mp3`` you call ``play_music("song")``.

**Signature:**

```
play_music(sound : String) -> None
```

</details>

## Game State

Functions that modify the general game state.

### Globals

None

### Functions
<details> <summary><b><code>get_event_history</code></b> </summary> <br/>

Gets the ids of all the encountered events in the order of occurrence.

**Signature:**

```
get_event_history() -> Array
```

</details>

<details> <summary><b><code>get_fight</code></b> </summary> <br/>

Gets the fight state. This contains the player hand, used, exhausted and round information.

**Signature:**

```
get_fight() -> Table
```

</details>

<details> <summary><b><code>get_fight_round</code></b> </summary> <br/>

Gets the number of stages cleared.

**Signature:**

```
get_fight_round() -> Number
```

</details>

<details> <summary><b><code>had_event</code></b> </summary> <br/>

Checks if the event happened at least once.

**Signature:**

```
had_event(eventId : String) -> Bool
```

</details>

<details> <summary><b><code>had_events</code></b> </summary> <br/>

Checks if all the events happened at least once.

**Signature:**

```
had_events(eventIds : Array) -> Bool
```

</details>

<details> <summary><b><code>had_events_any</code></b> </summary> <br/>

Checks if any of the events happened at least once.

**Signature:**

```
had_events_any(eventIds : Array) -> Bool
```

</details>

<details> <summary><b><code>set_event</code></b> </summary> <br/>

Set event by id.

**Signature:**

```
set_event(eventId : String) -> None
```

</details>

<details> <summary><b><code>set_fight_description</code></b> </summary> <br/>

Set the current fight description. This will be shown on the top right in the game.

**Signature:**

```
set_fight_description(desc : String) -> None
```

</details>

<details> <summary><b><code>set_game_state</code></b> </summary> <br/>

Set the current game state. See globals.

**Signature:**

```
set_game_state(state : String) -> None
```

</details>

## Actor Operations

Functions that modify or access the actors. Actors are either the player or enemies.

### Globals

None

### Functions
<details> <summary><b><code>actor_add_max_hp</code></b> </summary> <br/>

Increases the max hp value of a actor by a number. Can be negative value to decrease it.

**Signature:**

```
actor_add_max_hp(guid : String, amount : Number) -> None
```

</details>

<details> <summary><b><code>add_actor_by_enemy</code></b> </summary> <br/>

Creates a new enemy fighting against the player. Example ``add_actor_by_enemy("RUST_MITE")``.

**Signature:**

```
add_actor_by_enemy(enemyId : String) -> None
```

</details>

<details> <summary><b><code>get_actor</code></b> </summary> <br/>

Get a actor by guid.

**Signature:**

```
get_actor(guid : String) -> Table
```

</details>

<details> <summary><b><code>get_opponent_by_index</code></b> </summary> <br/>

Get opponent (actor) by index of a certain actor. ``get_opponent_by_index(PLAYER_ID, 2)`` would return the second alive opponent of the player.

**Signature:**

```
get_opponent_by_index(guid : String, index : Number) -> Table
```

</details>

<details> <summary><b><code>get_opponent_count</code></b> </summary> <br/>

Get the number of opponents (actors) of a certain actor. ``get_opponent_count(PLAYER_ID)`` would return 2 if the player had 2 alive enemies.

**Signature:**

```
get_opponent_count(guid : String) -> Table
```

</details>

<details> <summary><b><code>get_opponent_guids</code></b> </summary> <br/>

Get the guids of opponents (actors) of a certain actor. If the player had 2 enemies, ``get_opponent_guids(PLAYER_ID)`` would return a table with 2 strings containing the guids of these actors.

**Signature:**

```
get_opponent_guids(guid : String) -> Table
```

</details>

<details> <summary><b><code>get_player</code></b> </summary> <br/>

Get the player actor. Equivalent to ``get_actor(PLAYER_ID)``

**Signature:**

```
get_player() -> Table
```

</details>

<details> <summary><b><code>remove_actor</code></b> </summary> <br/>

Deletes a actor by id.

**Signature:**

```
remove_actor(guid : String) -> None
```

</details>

## Artifact Operations

Functions that modify or access the artifacts.

### Globals

None

### Functions
<details> <summary><b><code>get_artifact</code></b> </summary> <br/>

Returns the artifact definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.

**Signature:**

```
get_artifact(id : String) -> Table
```

</details>

<details> <summary><b><code>get_artifact_instance</code></b> </summary> <br/>

Returns the artifact instance by guid.

**Signature:**

```
get_artifact_instance(guid : String) -> None
```

</details>

<details> <summary><b><code>get_random_artifact_type</code></b> </summary> <br/>

Returns a random type id given a max gold price.

**Signature:**

```
get_random_artifact_type(maxGold : Number) -> None
```

</details>

<details> <summary><b><code>give_artifact</code></b> </summary> <br/>

Gives a actor a artifact. Returns the guid of the newly created artifact.

**Signature:**

```
give_artifact(typeId : String, actor : String) -> String
```

</details>

<details> <summary><b><code>remove_artifact</code></b> </summary> <br/>

Removes a artifact.

**Signature:**

```
remove_artifact(guid : String) -> None
```

</details>

## Status Effect Operations

Functions that modify or access the status effects.

### Globals

None

### Functions
<details> <summary><b><code>add_status_effect_stacks</code></b> </summary> <br/>

Adds to the stack count of a status effect. Negative values are also allowed.

**Signature:**

```
add_status_effect_stacks(guid : String, count : Number) -> None
```

</details>

<details> <summary><b><code>get_actor_status_effects</code></b> </summary> <br/>

Returns the guids of all status effects that belong to a actor.

**Signature:**

```
get_actor_status_effects(actorId : String) -> Array
```

</details>

<details> <summary><b><code>get_status_effect</code></b> </summary> <br/>

Returns the status effect definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.

**Signature:**

```
get_status_effect(id : String) -> Table
```

</details>

<details> <summary><b><code>get_status_effect_instance</code></b> </summary> <br/>

Returns the status effect instance.

**Signature:**

```
get_status_effect_instance(effectGuid : String) -> Table
```

</details>

<details> <summary><b><code>give_status_effect</code></b> </summary> <br/>

Gives a status effect to a actor. If count is not specified a stack of 1 is applied.

**Signature:**

```
give_status_effect(typeId : String, actorGuid : String, (optional) count : Number) -> None
```

</details>

<details> <summary><b><code>remove_status_effect</code></b> </summary> <br/>

Removes a status effect.

**Signature:**

```
remove_status_effect(guid : String) -> None
```

</details>

<details> <summary><b><code>set_status_effect_stacks</code></b> </summary> <br/>

Sets the stack count of a status effect by guid.

**Signature:**

```
set_status_effect_stacks(guid : String, count : Number) -> None
```

</details>

## Card Operations

Functions that modify or access the cards.

### Globals

None

### Functions
<details> <summary><b><code>cast_card</code></b> </summary> <br/>

Tries to cast a card with a guid and optional target. If the cast isn't successful returns false.

**Signature:**

```
cast_card(cardGuid : String, (optional) targetActorGuid : String) -> Bool
```

</details>

<details> <summary><b><code>get_card</code></b> </summary> <br/>

Returns the card type definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.

**Signature:**

```
get_card(id : String) -> Table
```

</details>

<details> <summary><b><code>get_card_instance</code></b> </summary> <br/>

Returns the instance object of a card.

**Signature:**

```
get_card_instance(cardGuid : String) -> Table
```

</details>

<details> <summary><b><code>get_cards</code></b> </summary> <br/>

Returns all the card guids from the given actor.

**Signature:**

```
get_cards(actorGuid : String) -> Array
```

</details>

<details> <summary><b><code>give_card</code></b> </summary> <br/>

Gives a card.

**Signature:**

```
give_card(cardTypeId : String, ownerActorId : String) -> String
```

</details>

<details> <summary><b><code>remove_card</code></b> </summary> <br/>

Removes a card.

**Signature:**

```
remove_card(cardGuid : String) -> None
```

</details>

<details> <summary><b><code>upgrade_card</code></b> </summary> <br/>

Upgrade a card without paying for it.

**Signature:**

```
upgrade_card(cardGuid : String) -> Bool
```

</details>

## Damage & Heal

Functions that deal damage or heal.

### Globals

None

### Functions
<details> <summary><b><code>deal_damage</code></b> </summary> <br/>

Deal damage to a enemy from one source. If flat is true the damage can't be modified by status effects or artifacts.

**Signature:**

```
deal_damage(source : String, target : String, damage : Number, flat : Bool) -> None
```

</details>

<details> <summary><b><code>deal_damage_multi</code></b> </summary> <br/>

Deal damage to multiple enemies from one source. If flat is true the damage can't be modified by status effects or artifacts.

**Signature:**

```
deal_damage_multi(source : String, targets : Array, damage : Number, flat : Bool) -> None
```

</details>

<details> <summary><b><code>heal</code></b> </summary> <br/>

Heals the target triggered by the source.

**Signature:**

```
heal(source : String, target : String, amount : Number) -> None
```

</details>

## Player Operations

Functions that are related to the player.

### Globals

None

### Functions
<details> <summary><b><code>give_player_gold</code></b> </summary> <br/>

Gives the player gold.

**Signature:**

```
give_player_gold(amount : Number) -> None
```

</details>

<details> <summary><b><code>player_buy_artifact</code></b> </summary> <br/>

Let the player buy the artifact with the given id. This will deduct the price form the players gold and return true if the buy was successful.

**Signature:**

```
player_buy_artifact(artifactId : String) -> Bool
```

</details>

<details> <summary><b><code>player_buy_card</code></b> </summary> <br/>

Let the player buy the card with the given id. This will deduct the price form the players gold and return true if the buy was successful.

**Signature:**

```
player_buy_card(cardId : String) -> Bool
```

</details>

<details> <summary><b><code>player_draw_card</code></b> </summary> <br/>

Let the player draw additional cards for this turn.

**Signature:**

```
player_draw_card(amount : Number) -> None
```

</details>

<details> <summary><b><code>player_give_action_points</code></b> </summary> <br/>

Gives the player more action points for this turn.

**Signature:**

```
player_give_action_points(points : Number) -> None
```

</details>

## Merchant Operations

Functions that are related to the merchant.

### Globals

None

### Functions
<details> <summary><b><code>add_merchant_artifact</code></b> </summary> <br/>

Adds another random artifact to the merchant

**Signature:**

```
add_merchant_artifact() -> None
```

</details>

<details> <summary><b><code>add_merchant_card</code></b> </summary> <br/>

Adds another random card to the merchant

**Signature:**

```
add_merchant_card() -> None
```

</details>

<details> <summary><b><code>get_merchant</code></b> </summary> <br/>

Returns the merchant state.

**Signature:**

```
get_merchant() -> Table
```

</details>

<details> <summary><b><code>get_merchant_gold_max</code></b> </summary> <br/>

Returns the maximum value of artifacts and cards that the merchant will sell. Good to scale ``random_card`` and ``random_artifact``.

**Signature:**

```
get_merchant_gold_max() -> Number
```

</details>

## Random Utility

Functions that help with random generation.

### Globals

None

### Functions
<details> <summary><b><code>gen_face</code></b> </summary> <br/>

Generates a random face.

**Signature:**

```
gen_face((optional) category : Number) -> String
```

</details>

<details> <summary><b><code>random_artifact</code></b> </summary> <br/>

Returns the type id of a random artifact.

**Signature:**

```
random_artifact(maxPrice : Number) -> String
```

</details>

<details> <summary><b><code>random_card</code></b> </summary> <br/>

Returns the type id of a random card.

**Signature:**

```
random_card(maxPrice : Number) -> String
```

</details>

## Content Registry

These functions are used to define new content in the base game and in mods.

### Globals

None

### Functions
<details> <summary><b><code>register_artifact</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)
```

**Signature:**

```
register_artifact(id : String, definition : Table) -> None
```

</details>

<details> <summary><b><code>register_card</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_card("MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 5 (+3 for each upgrade) damage.",
        state = function(ctx)
            return "Use your bare hands to deal " .. highlight(5 + ctx.level * 3) .. " damage."
        end,
        max_level = 1,
        color = "#2f3e46",
        need_target = true,
        point_cost = 1,
        price = 30,
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 5 + ctx.level * 3)
                return nil
            end,
        }
    }
)
```

**Signature:**

```
register_card(id : String, definition : Table) -> None
```

</details>

<details> <summary><b><code>register_enemy</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_enemy("RUST_MITE",
    {
        name = "Rust Mite",
        description = "Loves to eat metal.",
        look = "/v\\",
        color = "#e6e65a",
        initial_hp = 22,
        max_hp = 22,
        gold = 10,
        callbacks = {
            on_turn = function(ctx)
                if ctx.round % 4 == 0 then
                    give_status_effect("RITUAL", ctx.guid)
                else
                    deal_damage(ctx.guid, PLAYER_ID, 6)
                end

                return nil
            end
        }
    }
)
```

**Signature:**

```
register_enemy(id : String, definition : Table) -> None
```

</details>

<details> <summary><b><code>register_event</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_event("SOME_EVENT",
	{
		name = "Event Name",
		description = [[Flavor Text... Can include **Markdown** Syntax!]],
		choices = {
			{
				description = "Go...",
				callback = function()
					-- If you return nil on_end will decide the next game state
					return nil 
				end
			},
			{
				description = "Other Option",
				callback = function() return GAME_STATE_FIGHT end
			}
		},
		on_enter = function()
			play_music("energetic_orthogonal_expansions")
	
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("MELEE_HIT", PLAYER_ID)
			give_card("RUPTURE", PLAYER_ID)
			give_card("BLOCK", PLAYER_ID)
			give_artifact(get_random_artifact_type(150), PLAYER_ID)
		end,
		on_end = function(choice)
			-- Choice will be nil or the index of the choice taken
			return GAME_STATE_RANDOM
		end,
	}
)
```

**Signature:**

```
register_event(id : String, definition : Table) -> None
```

</details>

<details> <summary><b><code>register_status_effect</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)
```

**Signature:**

```
register_status_effect(id : String, definition : Table) -> None
```

</details>

<details> <summary><b><code>register_story_teller</code></b> </summary> <br/>

Registers a new artifact.

```lua
register_artifact("REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
)
```

**Signature:**

```
register_story_teller(id : String, definition : Table) -> None
```

</details>

