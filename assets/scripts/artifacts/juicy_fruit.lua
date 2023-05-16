register_artifact("JUICY_FRUIT", {
    name = "Juicy Fruit",
    description = "Tastes good and boosts your HP.",
    price = 80,
    order = 0,
    callbacks = {
        on_pick_up = function(ctx)
            actor_add_max_hp(ctx.owner, 10)
            return nil
        end
    }
});