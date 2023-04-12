register_card("MELEE_HIT",
    {
        Name = "Melee Hit",
        Description = "Use your bare hands to deal 2 damage.",
        Color = "#2f3e46",
        NeedTarget = true,
        PointCost = 1,
        Callbacks = {
            OnCast = function(type, guid, caster, target)
                deal_damage(caster, target, 2)
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
                OnCast = function(type, guid, caster, target)
                    -- Deal 1 damage from caster to target
                    deal_damage(caster, target, 1)
                    return nil
                end,
            }
        }
);

register_card("GATHER_HEALTH",
        {
            Name = "Gather Health",
            Description = "Try to heal up a bit and restore 5 hp.",
            Color = "#006400",
            NeedTarget = false,
            PointCost = 2,
            Callbacks = {
                OnCast = function(type, guid, caster, target)
                    heal(caster, caster, 5)
                    return nil
                end,
            }
        }
);