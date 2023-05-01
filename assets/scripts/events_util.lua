function create_artifact_choice(artifacts, options)
    local id = guid()
    local choices = fun.iter(artifacts):map(function(type)
        local art = get_artifact(type)

        return {
            description = "Take " .. text_bold(art.name) .. text_italic(" (" .. art.description .. ")"),
            callback = function()
                give_artifact(type, PLAYER_ID)
                return nil
            end
        }
    end):totable()
    choices[#choices + 1] = {
        description = "Skip...",
        callback = function()
            return nil
        end
    }

    local def = {
        name = "Artifact",
        description = [[As you journey through the desolate land, you come across a hidden cache. Inside, you find an array of strange and wondrous artifacts, each with their own mysterious powers. You know that choosing just one could change the course of your journey forever. As you examine each item, you feel the weight of the responsibility resting on your shoulders...]],
        choices = choices,
        on_end = function()
            return GAME_STATE_RANDOM
        end
    }

    if options ~= nil then
        if options["name"] ~= nil then
            def.name = options["name"]
        end

        if options["description"] ~= nil then
            def.description = options["description"]
        end

        if options["on_end"] ~= nil then
            def.on_end = options["on_end"]
        end
    end

    register_event(id, def)

    return id
end

function create_card_choice(cards, options)
    local id = guid()
    local choices = fun.iter(cards):map(function(type)
        local art = get_artifact(type)

        return {
            description = "Take " .. text_bold(art.name) .. text_italic(" (" .. art.description .. ")"),
            callback = function()
                give_card(type, PLAYER_ID)
                return nil
            end
        }
    end):totable()
    choices[#choices + 1] = {
        description = "Skip...",
        callback = function()
            return nil
        end
    }

    local def = {
        name = "Cards",
        description = [[As you journey through the desolate land, you come across a hidden cache. Inside, you find an array of strange and wondrous artifacts, each with their own mysterious powers. You know that choosing just one could change the course of your journey forever. As you examine each item, you feel the weight of the responsibility resting on your shoulders...]],
        choices = choices,
        on_end = function()
            return GAME_STATE_RANDOM
        end
    }

    if options ~= nil then
        if options["name"] ~= nil then
            def.name = options["name"]
        end

        if options["description"] ~= nil then
            def.description = options["description"]
        end

        if options["on_end"] ~= nil then
            def.on_end = options["on_end"]
        end
    end

    register_event(id, def)

    return id
end
