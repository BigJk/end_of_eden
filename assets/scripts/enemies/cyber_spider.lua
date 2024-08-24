register_enemy("CYBER_SPIDER", {
    name = "CYBER Spider",
    description = "It waits for its prey to come closer",
    look = [[/\o^o/\]],
    color = "#ff4d6d",
    initial_hp = 8,
    max_hp = 8,
    gold = 40,
    intend = function(ctx)
        if ctx.round > 0 and ctx.round % 3 == 0 then
            return "Deal " .. highlight(5) .. " damage"
        end

        return "Wait..."
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round > 0 and ctx.round % 3 == 0 then
                deal_damage(ctx.guid, PLAYER_ID, 5)
            end

            return nil
        end
    }
})
