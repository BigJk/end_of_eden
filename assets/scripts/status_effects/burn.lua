register_status_effect("BURN", {
    name = "Burning",
    description = "The enemy burns and receives damage.",
    look = "Brn",
    foreground = "#d00000",
    state = function(ctx)
        return "Takes " .. highlight(ctx.stacks * 4) .. " damage per turn"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, ctx.owner, ctx.stacks * 2, true)
            return nil
        end
    }
})
