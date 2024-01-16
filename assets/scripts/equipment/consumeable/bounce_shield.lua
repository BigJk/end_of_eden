register_card("BOUNCE_SHIELD", {
    name = l("cards.BOUNCE_SHIELD.name", "Bounce Shield"),
    description = string.format(
        l("cards.BOUNCE_SHIELD.description","%s\n\nDeploy a temporary shield. %s bounces the damage back, but still takes damage."),
        highlight("One-Time"), highlight("Negates")
    ),
    tags = { "DEF", "_ACT_0" },
    max_level = 0,
    color = COLOR_BLUE,
    need_target = false,
    does_consume = true,
    point_cost = 0,
    price = 150,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("BOUNCE_SHIELD", ctx.caster, 1 + ctx.level)
            return nil
        end
    }
})

register_status_effect("BOUNCE_SHIELD", {
    name = l("status_effects.BOUNCE_SHIELD.name", "Bounce Shield"),
    description = l("status_effects.BOUNCE_SHIELD.description", "Bounces back the next damage. Still takes damage."),
    look = "BS",
    foreground = COLOR_BLUE,
    can_stack = false,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.simulated then
                return ctx.damage
            end
            
            if ctx.target == ctx.owner then
                add_status_effect_stacks(ctx.guid, -1)
                deal_damage(ctx.target, ctx.source, ctx.damage)
            end
            return ctx.damage
        end,
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("BOUNCE_SHIELD", 1) end,
            function ()
                local hp_before = get_actor(PLAYER_ID).hp
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(dummy, PLAYER_ID, 100)
                if damage ~= 100 then
                    return "Expected 100 damage, got " .. damage
                end

                local hp_after = get_actor(PLAYER_ID).hp
                if hp_before - hp_after ~= 100 then
                    return "Expected 100 damage, got " .. (hp_before - hp_after)
                end
            end
        })
    end
})
