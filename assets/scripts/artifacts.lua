register_artifact(
    "DOUBLE_DAMAGE",
    {
        name = "Stone Of Gigantic Strength",
        description = "Double all damage dealt.",
        price = 1000,
        order = -10,
        callbacks = {
            on_damage_calc = function(ctx)
                if ctx.target == ctx.owner then
                    return ctx.damage * 2
                end
                return nil
            end,
        }
    }
);

register_artifact(
        "LESSER_DAMAGE_HEAL",
        {
            name = "Repulsion Stone",
            description = "For each damage taken heal for 2",
            price = 200,
            order = 0,
            callbacks = {
                on_damage = function(ctx)
                    if ctx.target == ctx.owner then
                        heal(ctx.owner, 2)
                    end
                    return nil
                end,
            }
        }
);

register_artifact(
        "RADIANT_SEED",
        {
            name = "Radiant Seed",
            description = "A small glowing seed.",
            price = 50,
            order = 0,
            callbacks = {
                on_pick_up = function(ctx)
                    give_card("RADIANT_SEED", ctx.owner)
                    return nil
                end,
            }
        }
);