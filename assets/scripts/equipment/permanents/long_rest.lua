register_card("LONG_REST", {
    name = l("cards.LONG_REST.name", "Long Rest"),
    description = l("cards.REST.description", "Heal for 4 (+1 per level)."),
    state = function(ctx)
        return string.format(l("cards.REST.state", "Take a short rest. Heal for %s."), highlight(4 + ctx.level))
    end,
    tags = { "HEAL" },
    max_level = 2,
    color = COLOR_GREEN,
    need_target = false,
    does_exhaust = true,
    point_cost = 3,
    price = 300,
    callbacks = {
        on_cast = function(ctx)
            heal(ctx.caster, ctx.caster, 4 + ctx.level)
            return nil
        end
    },
})
