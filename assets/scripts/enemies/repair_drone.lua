REPAIR_DRONE_HEAL = 2

register_enemy("REPAIR_DRONE", {
    name = "Repair Drone",
    description = "A drone designed to repair and support other machines.",
    look = [[]rr[]],
    color = "#00ff00",
    initial_hp = 10,
    max_hp = 10,
    gold = 50,
    intend = function(ctx)
        local opponents = get_opponent_guids(PLAYER_ID)

        -- Check if any opponent needs healing
        for _, opponent_guid in ipairs(opponents) do
            local opponent = get_actor(opponent_guid)
            if opponent_guid ~= ctx.guid and opponent.hp < opponent.max_hp then
                return "Heal " .. highlight(REPAIR_DRONE_HEAL) .. " HP to an ally"
            end
        end

        return "Standby..."
    end,
    callbacks = {
        on_turn = function(ctx)
            local opponents = get_opponent_guids(PLAYER_ID)

            -- Check if any opponent needs healing
            for _, opponent_guid in ipairs(opponents) do
                local opponent = get_actor(opponent_guid)
                if opponent_guid ~= ctx.guid and opponent.hp < opponent.max_hp then
                    heal(ctx.guid, opponent_guid, REPAIR_DRONE_HEAL)
                    return nil
                end
            end

            return nil
        end
    }
})
