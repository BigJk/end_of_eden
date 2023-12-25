register_event("RAISING_THE_BAR", {
    name = "Raising The Bar",
    description = [[!!red_room.png

...]],
    choices = {
        {
            description_fn = function()
                return "Take Crowbar... (" .. registered.card["CROWBAR"].description .. ")"
            end,
            callback = function(ctx)
                give_card("CROWBAR", PLAYER_ID)
                return nil
            end
        }, {
        description = "Leave...",
        callback = function()
            return nil
        end
    }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_card("CROWBAR", {
    name = "Crowbar",
    description = "Deal " .. highlight(22) .. " damage.",
    state = function(ctx)
        return nil
    end,
    max_level = 0,
    color = "#f37b21",
    need_target = true,
    exhaust = true,
    point_cost = 3,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            deal_damage(ctx.caster, ctx.target, 22)
            return nil
        end
    }
})
