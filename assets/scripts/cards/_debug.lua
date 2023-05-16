register_card("KILL", {
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
        end
    }
})