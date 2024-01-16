register_card("ULTRA_FLASH_SHIELD", {
    name = l("cards.ULTRA_FLASH_SHIELD.name", "Ultra Flash Shield"),
    description = string.format(
        l("cards.ULTRA_FLASH_SHIELD.description","%s\n\nDeploy a temporary shield. %s all attack this turn."),
        highlight("One-Time"), highlight("Negates")
    ),
    tags = { "DEF", "_ACT_0" },
    max_level = 0,
    color = COLOR_BLUE,
    need_target = false,
    does_consume = true,
    point_cost = 3,
    price = 250,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("ULTRA_FLASH_SHIELD", ctx.caster, 1 + ctx.level)
            return nil
        end
    }
})

register_status_effect("ULTRA_FLASH_SHIELD", {
    name = l("status_effects.ULTRA_FLASH_SHIELD.name", "Ultra Flash Shield"),
    description = l("status_effects.ULTRA_FLASH_SHIELD.description", "Negates all attacks."),
    look = "UFS",
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
                return 0
            end
            return ctx.damage
        end,
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("ULTRA_FLASH_SHIELD", 1) end,
            function ()
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(dummy, PLAYER_ID, 100)
                if damage ~= 0 then
                    return "Expected 0 damage, got " .. damage
                end

                damage = deal_damage(dummy, PLAYER_ID, 2)
                if damage ~= 0 then
                    return "Expected 0 damage, got " .. damage
                end
            end
        })
    end
})
