---@meta

---@class card_state_ctx
---@field type_id type_id
---@field guid guid
---@field level number
---@field owner guid

---Card represents a playable card definition.
---@class card
---@field id? type_id
---@field name string
---@field description string
---@field tags? string[]
---@field state? fun(ctx:card_state_ctx):nil
---@field color string
---@field point_cost number
---@field max_level number
---@field does_exhaust? boolean
---@field need_target boolean
---@field price number
---@field callbacks callbacks
---@field test? fun():nil|string
---@field base_game? boolean

---CardInstance represents an instance of a card owned by some actor.
---@class card_instance
---@field guid guid
---@field type_id type_id
---@field level number
---@field owner guid