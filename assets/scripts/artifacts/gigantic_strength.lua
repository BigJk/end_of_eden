register_artifact("GIGANTIC_STRENGTH", {
    name = "Stone Of Gigantic Strength",
    description = "Double all damage dealt.",
    price = 250,
    order = 0,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                return ctx.damage * 2
            end
            return nil
        end
    }
});
