register_card("MELEE_HIT", {
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
        end
    }
})
