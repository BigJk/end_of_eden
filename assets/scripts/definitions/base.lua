---@meta

---GUID of an actor, artifact, or status effect. References an active instance of the object.
---@alias guid string

---ID of an actor, artifact, or status effect. References the definition of the object.
---@alias type_id string

---Game state. Used to determine what the game should be doing at the moment. Can be one of: GAME_STATE_EVENT, GAME_STATE_FIGHT, GAME_STATE_MERCHANT, GAME_STATE_RANDOM
---@alias game_state string

---Next game state. Used to determine what the game should be doing next. Can be one of: GAME_STATE_EVENT, GAME_STATE_FIGHT, GAME_STATE_MERCHANT, GAME_STATE_RANDOM
---@alias next_game_state string

---Registered objects.
---@class registered
---@field card { [string]: card }
---@field artifact { [string]: artifact }
---@field event { [string]: event }
---@field story_teller { [string]: story_teller }
---@field status_effect { [string]: status_effect }
registered = {
    ["card"] = {},
    ["artifact"] = {},
    ["event"] = {},
    ["story_teller"] = {},
    ["status_effect"] = {},
}