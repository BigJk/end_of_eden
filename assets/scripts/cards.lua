register_card("MELEE_HIT",
    {
        Name = "Melee Hit",
        Description = "Use your bare hands to deal 2 damage.",
        Color = "#2f3e46",
        NeedTarget = true,
        PointCost = 1,
        Callbacks = {
            OnCast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 2)
                return nil
            end,
        }
    }
);

register_card("SLICE",
        {
            Name = "Slice",
            Description = "Try to inflict a wound on the enemy.",
            Color = "#cf532d",
            NeedTarget = true,
            PointCost = 1,
            Callbacks = {
                OnCast = function(ctx)
                    give_status_effect("VULNERABLE", ctx.target)
                    return nil
                end,
            }
        }
);

register_card("BITE",
        {
            Name = "Bite",
            Description = "Nom nom...",
            Color = "#2f3e46",
            NeedTarget = true,
            PointCost = 1,
            Callbacks = {
                OnCast = function(ctx)
                    -- Deal 1 damage from caster to target
                    deal_damage(ctx.caster, ctx.target, 1)
                    return nil
                end,
            }
        }
);

register_card("RADIANT_SEED",
        {
            Name = "Radiant Seed",
            Description = "Inflict 10 damage to all enemies, but also causes 5 damage to the caster.",
            Color = "#82c93e",
            NeedTarget = false,
            PointCost = 2,
            Callbacks = {
                OnCast = function(ctx)
                    -- Deal 5 damage to caster without any modifiers applying
                    deal_damage(ctx.caster, ctx.caster, 5, true)
                    -- Deal 10 damage to opponents
                    deal_damage_multi(ctx.caster, get_opponent_guids(ctx.caster), 10)
                    return nil
                end,
            }
        }
);