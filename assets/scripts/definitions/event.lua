---@meta

---@class event_on_enter_ctx
---@field type_id type_id

---@class event_choice_ctx
---@field type_id type_id
---@field choice number

---EventChoice represents a possible choice in the Event.
---@class event_choice
---@field description? string
---@field description_fn? fun():nil|string
---@field callback fun(ctx:event_choice_ctx):next_game_state|nil

---Event represents a encounter-able event.
---@class event
---@field id? string
---@field name string
---@field description string
---@field choices event_choice[]
---@field on_enter? fun(ctx:event_on_enter_ctx):nil
---@field on_end fun(ctx:event_choice_ctx):next_game_state|nil
---@field test? fun():nil|string
---@field base_game? boolean