register_enemy("CYBER_SLIME_MINION", {
    name = "Cyber Slime Offspring",
    description = "A smaller version of the Cyber Slime.",
    look = [[ o ]],
    color = "#00ff00",
    initial_hp = 4,
    max_hp = 4,
    gold = 10,
    intend = function(ctx)
        return "Deal " .. highlight(1) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, PLAYER_ID, 1)
            return nil
        end
    }
})

register_enemy("CYBER_SLIME", {
    name = "Cyber Slime",
    description = "A cybernetic slime that splits into smaller slimes when defeated.",
    look = [[ (O) ]],
    color = "#00ff00",
    initial_hp = 10,
    max_hp = 10,
    gold = 50,
    intend = function(ctx)
        return "Deal " .. highlight(2) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, PLAYER_ID, 2)
            return nil
        end,
        on_actor_die = function(ctx)
            if get_actor(ctx.target).type_id ~= "CYBER_SLIME" then
                return nil
            end

            add_actor_by_enemy("CYBER_SLIME_MINION")
            add_actor_by_enemy("CYBER_SLIME_MINION")
            if math.random() < 0.25 then
                add_actor_by_enemy("CYBER_SLIME_MINION")
            end
            return nil
        end
    }
})
