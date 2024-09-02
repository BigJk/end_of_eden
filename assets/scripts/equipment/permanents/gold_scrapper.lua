register_artifact("GOLD_SCRAPPER", {
    name = "Gold Scrapper",
    description = "Gain 15 gold on kill.",
    tags = { "_ACT_0" },
    price = 200,
    order = 0,
    callbacks = {
        on_actor_die = function(ctx)
            if ctx.source == PLAYER_ID then
                give_player_gold(15)
            end
            return nil
        end
    },
    test = function()
        local dummy = add_actor_by_enemy("DUMMY")
        deal_damage(PLAYER_ID, dummy, 100)
        if get_player().gold ~= 15 then
            return "Expected 15 gold, got " .. get_player().gold
        end
    end
})
