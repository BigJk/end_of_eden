register_card("FLASH_BANG", {
    name = l("cards.FLASH_BANG.name", "Flash Bang"),
    description = l("cards.FLASH_BANG.description", highlight("One-Time") .. "\n\nInflicts " .. highlight("Blinded") .. " on the target, causing them to deal less damage."),
    tags = { "CC" },
    max_level = 0,
    color = COLOR_PURPLE,
    need_target = true,
    does_consume = true,
    point_cost = 1,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("FLASH_BANG", ctx.target)
            return nil
        end
    }
})

register_status_effect("FLASH_BANG", {
    name = l("cards.FLASH_BANG.name", "Blinded"),
    description = l("cards.FLASH_BANG.description", "Causing " .. highlight("25%") .. " less damage."),
    look = "FL",
    foreground = COLOR_PURPLE,
    state = function(ctx) return nil end,
    can_stack = true,
    decay = DECAY_ONE,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage * 0.75
            end
            return ctx.damage
        end
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("FLASH_BANG", 1) end,
            function ()
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(PLAYER_ID, dummy, 10)
                if damage ~= 7 then
                    return "Expected 7 damage, got " .. damage
                end
            end
        })
    end
})
