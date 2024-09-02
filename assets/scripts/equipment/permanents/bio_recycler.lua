register_artifact("BIO_RECYCLER", {
    name = "Bio Recycler",
    description = "Heal 1 on kill.",
    tags = { "_ACT_0" },
    price = 200,
    order = 0,
    callbacks = {
        on_actor_die = function(ctx)
            if ctx.source == PLAYER_ID then
                heal(PLAYER_ID, PLAYER_ID, 1)
            end
            return nil
        end
    },
    test = function()
        local dummy = add_actor_by_enemy("DUMMY")
        deal_damage(dummy, PLAYER_ID, 1, true)
        local hp_before = get_player().hp
        deal_damage(PLAYER_ID, dummy, 100)
        local hp_after = get_player().hp
        if hp_after - hp_before ~= 1 then
            return "Expected 1 HP heal, got " .. (hp_after - hp_before)
        end
    end
})
