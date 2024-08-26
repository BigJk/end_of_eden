register_card("ADRENALINE_SHOT", {
    name = "Adrenaline Shot",
    description = "Gain 2 additional action points for the next 3 turns.",
    tags = { "BUFF", "_ACT_0" },
    max_level = 0,
    color = COLOR_RED,
    need_target = false,
    does_consume = true,
    point_cost = 0,
    price = 400,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("ADRENALINE_SHOT", ctx.caster, 3)
            return nil
        end
    }
})

register_status_effect("ADRENALINE_SHOT", {
    name = "Adrenaline Shot",
    description = "Gain 2 additional action points.",
    look = "AS",
    foreground = COLOR_RED,
    can_stack = false,
    decay = DECAY_ONE,
    rounds = 3,
    order = 100,
    callbacks = {
        on_player_turn = function(ctx)
            if ctx.owner == PLAYER_ID then
                player_give_action_points(2)
            end
            return nil
        end
    }
})
