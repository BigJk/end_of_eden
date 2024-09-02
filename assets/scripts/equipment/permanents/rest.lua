register_card("REST", {
    name = l("cards.REST.name", "Short Rest"),
    description = l("cards.REST.description", "Heal for 1 (+1 per level)."),
    state = function(ctx)
        return string.format(l("cards.REST.state", "Take a short rest. Heal for %s."), highlight(1 + ctx.level))
    end,
    tags = { "HEAL" },
    max_level = 3,
    color = COLOR_GREEN,
    need_target = false,
    point_cost = 1,
    price = 120,
    callbacks = {
        on_cast = function(ctx)
            heal(ctx.caster, ctx.caster, 1 + ctx.level)
            return nil
        end
    },
})
