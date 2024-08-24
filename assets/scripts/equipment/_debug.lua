register_card("DEBUG_INSTA_KILL", {
    name = l("cards.DEBUG_INSTA_KILL.name", "DEBUG Insta Kill"),
    description = l("cards.DEBUG_INSTA_KILL.description", "..."),
    state = function(ctx)
        return "Kill"
    end,
    tags = {},
    max_level = 1,
    color = COLOR_GRAY,
    need_target = true,
    point_cost = 0,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            deal_damage_card(ctx.caster, ctx.guid, ctx.target, 10000)
            return nil
        end
    },
})
