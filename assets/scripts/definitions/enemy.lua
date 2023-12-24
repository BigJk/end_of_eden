---@meta

---@class enemy_intend_ctx
---@field type_id type_id
---@field guid guid
---@field round number

---Enemy represents a definition of a enemy that can be linked from a Actor.
---@class enemy
---@field id? type_id
---@field name string
---@field description string
---@field initial_hp number
---@field max_hp number
---@field look string
---@field color string
---@field intend? fun(ctx:enemy_intend_ctx):nil
---@field callbacks callbacks
---@field test? fun():nil|string
---@field base_game? boolean