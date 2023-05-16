register_enemy("SAND_STALKER", {
    name = "Sand Stalker",
    description = "It waits for its prey to come closer.",
    look = "( ° ° )",
    color = "#8e4028",
    initial_hp = 25,
    max_hp = 25,
    gold = 20,
    intend = function(ctx)
        if ctx.round % 4 == 0 then
            return "Weaken your resolve"
        end

        return "Deal " .. highlight(7) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round % 4 == 0 then
                if deal_damage(ctx.guid, PLAYER_ID, 5) > 0 then
                    give_status_effect("WEAKEN", PLAYER_ID, 1)
                end
            else
                deal_damage(ctx.guid, PLAYER_ID, 7)
            end

            return nil
        end
    }
})
