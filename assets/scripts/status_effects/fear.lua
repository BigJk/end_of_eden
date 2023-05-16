register_status_effect("FEAR", {
    name = "Fear",
    description = "Can't act.",
    look = "Fear",
    foreground = "#bb3e03",
    state = function(ctx)
        return "Can't act for " .. highlight(ctx.stacks) .. " turns"
    end,
    can_stack = true,
    decay = DECAY_ONE,
    rounds = 1,
    callbacks = {
        on_turn = function(ctx)
            return true
        end
    }
})
