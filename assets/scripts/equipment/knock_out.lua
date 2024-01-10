register_card("KNOCK_OUT", {
    name = l("cards.KNOCK_OUT.name", "Knock Out"),
    description = l("cards.KNOCK_OUT.description", "Inflicts " .. highlight("Knock Out") .. " on the target, causing them to miss their next turn."),
    tags = { "CC" },
    max_level = 0,
    color = COLOR_PURPLE,
    need_target = true,
    point_cost = 2,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("KNOCK_OUT", ctx.target)
            return nil
        end
    }
})

register_status_effect("KNOCK_OUT", {
    name = l("status_effects.KNOCK_OUT.name", "Knock Out"),
    description = l("status_effects.KNOCK_OUT.description", "Can't act"),
    look = "K",
    foreground = COLOR_PURPLE,
    state = function(ctx)
        return string.format(l("status_effects.KNOCK_OUT.state", "Can't act for %s turns"), highlight(ctx.stacks))
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
