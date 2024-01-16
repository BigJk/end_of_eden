---@meta

-- #####################################
-- Game Constants
-- #####################################

--- Status effect decays by all stacks per turn.
DECAY_ALL = ""

--- Status effect never decays.
DECAY_NONE = ""

--- Status effect decays by 1 stack per turn.
DECAY_ONE = ""

--- Represents the event game state.
GAME_STATE_EVENT = ""

--- Represents the fight game state.
GAME_STATE_FIGHT = ""

--- Represents the merchant game state.
GAME_STATE_MERCHANT = ""

--- Represents the random game state in which the active story teller will decide what happens next.
GAME_STATE_RANDOM = ""

--- Player actor id for use in functions where the guid is needed, for example: ``deal_damage(PLAYER_ID, enemy_guid, 10)``.
PLAYER_ID = ""

-- #####################################
-- Utility
-- #####################################

--- Fetches a value from the persistent store
---@param key string
---@return any
function fetch(key) end

--- returns a new random guid.
---@return guid
function guid() end

--- Stores a persistent value for this run that will be restored after a save load. Can store any lua basic value or table.
---@param key string
---@param value any
function store(key, value) end

-- #####################################
-- Styling
-- #####################################

--- Makes the text background colored. Takes hex values like #ff0000.
---@param color string
---@param value any
---@return string
function text_bg(color, value) end

--- Makes the text bold.
---@param value any
---@return string
function text_bold(value) end

--- Makes the text italic.
---@param value any
---@return string
function text_italic(value) end

--- Makes the text underlined.
---@param value any
---@return string
function text_underline(value) end

-- #####################################
-- Logging
-- #####################################

--- Log at **danger** level to player log.
---@param value any
function log_d(value) end

--- Log at **information** level to player log.
---@param value any
function log_i(value) end

--- Log at **success** level to player log.
---@param value any
function log_s(value) end

--- Log at **warning** level to player log.
---@param value any
function log_w(value) end

--- Log to session log.
---@param ... any
function print(...) end

-- #####################################
-- Audio
-- #####################################

--- Plays a sound effect. If you want to play ``button.mp3`` you call ``play_audio("button")``.
---@param sound string
function play_audio(sound) end

--- Start a song for the background loop. If you want to play ``song.mp3`` you call ``play_music("song")``.
---@param sound string
function play_music(sound) end

-- #####################################
-- Game State
-- #####################################

--- Gets the ids of all the encountered events in the order of occurrence.
---@return string[]
function get_event_history() end

--- Gets the fight state. This contains the player hand, used, exhausted and round information.
---@return fight_state
function get_fight() end

--- Gets the fight round.
---@return number
function get_fight_round() end

--- Gets the number of stages cleared.
---@return number
function get_stages_cleared() end

--- Checks if the event happened at least once.
---@param event_id type_id
---@return boolean
function had_event(event_id) end

--- Checks if all the events happened at least once.
---@param event_ids type_id[]
---@return boolean
function had_events(event_ids) end

--- Checks if any of the events happened at least once.
---@param eventIds string[]
---@return boolean
function had_events_any(eventIds) end

--- Set event by id.
---@param event_id type_id
function set_event(event_id) end

--- Set the current fight description. This will be shown on the top right in the game.
---@param desc string
function set_fight_description(desc) end

--- Set the current game state. See globals.
---@param state next_game_state
function set_game_state(state) end

-- #####################################
-- Actor Operations
-- #####################################

--- Increases the hp value of a actor by a number. Can be negative value to decrease it. This won't trigger any on_damage callbacks
---@param guid guid
---@param amount number
function actor_add_hp(guid, amount) end

--- Increases the max hp value of a actor by a number. Can be negative value to decrease it.
---@param guid guid
---@param amount number
function actor_add_max_hp(guid, amount) end

--- Sets the hp value of a actor to a number. This won't trigger any on_damage callbacks
---@param guid guid
---@param amount number
function actor_set_hp(guid, amount) end

--- Sets the max hp value of a actor to a number.
---@param guid guid
---@param amount number
function actor_set_max_hp(guid, amount) end

--- Creates a new enemy fighting against the player. Example ``add_actor_by_enemy("RUST_MITE")``.
---@param enemy_guid type_id
---@return string
function add_actor_by_enemy(enemy_guid) end

--- Get a actor by guid.
---@param guid guid
---@return actor
function get_actor(guid) end

