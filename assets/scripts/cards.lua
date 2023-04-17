register_card("MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 2 damage.",
        color = "#2f3e46",
        need_target = true,
        point_cost = 1,
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 2)
                return nil
            end,
        }
    }
);

register_card("SLICE",
        {
            name = "Slice",
            description = "Try to inflict a wound on the enemy.",
            color = "#cf532d",
            need_target = true,
            point_cost = 1,
            callbacks = {
                on_cast = function(ctx)
                    give_status_effect("VULNERABLE", ctx.target)
                    return nil
                end,
            }
        }
);

register_card("BITE",
        {
            name = "Bite",
            description = "Nom nom...",
            color = "#2f3e46",
            need_target = true,
            point_cost = 1,
            callbacks = {
                on_cast = function(ctx)
                    -- Deal 1 damage from caster to target
                    deal_damage(ctx.caster, ctx.target, 1)
                    return nil
                end,
            }
        }
);

register_card("RADIANT_SEED",
        {
            name = "Radiant Seed",
            description = "Inflict 10 damage to all enemies, but also causes 5 damage to the caster.",
            color = "#82c93e",
            need_target = false,
            point_cost = 2,
            callbacks = {
                on_cast = function(ctx)
                    -- Deal 5 damage to caster without any modifiers applying
                    deal_damage(ctx.caster, ctx.caster, 5, true)
                    -- Deal 10 damage to opponents
                    deal_damage_multi(ctx.caster, get_opponent_guids(ctx.caster), 10)
                    return nil
                end,
            }
        }
);