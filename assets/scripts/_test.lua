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

---assert_card_present asserts that the player's first card is of a certain type, returning an error message if not
---@param id type_id
---@return string|nil
function assert_card_present(id)
    local cards = get_cards(PLAYER_ID)

    if not cards[1] then
        return "Card not in hand"
    end

    local card = get_card_instance(cards[1])
    if card.type_id ~= id then
        return "Card has wrong type: " .. card.type_id
    end

    return nil
end


function assert_cast_card(id, target)
    local cards = get_cards(PLAYER_ID)

    if not cards[1] then
        return "Card not in hand"
    end

    local card = get_card_instance(cards[1])
    if card.type_id ~= id then
        return "Card has wrong type: " .. card.type_id
    end

    if not target then
        cast_card(cards[1])
    else
        cast_card(cards[1], target)
    end
end

---assert_cast_damage asserts that the player's first card deals a certain amount of damage, returning an error message if not
---@param id type_id
---@param dmg number
---@return string|nil
function assert_cast_damage(id, dmg)
    local dummy = add_actor_by_enemy("DUMMY")
    local cards = get_cards(PLAYER_ID)

    if not cards[1] then
        return "Card not in hand"
    end

    local card = get_card_instance(cards[1])
    if card.type_id ~= id then
        return "Card has wrong type: " .. card.type_id
    end

    cast_card(cards[1], dummy)

    if get_actor(dummy).hp ~= 100 - dmg then
        return "Expected " .. tostring(100 - dmg) .. " health, got " .. get_actor(dummy).hp
    end
end

function assert_cast_heal(id, heal)
    local cards = get_cards(PLAYER_ID)

    if not cards[1] then
        return "Card not in hand"
    end

    local card = get_card_instance(cards[1])
    if card.type_id ~= id then
        return "Card has wrong type: " .. card.type_id
    end

    deal_damage(PLAYER_ID, PLAYER_ID, 5, true)
    local hp_before = get_actor(PLAYER_ID).hp
    cast_card(cards[1], PLAYER_ID)

    if get_actor(PLAYER_ID).hp ~= hp_before + heal then
        return "Expected " .. tostring(hp_before + heal) .. " health, got " .. get_actor(PLAYER_ID).hp
    end
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