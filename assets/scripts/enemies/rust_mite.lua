register_enemy("RUST_MITE", {
    name = l("enemies.RUST_MITE.name", "Rust Mite"),
    description = l("enemies.RUST_MITE.description", "A small robot that eats metal."),
    look = "/v\\",
    color = "#e6e65a",
    initial_hp = 12,
    max_hp = 12,
    gold = 35,
    intend = function(ctx)
        if ctx.round % 4 == 0 then
            return "Load battery"
        end

        return "Deal " .. highlight(simulate_deal_damage(ctx.guid, PLAYER_ID, 1)) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round % 4 == 0 then
                give_status_effect("CHARGED", ctx.guid)
            else
                deal_damage(ctx.guid, PLAYER_ID, 1)
            end

            return nil
        end
    }
})

register_status_effect("CHARGED", {
    name = l("status_effects.CHARGED.name", "Charged"),
    description = l("status_effects.CHARGED.description", "Attacks will deal more damage per stack."),
    look = "CHRG",
    foreground = "#207BE7",
    state = function(ctx)
        return string.format(l("status_effects.CHARGED.state", "Attacks deal %s more damage"), highlight(ctx.stacks * 1))
    end,
    can_stack = true,
    decay = DECAY_NONE,
    rounds = 0,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage + 1 * ctx.stacks
            end
            return nil
        end
    }
})
