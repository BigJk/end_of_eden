function assert_chain(tests)
    for i, test in ipairs(tests) do
        result = test()
        if result ~= nil then
            return result
        end
    end
    return nil
end

function assert_status_effect_count(count)
    status_effects = get_actor_status_effects(PLAYER_ID)

    -- check if length of status_effects is 1
    if #status_effects ~= count then
        return "Expected " .. count .. " status effects, got " .. #status_effects
    end

    return nil
end

function assert_status_effect(type, number)
    status_effects = get_actor_status_effects(PLAYER_ID)

    -- find the status effect
    for i, guid in ipairs(status_effects) do
        instance = get_status_effect_instance(guid)
        if instance.type_id == type then
            -- check if the stacks are equal to the number
            if instance.stacks ~= number then
                return "Expected " .. number .. " block, got " .. tostring(block.stacks)
            end

            return nil
        end
    end

    return "Status effect not found"
end