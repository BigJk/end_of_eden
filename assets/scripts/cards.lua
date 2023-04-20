function highlight(dmg)
    return text_underline(text_bold("[" .. tostring(dmg) .. "]"))
end

register_card("KILL",
    {
        name = "Kill",
        description = "Debug Card",
        state = function(ctx)
            return nil
        end,
        max_level = 0,
        color = "#2f3e46",
        need_target = true,
        point_cost = 0,
        price = -1,
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 1000, true)
                return nil
            end,
        }
    }
);

register_card("MELEE_HIT",
    {
        name = "Melee Hit",
        description = "Use your bare hands to deal 5 (+3 for each upgrade) damage.",
        state = function(ctx)
            return "Use your bare hands to deal " .. highlight(5 + ctx.level * 3) .. " damage."
        end,
        max_level = 1,
        color = "#2f3e46",
        need_target = true,
        point_cost = 1,
        price = 30,
        callbacks = {
            on_cast = function(ctx)
                deal_damage(ctx.caster, ctx.target, 5 + ctx.level * 3)
                return nil
            end,
        }
    }
);

register_card("RUPTURE",
    {
        name = "Rupture",
        description = "Inflict your enemy with " .. highlight("Vulnerable") .. ".",
        state = function(ctx)
            return nil
        end,
        max_level = 0,
        color = "#cf532d",
        need_target = true,
        point_cost = 1,
        price = 30,
        callbacks = {
            on_cast = function(ctx)
                give_status_effect("VULNERABLE", ctx.target)
                return nil
            end,
        }
    }
);

register_card("BLOCK",
    {
        name = "Block",
        description = "Shield yourself and gain 5 block.",
        state = function(ctx)
            return "Shield yourself and gain " .. highlight(5 + ctx.level * 3) .. " block."
        end,
        max_level = 1,
        color = "#219ebc",
        need_target = false,
        point_cost = 1,
        price = 40,
        callbacks = {
            on_cast = function(ctx)
                give_status_effect("BLOCK", ctx.caster, 5 + ctx.level * 3)
                return nil
            end,
        }
    }
);

register_card("BLOCK_SPIKES",
    {
        name = "Block Spikes",
        description = "Transforms block to damage.",
        state = function(ctx)
            -- Fetch all BLOCK instances of owner
            local blocks = fun.iter(pairs(get_actor_status_effects(ctx.owner)))
                              :map(get_status_effect_instance)
                              :filter(function(val) return val.type_id == "BLOCK" end)
                              :totable()

            -- Sum stacks to get damage
            local damage = fun.iter(pairs(blocks))
                              :reduce(function(acc, val) return acc + val.stacks end, 0)

            return "Transforms block to " .. highlight(damage) .. " damage."
        end,
        max_level = 0,
        color = "#895cd6",
        need_target = true,
        point_cost = 1,
        price = 100,
        callbacks = {
            on_cast = function(ctx)
                -- Fetch all BLOCK instances of caster
                local blocks = fun.iter(pairs(get_actor_status_effects(ctx.caster)))
                        :map(get_status_effect_instance)
                        :filter(function(val) return val.type_id == "BLOCK" end)
                        :totable()

                -- Sum stacks to get damage
                local damage = fun.iter(pairs(blocks))
                        :reduce(function(acc, val) return acc + val.stacks end, 0)

                if damage == 0 then
                    return "No block status effect present!"
                end

                -- Remove BLOCKs
                fun.iter(pairs(blocks)):for_each(function(val) remove_status_effect(val.guid) end)

                -- Deal Damage
                deal_damage(ctx.caster, ctx.target, damage)

                return nil
            end,
        }
    }
);

register_card("RADIANT_SEED",
    {
        name = "Radiant Seed",
        description = "Inflict 10 (+2 for each upgrade) damage to all enemies, but also causes 5 (-2 for each upgrade) damage to the caster.",
        state = function(ctx)
            return "Inflict " .. highlight(10 + ctx.level * 2) .. " damage to all enemies, but also causes " .. highlight(5 - ctx.level * 2) .. " damage to the caster."
        end,
        max_level = 1,
        color = "#82c93e",
        need_target = false,
        point_cost = 2,
        price = 120,
        callbacks = {
            on_cast = function(ctx)
                -- Deal damage to caster without any modifiers applying
                deal_damage(ctx.caster, ctx.caster, 5 - ctx.level * 2, true)
                -- Deal damage to opponents
                deal_damage_multi(ctx.caster, get_opponent_guids(ctx.caster), 10 + ctx.level * 2)
                return nil
            end,
        }
    }
);