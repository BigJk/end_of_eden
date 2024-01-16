register_card("NANO_CHARGER", {
    name = l("cards.NANO_CHARGER.name", "Nano Charger"),
    description = string.format(
        l("cards.NANO_CHARGER.description","%s\n\nSupercharge your next attack. Deals %s damage."),
        highlight("One-Time"), highlight("Double")
    ),
    tags = { "BUFF", "_ACT_0" },
    max_level = 0,
    color = COLOR_RED,
    need_target = false,
    does_consume = true,
    point_cost = 0,
    price = 150,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("NANO_CHARGER", ctx.caster, 1 + ctx.level)
            return nil
        end
    },
    test = function()
        return assert_chain({
            function() return assert_cast_card("NANO_CHARGER") end,
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("NANO_CHARGER", 1) end,
        })
    end
})

register_status_effect("NANO_CHARGER", {
    name = l("status_effects.NANO_CHARGER.name", "Nano Charge"),
    description = string.format(l("status_effects.NANO_CHARGER.description", "Next attack deals %s damage."), highlight("Double")),
    look = "NC",
    foreground = COLOR_RED,
    can_stack = false,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.simulated then
                return ctx.damage * 2
            end
            
            if ctx.source == ctx.owner and ctx.target ~= ctx.owner then
                add_status_effect_stacks(ctx.guid, -1)
                return ctx.damage * 2
            end
            return ctx.damage
        end,
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("NANO_CHARGER", 1) end,
            function ()
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(PLAYER_ID, dummy, 5)
                if damage ~= 10 then
                    return "Expected 10 damage, got " .. damage
                end

                local hp_after = get_actor(dummy).hp
                if hp_after ~= 90 then
                    return "Expected 100 damage, got " .. hp_after
                end
            end
        })
    end
})
