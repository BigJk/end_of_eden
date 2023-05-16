register_artifact("SPIKED_PLANT", {
    name = "Spiked Plant",
    description = "Deal 2 damage back to enemy attacks.",
    price = 50,
    order = 0,
    callbacks = {
        on_damage = function(ctx)
            if ctx.source ~= ctx.owner and ctx.owner == ctx.target then
                deal_damage(ctx.owner, ctx.source, 2)
            end
            return nil
        end
    }
});