--- Get opponent (actor) by index of a certain actor. ``get_opponent_by_index(PLAYER_ID, 2)`` would return the second alive opponent of the player.
---@param guid guid
---@param index number
---@return actor
function get_opponent_by_index(guid, index) end

--- Get the number of opponents (actors) of a certain actor. ``get_opponent_count(PLAYER_ID)`` would return 2 if the player had 2 alive enemies.
---@param guid guid
---@return number
function get_opponent_count(guid) end

--- Get the guids of opponents (actors) of a certain actor. If the player had 2 enemies, ``get_opponent_guids(PLAYER_ID)`` would return a table with 2 strings containing the guids of these actors.
---@param guid guid
---@return guid[]
function get_opponent_guids(guid) end

--- Get the player actor. Equivalent to ``get_actor(PLAYER_ID)``
---@return actor
function get_player() end

--- Deletes a actor by id.
---@param guid guid
function remove_actor(guid) end

-- #####################################
-- Artifact Operations
-- #####################################

--- Returns the artifact definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.
---@param id string
---@return artifact
function get_artifact(id) end

--- Returns the artifact instance by guid.
---@param guid guid
---@return artifact_instance
function get_artifact_instance(guid) end

--- Returns all the artifacts guids from the given actor.
---@param actor_guid string
---@return guid[]
function get_artifacts(actor_guid) end

--- Gives a actor a artifact. Returns the guid of the newly created artifact.
---@param type_id type_id
---@param actor guid
---@return string
function give_artifact(type_id, actor) end

--- Removes a artifact.
---@param guid guid
function remove_artifact(guid) end

-- #####################################
-- Status Effect Operations
-- #####################################

--- Adds to the stack count of a status effect. Negative values are also allowed.
---@param guid guid
---@param count number
function add_status_effect_stacks(guid, count) end

--- Returns the guids of all status effects that belong to a actor.
---@param actor_guid string
---@return guid[]
function get_actor_status_effects(actor_guid) end

--- Returns the status effect definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.
---@param id string
---@return status_effect
function get_status_effect(id) end

--- Returns the status effect instance.
---@param effect_guid guid
---@return status_effect_instance
function get_status_effect_instance(effect_guid) end

--- Gives a status effect to a actor. If count is not specified a stack of 1 is applied.
---@param type_id string
---@param actor_guid string
---@param count? number
function give_status_effect(type_id, actor_guid, count) end

--- Removes a status effect.
---@param guid guid
function remove_status_effect(guid) end

--- Sets the stack count of a status effect by guid.
---@param guid guid
---@param count number
function set_status_effect_stacks(guid, count) end

-- #####################################
-- Card Operations
-- #####################################

--- Tries to cast a card with a guid and optional target. If the cast isn't successful returns false.
---@param card_guid guid
---@param target_actor_guid? guid
---@return boolean
function cast_card(card_guid, target_actor_guid) end

--- Returns the card type definition. Can take either a guid or a typeId. If it's a guid it will fetch the type behind the instance.
---@param id type_id
---@return card
function get_card(id) end

--- Returns the instance object of a card.
---@param card_guid guid
---@return card_instance
function get_card_instance(card_guid) end

--- Returns all the card guids from the given actor.
---@param actor_guid string
---@return guid[]
function get_cards(actor_guid) end

--- Gives a card.
---@param card_type_id type_id
---@param owner_actor_guid guid
---@return string
function give_card(card_type_id, owner_actor_guid) end

--- Removes a card.
---@param card_guid string
function remove_card(card_guid) end

--- Upgrade a card without paying for it.
---@param card_guid guid
---@return boolean
function upgrade_card(card_guid) end

--- Upgrade a random card without paying for it.
---@param actor_guid guid
---@return boolean
function upgrade_random_card(actor_guid) end

-- #####################################
-- Damage & Heal
-- #####################################

--- Deal damage from one source to a target. If flat is true the damage can't be modified by status effects or artifacts. Returns the damage that was dealt.
---@param source guid
---@param target guid
---@param damage number
---@param flat? boolean
---@return number
function deal_damage(source, target, damage, flat) end

--- Deal damage from one source to a target from a card. If flat is true the damage can't be modified by status effects or artifacts. Returns the damage that was dealt.
---@param source guid
---@param card guid
---@param target guid
---@param damage number
---@param flat? boolean
---@return number
function deal_damage_card(source, card, target, damage, flat) end

--- Deal damage to multiple enemies from one source. If flat is true the damage can't be modified by status effects or artifacts. Returns a array of damages for each actor hit.
---@param source guid
---@param targets guid[]
---@param damage number
---@param flat? boolean
---@return number[]
function deal_damage_multi(source, targets, damage, flat) end

