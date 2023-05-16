register_card("BLOCK_SPIKES", {
    name = "Block Spikes",
    description = "Transforms " .. highlight("block") .. " to damage.",
    state = function(ctx)
        -- Fetch all BLOCK instances of owner
        local blocks = fun.iter(pairs(get_actor_status_effects(ctx.owner))):map(get_status_effect_instance):filter(function(val)
            return val.type_id == "BLOCK"
        end):totable()

        -- Sum stacks to get damage
        local damage = fun.iter(pairs(blocks)):reduce(function(acc, val)
            return acc + val.stacks
        end, 0)

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
            local blocks = fun.iter(pairs(get_actor_status_effects(ctx.caster))):map(get_status_effect_instance):filter(function(val)
                return val.type_id == "BLOCK"
            end):totable()

            -- Sum stacks to get damage
            local damage = fun.iter(pairs(blocks)):reduce(function(acc, val)
                return acc + val.stacks
            end, 0)

            if damage == 0 then
                return "No block status effect present!"
            end

            -- Remove BLOCKs
            fun.iter(pairs(blocks)):for_each(function(val)
                remove_status_effect(val.guid)
            end)

            -- Deal Damage
            deal_damage(ctx.caster, ctx.target, damage)

            return nil
        end
    }
})