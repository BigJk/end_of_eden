register_status_effect("WEAKEN", {
    name = "Weaken",
    description = "Weakens damage for each stack",
    look = "W",
    foreground = "#ed985f",
    background = "#8f491b",
    state = function()
        return "-" .. tostring(ctx.stacks * 2) .. " damage"
    end,
    can_stack = true,
    rounds = 2,
    callbacks = {
        on_damage_calc = function()
            if ctx.source == ctx.owner then
                return ctx.damage - ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("STRENGTH", {
    name = "Strength",
    description = "Increases damage for each stack",
    look = "S",
    foreground = "#ed985f",
    background = "#8f491b",
    state = function(ctx)
        return tostring(ctx.stacks * 2) .. " damage"
    end,
    can_stack = true,
    rounds = 2,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage + ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("VULNERABLE", {
    name = "Vulnerable",
    description = "Increases received damage for each stack",
    look = "V",
    foreground = "#ed985f",
    background = "#8f491b",
    state = function(ctx)
        return tostring(ctx.stacks * 2) .. " damage"
    end,
    can_stack = true,
    rounds = 2,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                return ctx.damage + ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("BURN", {
    name = "Burning",
    description = "The enemy burns and receives damage.",
    look = "F",
    foreground = "#ed985f",
    background = "#8f491b",
    state = function(ctx)
        return tostring(ctx.stacks * 4) .. " damage per turn"
    end,
    can_stack = true,
    rounds = 2,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, ctx.owner, ctx.stacks * 2, true)
            return nil
        end,
    }
})