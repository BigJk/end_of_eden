register_artifact("GIGANTIC_STRENGTH", {
    name = "Stone Of Gigantic Strength",
    description = "Double all damage dealt.",
    price = 250,
    order = 0,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage * 2
            end
            return nil
        end
    },
    test = function()
        dummy = add_actor_by_enemy("DUMMY")

        hp_before = get_actor(dummy).hp
        deal_damage(PLAYER_ID, dummy, 1)
        hp_after = get_actor(dummy).hp

        if hp_after == hp_before - 2 then
            return nil
        end
        return "Damage was not doubled. Before:" .. hp_before .. " After:" .. hp_after
    end
})
