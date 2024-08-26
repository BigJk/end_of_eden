register_card("SMOKE_BOMB", {
    name = "Smoke Bomb",
    description = "Reduces the accuracy of all enemies for 1 turn.",
    tags = { "CC", "_ACT_0" },
    max_level = 0,
    color = COLOR_GRAY,
    need_target = false,
    does_consume = true,
    point_cost = 0,
    price = 150,
    callbacks = {
        on_cast = function(ctx)
            local enemies = get_opponent_guids(PLAYER_ID)
            for _, enemy in ipairs(enemies) do
                give_status_effect("SMOKE_BOMB", enemy, 1)
            end
            return nil
        end
    }
})

register_status_effect("SMOKE_BOMB", {
    name = "Smoke Bomb",
    description = "Reduces accuracy by 50%.",
    look = "SB",
    foreground = COLOR_GRAY,
    can_stack = false,
    decay = DECAY_ONE,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage * 0.5
            end
            return ctx.damage
        end
    }
})
