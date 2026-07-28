register_enemy("LASER_DRONE", {
    name = l("enemies.LASER_DRONE.name", "Laser Drone"),
    description = "A drone equipped with a powerful laser cannon.",
    look = [[|o|]],
    color = "#ff0000",
    initial_hp = 7,
    max_hp = 7,
    gold = 40,
    intend = function(ctx)
        if ctx.round % 3 == 0 then
            return l("intents.charge_laser", "Charge up for a powerful laser attack")
        elseif ctx.round % 3 == 1 then
            return l("intents.deal_damage", "Deal ") .. highlight(2) .. " damage"
        else
            return l("intents.deal_damage", "Deal ") .. highlight(5) .. " damage"
        end
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round % 3 == 0 then
                give_status_effect("CHARGING", ctx.guid)
            elseif ctx.round % 3 == 1 then
                deal_damage(ctx.guid, PLAYER_ID, 2)
            else
                deal_damage(ctx.guid, PLAYER_ID, 5)
            end
            return nil
        end
    }
})

register_status_effect("CHARGING", {
    name = l("status_effects.CHARGING.name", "Charging"),
    description = l("status_effects.CHARGING.description", "The drone is charging up for a powerful attack."),
    look = "CHRG",
    foreground = "#ff0000",
    state = function(ctx)
        return l("status_effects.CHARGING.state", "Charging up for a powerful attack.")
    end,
    can_stack = false,
    decay = DECAY_NONE,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage
            end
            return nil
        end
    }
})
