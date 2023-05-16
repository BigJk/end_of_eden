register_status_effect("BLOCK", {
    name = "Block",
    description = "Decreases incoming damage for each stack",
    look = "Blk",
    foreground = "#219ebc",
    state = function(ctx)
        return "Takes " .. highlight(ctx.stacks) .. " less damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                add_status_effect_stacks(ctx.guid, -ctx.damage)
                return ctx.damage - ctx.stacks
            end
            return ctx.damage
        end
    }
})
