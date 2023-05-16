register_card("BLOCK", {
    name = "Block",
    description = "Shield yourself and gain 5 " .. highlight("block") .. ".",
    state = function(ctx)
        return "Shield yourself and gain " .. highlight(5 + ctx.level * 3) .. " block."
    end,
    max_level = 1,
    color = "#219ebc",
    need_target = false,
    point_cost = 1,
    price = 40,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("BLOCK", ctx.caster, 5 + ctx.level * 3)
            return nil
        end
    }
})
