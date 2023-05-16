register_card("RUPTURE", {
    name = "Rupture",
    description = "Inflict your enemy with " .. highlight("Vulnerable") .. ".",
    state = function(ctx)
        return "Inflict your enemy with " .. highlight(tostring(1 + ctx.level) .. " Vulnerable") .. "."
    end,
    max_level = 3,
    color = "#cf532d",
    need_target = true,
    point_cost = 1,
    price = 30,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("VULNERABLE", ctx.target, 1 + ctx.level)
            return nil
        end
    }
})
