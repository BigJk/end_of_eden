register_card("MELEE_HIT", {
    name = l("cards.MELEE_HIT.name", "Melee Hit"),
    description = l("cards.MELEE_HIT.description", "Use your bare hands to deal 1 (+1 for each upgrade) damage."),
    state = function(ctx)
        return string.format(l("cards.MELEE_HIT.state", "Use your bare hands to deal %s damage."), highlight(1 + ctx.level))
    end,
    tags = { "ATK", "M", "HND" },
    max_level = 1,
    color = COLOR_GRAY,
    need_target = true,
    point_cost = 1,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            deal_damage_card(ctx.caster, ctx.guid, ctx.target, 1 + ctx.level)
            return nil
        end
    },
    test = function ()
        return assert_cast_damage("MELEE_HIT", 1)
    end
})
