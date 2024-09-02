register_card("NULLIFY", {
    name = l("cards.NULLIFY.name", "Nullify"),
    description = string.format(
        l("cards.NULLIFY.description", "%s\n\nDeploy a temporary damage nullifier. %s all damage this round."),
        highlight("Exhaust"), highlight("Negates")
    ),
    tags = { "DEF", "_ACT_0" },
    max_level = 0,
    color = COLOR_BLUE,
    need_target = false,
    does_exhaust = true,
    point_cost = 3,
    price = 200,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("NULLIFY", ctx.caster, 1 + ctx.level)
            return nil
        end
    }
})

register_status_effect("NULLIFY", {
    name = l("status_effects.NULLIFY.name", "Nullify Field"),
    description = l("status_effects.NULLIFY.description", "Negates all damage this round."),
    look = "NF",
    foreground = COLOR_BLUE,
    can_stack = false,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                return 0
            end
            return ctx.damage
        end,
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("NULLIFY", 1) end,
            function()
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(dummy, PLAYER_ID, 100)
                if damage ~= 0 then
                    return "Expected 0 damage, got " .. damage
                end
            end
        })
    end
})
