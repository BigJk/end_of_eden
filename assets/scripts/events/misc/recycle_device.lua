register_event("RECYCLE_DEVICE", {
    name = "Talking Being",
    description = [[!!artifact_chest.png

...]],
    choices = {
        {
            description_fn = function()
                return "Take Device... (" .. registered.card["RECYCLE"].description .. ")"
            end,
            callback = function(ctx)
                give_card("RECYCLE", PLAYER_ID)
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

register_card("RECYCLE", {
    name = "Recycle",
    description = "Deal " ..
        highlight(12) .. " damage. If " .. highlight("fatal") .. " upgrade random card. " .. highlight("Exhaust") ..
        ".",
    state = function(ctx)
        return nil
    end,
    max_level = 0,
    color = "#d8a448",
    need_target = true,
    exhaust = true,
    point_cost = 2,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            local op_before = #get_opponent_guids(ctx.caster)
            deal_damage(ctx.caster, ctx.target, 12)

            if op_before > #get_opponent_guids(ctx.caster) then
                upgrade_random_card(ctx.caster)
            end

            return nil
        end
    }
})
