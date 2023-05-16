register_artifact("GOLD_CONVERTER", {
    name = "Gold Converter",
    description = "Gain 10 extra gold for each killed enemy.",
    price = 50,
    order = 0,
    callbacks = {
        on_actor_die = function(ctx)
            if ctx.owner == PLAYER_ID and ctx.owner == ctx.source then
                give_player_gold(10)
            end
            return nil
        end
    }
});