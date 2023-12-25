---assert_chain executes a chain of tests, returning the first non-nil result
---@param tests function[]
---@return string|nil
function assert_chain(tests)
    for i, test in ipairs(tests) do
        local result = test()
        if result ~= nil then
            return result
        end
    end
    return nil
end

---assert_status_effect_count asserts that the player has a certain number of status effects, returning an error message if not
---@param count number
---@return string|nil
function assert_status_effect_count(count)
    local status_effects = get_actor_status_effects(PLAYER_ID)

    -- check if length of status_effects is 1
    if #status_effects ~= count then
        return "Expected " .. count .. " status effects, got " .. #status_effects
    end

    return nil
end

---assert_status_effect asserts that the player has a status effect of a certain type, returning an error message if not
---@param type type_id
---@param number number
---@return string|nil
function assert_status_effect(type, number)
    local status_effects = get_actor_status_effects(PLAYER_ID)

    -- find the status effect
    for i, guid in ipairs(status_effects) do
        local instance = get_status_effect_instance(guid)
        if instance.type_id == type then
            -- check if the stacks are equal to the number
            if instance.stacks ~= number then
                return "Expected " .. number .. " block, got " .. tostring(instance.stacks)
            end

            return nil
        end
    end

    return "Status effect not found"
end