--- Heals the target triggered by the source.
---@param source guid
---@param target guid
---@param amount number
function heal(source, target, amount) end

--- Simulate damage from a source to a target. If flat is true the damage can't be modified by status effects or artifacts. Returns the damage that would be dealt.
---@param source guid
---@param target guid
---@param damage number
---@param flat? boolean
---@return number
function simulate_deal_damage(source, target, damage, flat) end

-- #####################################
-- Player Operations
-- #####################################

--- Finishes the player turn.
function finish_player_turn() end

--- Gives the player gold.
---@param amount number
function give_player_gold(amount) end

--- Let the player buy the artifact with the given id. This will deduct the price form the players gold and return true if the buy was successful.
---@param card_id type_id
---@return boolean
function player_buy_artifact(card_id) end

--- Let the player buy the card with the given id. This will deduct the price form the players gold and return true if the buy was successful.
---@param card_id type_id
---@return boolean
function player_buy_card(card_id) end

--- Let the player draw additional cards for this turn.
---@param amount number
function player_draw_card(amount) end

--- Gives the player more action points for this turn.
---@param points number
function player_give_action_points(points) end

-- #####################################
-- Merchant Operations
-- #####################################

--- Adds another random artifact to the merchant
function add_merchant_artifact() end

--- Adds another random card to the merchant
function add_merchant_card() end

--- Returns the merchant state.
---@return merchant_state
function get_merchant() end

--- Returns the maximum value of artifacts and cards that the merchant will sell. Good to scale ``random_card`` and ``random_artifact``.
---@return number
function get_merchant_gold_max() end

-- #####################################
-- Random Utility
-- #####################################

--- Generates a random face.
---@param category? number
---@return string
function gen_face(category) end

--- Returns the type id of a random artifact.
---@param max_price number
---@return type_id
function random_artifact(max_price) end

--- Returns the type id of a random card.
---@param max_price number
---@return type_id
function random_card(max_price) end

-- #####################################
-- Localization
-- #####################################

--- Returns the localized string for the given key. Examples on locals definition can be found in `/assets/locals`. Example: ``
--- l('cards.MY_CARD.name', "English Default Name")``
---@param key string
---@param default? string
---@return string
function l(key, default) end

-- #####################################
-- Content Registry
-- #####################################

--- Deletes all base game content. Useful if you don't want to include base game content in your mod.
--- 
--- ```lua
--- delete_base_game() -- delete all base game content
--- delete_base_game("artifact") -- deletes all artifacts
--- delete_base_game("card") -- deletes all cards
--- delete_base_game("enemy") -- deletes all enemies
--- delete_base_game("event") -- deletes all events
--- delete_base_game("status_effect") -- deletes all status effects
--- delete_base_game("story_teller") -- deletes all story tellers
--- 
--- ```
---@param type? string
function delete_base_game(type) end

--- Deletes a card.
--- 
--- ```lua
--- delete_card("SOME_CARD")
--- ```
---@param id type_id
function delete_card(id) end

--- Deletes an enemy.
--- 
--- ```lua
--- delete_enemy("SOME_ENEMY")
--- ```
---@param id type_id
function delete_enemy(id) end

--- Deletes an event.
--- 
--- ```lua
--- delete_event("SOME_EVENT")
--- ```
---@param id type_id
function delete_event(id) end

--- Deletes a status effect.
--- 
--- ```lua
--- delete_status_effect("SOME_STATUS_EFFECT")
--- ```
---@param id type_id
function delete_status_effect(id) end

--- Deletes a story teller.
--- 
--- ```lua
--- delete_story_teller("SOME_STORY_TELLER")
--- ```
---@param id type_id
function delete_story_teller(id) end

--- Registers a new artifact.
--- 
--- ```lua
--- register_artifact("REPULSION_STONE",
---     {
---         name = "Repulsion Stone",
---         description = "For each damage taken heal for 2",
---         price = 100,
---         order = 0,
---         callbacks = {
---             on_damage = function(ctx)
---                 if ctx.target == ctx.owner then
---                     heal(ctx.owner, 2)
---                 end
---                 return nil
---             end,
---         }
---     }
--- )
--- ```
---@param id type_id
---@param definition artifact
function register_artifact(id, definition) end

