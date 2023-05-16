register_artifact("REPULSION_STONE", {
    name = "Repulsion Stone",
    description = "For each damage taken heal for 2",
    price = 100,
    order = 0,
    callbacks = {
        on_damage = function(ctx)
            if ctx.target == ctx.owner then
                heal(ctx.owner, 2)
            end
            return nil
        end
    }
});
