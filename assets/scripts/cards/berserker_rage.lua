register_card("BERSERKER_RAGE", {
    name = "Berserker Rage",
    description = "Gain " .. highlight("3 action points") .. ", but take 30% (-10% per level) of your HP as damage.",
    state = function(ctx)
        return "Gain " .. highlight("3 action points") .. ", but take " .. highlight(tostring(30 - ctx.level * 10) .. "%") .. " (" ..
                   tostring(get_player().hp * (0.3 - ctx.level * 0.1)) .. ") of your HP as damage."
    end,
    max_level = 0,
    color = "#d8a448",
    need_target = false,
    point_cost = 0,
    price = 100,
    callbacks = {
        on_cast = function(ctx)
            player_give_action_points(3)
            deal_damage(ctx.caster, ctx.caster, get_player().hp * (0.3 - ctx.level * 0.1), true)
            return nil
        end
    }
})
