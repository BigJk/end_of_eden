register_enemy("CLEAN_BOT", {
    name = "Cleaning Bot",
    description = "It never stopped cleaning...",
    look = [[ \_/
(* *)
 )#(]],
    color = "#32a891",
    initial_hp = 13,
    max_hp = 13,
    gold = 15,
    intend = function(ctx)
        local self = get_actor(ctx.guid)
        if self.hp <= 4 then
            return "Block " .. highlight(2)
        end

        return "Deal " .. highlight(2) .. " damage"
    end,
    callbacks = {
        on_player_turn = function(ctx)
            local self = get_actor(ctx.guid)

            if self.hp <= 4 then
                give_status_effect("BLOCK", ctx.guid, 2)
            end
        end,
        on_turn = function(ctx)
            local self = get_actor(ctx.guid)

            if self.hp > 4 then
                deal_damage(ctx.guid, PLAYER_ID, 2)
            end

            return nil
        end
    }
})
