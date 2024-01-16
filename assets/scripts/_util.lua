local function _escape_color(number)
    local escapeString = string.char(27) .. '[%sm'
    return escapeString:format(number)
end

---highlight some value
---@param val any
function highlight(val)
    return text_underline(text_bold("[" .. tostring(val) .. "]"))
end

---highlight_warn some value with warning colors
---@param val any
function highlight_warn(val)
    return  text_underline(text_bold(_escape_color("38;5;161") .. "[" .. tostring(val) .. "]" .. string.char(27) .. "[0m"))
end

---highlight_success some value with success colors
---@param val any
function highlight_success(val)
    return  text_underline(text_bold(_escape_color("38;5;119") .. "[" .. tostring(val) .. "]" .. string.char(27) .. "[0m"))
end

---choose_weighted chooses an item from a list of choices, with a weight for each item.
---@param choices table
---@param weights number[]
---@return string
function choose_weighted(choices, weights)
    print(choices, weights)

    local total_weight = 0
    for _, weight in ipairs(weights) do
        total_weight = total_weight + weight
    end

    local random = math.random() * total_weight
    for i, weight in ipairs(weights) do
        random = random - weight
        if random <= 0 then
            return choices[i]
        end
    end

    return choices[#choices]
end

---table.contains check if a table contains an element.
function table.contains(table, element)
    if table == nil then
        return false
    end
    for _, value in pairs(table) do
        if value == element then
            return true
        end
    end
    return false
end

---find_by_tags find all items with the given tags.
---@param items artifact|card
---@param tags string[]
function find_by_tags(items, tags)
    local found = {}
    for _, item in pairs(items) do
        for _, tag in pairs(tags) do
            if item.tags == nil then
                goto continue
            end
            if not table.contains(item.tags, tag) then
                goto continue
            end
        end

        table.insert(found, item)

        ::continue::
    end
    return found
end

---find_artifacts_by_tags find all artifacts with the given tags.
---@param tags string[]
---@return artifact[]
function find_artifacts_by_tags(tags)
    return find_by_tags(registered.artifact, tags)
end

---find_cards_by_tags find all cards with the given tags.
---@param tags string[]
---@return card[]
function find_cards_by_tags(tags)
    return find_by_tags(registered.card, tags)
end

---find_events_by_tags find all events with the given tags.
---@param tags string[]
---@return event[]
function find_events_by_tags(tags)
    return find_by_tags(registered.event, tags)
end

---choose_weighted_by_price choose a random item from the given list, weighted by price.
---@param items artifact|card
---@return string
function choose_weighted_by_price(items)
    return choose_weighted(
        fun.iter(items):map(function(item) return item.id or item.type_id end):totable(),
        fun.iter(items):map(function(item) return item.price end):totable()
    )
end

---clear_cards_by_tag remove all cards with tag.
---@param tag string tag to remove
---@param excluded? table optional table of guids to exclude.
function clear_cards_by_tag(tag, excluded)
    for _, guid in pairs(get_cards(PLAYER_ID)) do
        if excluded and table.contains(excluded, guid) then
            goto continue
        end

        local tags = get_card(guid).tags
        if table.contains(tags, tag) then
            remove_card(guid)
        end

        ::continue::
    end
end

---clear_artifacts_by_tag remove all artifacts with tag.
---@param tag string tag to remove
---@param excluded table optional table of guids to exclude.
function clear_artifacts_by_tag(tag, excluded)
    for _, guid in pairs(get_artifacts(PLAYER_ID)) do
        if excluded and table.contains(excluded, guid) then
            goto continue
        end

        local tags = get_artifact(guid).tags
        if table.contains(tags, tag) then
            remove_artifact(guid)
        end

        ::continue::
    end
end
