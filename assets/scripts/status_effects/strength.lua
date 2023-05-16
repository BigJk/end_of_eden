register_status_effect("STRENGTH", {
    name = "Strength",
    description = "Increases damage for each stack",
    look = "Str",
    foreground = "#d00000",
    state = function(ctx)
        return "Deal " .. highlight(ctx.stacks) .. " more damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage + ctx.stacks
            end
            return ctx.damage
        end
    }
})