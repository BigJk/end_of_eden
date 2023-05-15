register_enemy("CLEAN_BOT", {
    name = "Cleaning Bot",
    description = "It never stopped cleaning...",
    look = [[ \_/
(* *)
 )#(]],
    color = "#32a891",
    initial_hp = 25,
    max_hp = 25,
    gold = 15,
    intend = function(ctx)
        local self = get_actor(ctx.guid)
        if self.hp <= 8 then
            return "Block " .. highlight(4)
        end

        return "Deal " .. highlight(7) .. " damage"
    end,
    callbacks = {
        on_player_turn = function(ctx)
            local self = get_actor(ctx.guid)

            if self.hp <= 8 then
                give_status_effect("BLOCK", ctx.guid, 4)
            end
        end,
        on_turn = function(ctx)
            local self = get_actor(ctx.guid)

            if self.hp > 8 then
                deal_damage(ctx.guid, PLAYER_ID, 7)
            end

            return nil
        end
    }
})