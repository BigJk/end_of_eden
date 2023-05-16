register_card("SHIELD_BASH", {
    name = "Shield Bash",
    description = "Deal 4 (+2 for each upgrade) damage to the enemy and gain " .. highlight("block") ..
        " status effect equal to the damage dealt.",
    state = function(ctx)
        return "Deal " .. highlight(4 + ctx.level * 2) .. " damage to the enemy and gain " .. highlight("block") ..
                   " status effect equal to the damage dealt."
    end,
    max_level = 1,
    color = "#ff5722",
    need_target = true,
    point_cost = 1,
    price = 40,
    callbacks = {
        on_cast = function(ctx)
            local damage = deal_damage(ctx.caster, 4 + ctx.level * 2)
            give_status_effect("BLOCK", ctx.caster, damage)
            return nil
        end
    }
})
