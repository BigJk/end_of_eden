register_card("RADIANT_SEED", {
    name = "Radiant Seed",
    description = "Inflict 10 (+2 for each upgrade) damage to all enemies, but also causes 5 (-2 for each upgrade) damage to the caster.",
    state = function(ctx)
        return "Inflict " .. highlight(10 + ctx.level * 2) .. " damage to all enemies, but also causes " .. highlight(5 - ctx.level * 2) ..
                " damage to the caster."
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
        end
    }
})