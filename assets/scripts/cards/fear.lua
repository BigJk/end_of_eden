register_card("FEAR", {
    name = "Fear",
    description = "Inflict " .. highlight("fear") .. " on the target, causing them to miss their next turn.",
    state = function(ctx)
        return nil
    end,
    max_level = 0,
    color = "#725e9c",
    need_target = true,
    point_cost = 2,
    price = 80,
    callbacks = {
        on_cast = function(ctx)
            give_status_effect("FEAR", ctx.target)
            return nil
        end
    }
})