--- Registers a new card.
--- 
--- ```lua
--- register_card("MELEE_HIT",
---     {
---         name = "Melee Hit",
---         description = "Use your bare hands to deal 5 (+3 for each upgrade) damage.",
---         state = function(ctx)
---             return "Use your bare hands to deal " .. highlight(5 + ctx.level * 3) .. " damage."
---         end,
---         max_level = 1,
---         color = "#2f3e46",
---         need_target = true,
---         point_cost = 1,
---         price = 30,
---         callbacks = {
---             on_cast = function(ctx)
---                 deal_damage(ctx.caster, ctx.target, 5 + ctx.level * 3)
---                 return nil
---             end,
---         }
---     }
--- )
--- ```
---@param id type_id
---@param definition card
function register_card(id, definition) end

--- Registers a new enemy.
--- 
--- ```lua
--- register_enemy("RUST_MITE",
---     {
---         name = "Rust Mite",
---         description = "Loves to eat metal.",
---         look = "/v\\",
---         color = "#e6e65a",
---         initial_hp = 22,
---         max_hp = 22,
---         gold = 10,
---         callbacks = {
---             on_turn = function(ctx)
---                 if ctx.round % 4 == 0 then
---                     give_status_effect("RITUAL", ctx.guid)
---                 else
---                     deal_damage(ctx.guid, PLAYER_ID, 6)
---                 end
--- 
---                 return nil
---             end
---         }
---     }
--- )
--- ```
---@param id type_id
---@param definition enemy
function register_enemy(id, definition) end

--- Registers a new event.
--- 
--- ```lua
--- register_event("SOME_EVENT",
--- 	{
--- 		name = "Event Name",
--- 		description = "Flavor Text... Can include **Markdown** Syntax!",
--- 		choices = {
--- 			{
--- 				description = "Go...",
--- 				callback = function()
--- 					-- If you return nil on_end will decide the next game state
--- 					return nil 
--- 				end
--- 			},
--- 			{
--- 				description = "Other Option",
--- 				callback = function() return GAME_STATE_FIGHT end
--- 			}
--- 		},
--- 		on_enter = function()
--- 			play_music("energetic_orthogonal_expansions")
--- 	
--- 			give_card("MELEE_HIT", PLAYER_ID)
--- 			give_card("MELEE_HIT", PLAYER_ID)
--- 			give_card("MELEE_HIT", PLAYER_ID)
--- 			give_card("RUPTURE", PLAYER_ID)
--- 			give_card("BLOCK", PLAYER_ID)
--- 			give_artifact(get_random_artifact_type(150), PLAYER_ID)
--- 		end,
--- 		on_end = function(choice)
--- 			-- Choice will be nil or the index of the choice taken
--- 			return GAME_STATE_RANDOM
--- 		end,
--- 	}
--- )
--- ```
---@param id type_id
---@param definition event
function register_event(id, definition) end

--- Registers a new status effect.
--- 
--- ```lua
--- register_status_effect("BLOCK", {
---     name = "Block",
---     description = "Decreases incoming damage for each stack",
---     look = "Blk",
---     foreground = "#219ebc",
---     state = function(ctx)
---         return "Takes " .. highlight(ctx.stacks) .. " less damage"
---     end,
---     can_stack = true,
---     decay = DECAY_ALL,
---     rounds = 1,
---     order = 100,
---     callbacks = {
---         on_damage_calc = function(ctx)
---             if ctx.target == ctx.owner then
---                 add_status_effect_stacks(ctx.guid, -ctx.damage)
---                 return ctx.damage - ctx.stacks
---             end
---             return ctx.damage
---         end
---     }
--- })
--- ```
---@param id type_id
---@param definition status_effect
function register_status_effect(id, definition) end

--- Registers a new story teller.
--- 
--- ```lua
--- register_story_teller("STORY_TELLER_XYZ", {
---     active = function(ctx)
---         if not had_events_any({ "A", "B", "C" }) then
---             return 1
---         end
---         return 0
---     end,
---     decide = function(ctx)
---         local stage = get_stages_cleared()
--- 
---         if stage >= 3 then
---             set_event("SOME_EVENT")
---             return GAME_STATE_EVENT
---         end
--- 
---         -- Fight against rust mites or clean bots
---         local d = math.random(2)
---         if d == 1 then
---             add_actor_by_enemy("RUST_MITE")
---         elseif d == 2 then
---             add_actor_by_enemy("CLEAN_BOT")
---         end
--- 
---         return GAME_STATE_FIGHT
---     end
--- })
--- ```
---@param id type_id
---@param definition story_teller
function register_story_teller(id, definition) end

