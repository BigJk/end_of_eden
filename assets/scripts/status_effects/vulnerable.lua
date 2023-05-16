register_status_effect("VULNERABLE", {
    name = "Vulnerable",
    description = "Increases received damage for each stack",
    look = "Vur",
    foreground = "#ffba08",
    state = function(ctx)
        return "Takes " .. highlight(ctx.stacks * 25) .. "% more damage"
    end,
    can_stack = true,
    decay = DECAY_ONE,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                return ctx.damage * (1.0 + 0.25 * ctx.stacks)
            end
            return ctx.damage
        end
    }
})