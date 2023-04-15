register_status_effect("WEAKEN", {
    Name = "Weaken",
    Description = "Weakens damage for each stack",
    Look = "W",
    Foreground = "#ed985f",
    Background = "#8f491b",
    State = function()
        return "-" .. tostring(ctx.stacks * 2) .. " damage"
    end,
    CanStack = true,
    Rounds = 2,
    Callbacks = {
        OnDamageCalc = function()
            if ctx.source == ctx.owner then
                return ctx.damage - ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("STRENGTH", {
    Name = "Strength",
    Description = "Increases damage for each stack",
    Look = "S",
    Foreground = "#ed985f",
    Background = "#8f491b",
    State = function(ctx)
        return tostring(ctx.stacks * 2) .. " damage"
    end,
    CanStack = true,
    Rounds = 2,
    Callbacks = {
        OnDamageCalc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage + ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("VULNERABLE", {
    Name = "Vulnerable",
    Description = "Increases received damage for each stack",
    Look = "V",
    Foreground = "#ed985f",
    Background = "#8f491b",
    State = function(ctx)
        return tostring(ctx.stacks * 2) .. " damage"
    end,
    CanStack = true,
    Rounds = 2,
    Callbacks = {
        OnDamageCalc = function(ctx)
            if ctx.target == ctx.owner then
                return ctx.damage + ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("BURN", {
    Name = "Burning",
    Description = "The enemy burns and receives damage.",
    Look = "F",
    Foreground = "#ed985f",
    Background = "#8f491b",
    State = function(ctx)
        return tostring(ctx.stacks * 4) .. " damage per turn"
    end,
    CanStack = true,
    Rounds = 2,
    Callbacks = {
        OnTurn = function(ctx)
            deal_damage(ctx.guid, ctx.owner, ctx.stacks * 2, true)
            return nil
        end,
    }
})