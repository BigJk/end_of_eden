register_card("COMBINED_SHOT", {
    name = "Combined Shot",
    description = "Deal " .. highlight(5) .. " (+5 for each level) damage for each enemy.",
    state = function(ctx)
        return "Deal " .. highlight((5 + ctx.level * 5) * #get_opponent_guids(ctx.owner)) .. " damage for each enemy."
    end,
    max_level = 1,
    color = "#d8a448",
    need_target = true,
    point_cost = 1,
    price = 150,
    callbacks = {
        on_cast = function(ctx)
            deal_damage(ctx.caster, ctx.target, (5 + ctx.level * 5) * #get_opponent_guids(ctx.owner))
            return nil
        end
    }
})
