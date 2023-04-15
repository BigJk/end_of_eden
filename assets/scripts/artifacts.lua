register_artifact(
    "DOUBLE_DAMAGE",
    {
        Name = "Stone Of Gigantic Strength",
        Description = "Double all damage dealt.",
        Price = 1000,
        Order = -10,
        Callbacks = {
            OnDamageCalc = function(ctx)
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
            Name = "Repulsion Stone",
            Description = "For each damage taken heal for 2",
            Price = 200,
            Order = 0,
            Callbacks = {
                OnDamage = function(ctx)
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
            Name = "Radiant Seed",
            Description = "A small glowing seed.",
            Price = 50,
            Order = 0,
            Callbacks = {
                OnPickUp = function(ctx)
                    give_card("RADIANT_SEED", ctx.owner)
                    return nil
                end,
            }
        }
);