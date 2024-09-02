register_enemy("NANOBOT_SWARM", {
    name = l("enemies.NANOBOT_SWARM.name", "Nanobot Swarm"),
    description = l("enemies.NANOBOT_SWARM.description", "A growing swarm of nanobots."),
    look = ".*#.-",
    color = "#9b5de5",
    initial_hp = 15,
    max_hp = 100,
    gold = 50,
    intend = function(ctx)
        return "Deal " ..
        highlight(simulate_deal_damage(ctx.guid, PLAYER_ID, ctx.round + 1)) .. " damage. Heal " .. highlight(1) .. " HP."
    end,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, PLAYER_ID, ctx.round + 1)
            heal(ctx.guid, ctx.guid, 1)

            return nil
        end
    }
})
