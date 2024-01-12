register_card("STIM_PACK", {
    name = l("cards.STIM_PACK.name", "Stim Pack"),
    description = l("cards.STIM_PACK.description", highlight("One-Time") .. "\n\nRestores " .. highlight(5) .. " HP."),
    tags = { "HEAL" },
    max_level = 0,
    color = COLOR_BLUE,
    need_target = false,
    does_consume = true,
    point_cost = 0,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            heal(ctx.caster, ctx.caster, 5)
            return nil
        end
    },
    test = function ()
        return assert_cast_heal("STIM_PACK", 5)
    end
})