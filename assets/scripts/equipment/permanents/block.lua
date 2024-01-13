register_card("BLOCK", {
    name = "Block",
    description = "Shield yourself and gain 5 " .. highlight("block") .. ".",
    state = function(ctx)
        return "Shield yourself and gain " .. highlight(1 + ctx.level) .. " block."
    end,
    tags = { "DEF" },
    max_level = 1,
    color = COLOR_BLUE,
    need_target = false,
    point_cost = 1,
    price = 50,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("BLOCK", ctx.caster, 1 + ctx.level)
            return nil
        end
    }
})

register_status_effect("BLOCK", {
    name = "Block",
    description = "Decreases incoming damage for each stack",
    look = "B",
    foreground = COLOR_BLUE,
    state = function(ctx)
        return "Takes " .. highlight(ctx.stacks) .. " less damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.simulated then
                return ctx.damage
            end
            
            if ctx.target == ctx.owner then
                add_status_effect_stacks(ctx.guid, -ctx.damage)
                return ctx.damage - ctx.stacks
            end
            return ctx.damage
        end,
    },
    test = function()
        return assert_chain({
            function() return assert_status_effect_count(1) end,
            function() return assert_status_effect("BLOCK", 1) end,
            function ()
                local dummy = add_actor_by_enemy("DUMMY")
                local damage = deal_damage(dummy, PLAYER_ID, 1)
                if damage ~= 0 then
                    return "Expected 0 damage, got " .. damage
                end

                damage = deal_damage(dummy, PLAYER_ID, 2)
                if damage ~= 2 then
                    return "Expected 2 damage, got " .. damage
                end
            end
        })
    end
})
