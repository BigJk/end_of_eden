register_artifact("RADIANT_SEED", {
    name = "Radiant Seed",
    description = "A small glowing seed.",
    price = 140,
    order = 0,
    callbacks = {
        on_pick_up = function(ctx)
            give_card("RADIANT_SEED", ctx.owner)
            return nil
        end
    }
});