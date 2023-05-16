register_status_effect("WEAKEN", {
    name = "Weaken",
    description = "Weakens damage for each stack",
    look = "W",
    foreground = "#ed985f",
    state = function()
        return "Deals " .. highlight(ctx.stacks * 2) .. " less damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage - ctx.stacks * 2
            end
            return ctx.damage
        end
    }
})