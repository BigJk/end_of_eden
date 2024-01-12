---@meta

---@alias decay_type string

---@class status_effect_state_ctx
---@field type_id type_id
---@field guid guid
---@field stacks number
---@field owner guid

--- Status effect defintion
---@class status_effect
---@field id? type_id
---@field name string
---@field description string
---@field state? fun(ctx:status_effect_state_ctx):nil
---@field look string
---@field foreground string
---@field order? number
---@field can_stack boolean
---@field decay decay_type
---@field rounds number
---@field callbacks callbacks
---@field test? fun():nil|string
---@field base_game? boolean

--- Status effect instance
---@class status_effect_instance
---@field guid guid
---@field type_id type_id
---@field owner guid
---@field stacks number
---@field rounds_left number
---@field round_entered number
