function add_found_artifact_event(id, picture, description, choice_description, tags)
    register_event(id, {
        name = l("ui.found_template", "Found: ") .. registered.artifact[id].name,
        description = string.format(l("ui.found_event_template", "!!%s\n\n**You found something!** %s"), picture or "artifact_chest.jpg", description),
        tags = tags,
        choices = {
            {
                description_fn = function()
                    return l("ui.take_template", "Take ") .. registered.artifact[id].name .. "... (" .. choice_description .. ")"
                end,
                callback = function(ctx)
                    give_artifact(id, PLAYER_ID)
                    return nil
                end
            },
            {
                description = l("ui.leave", "Leave..."),
                callback = function()
                    return nil
                end
            }
        },
        on_end = function()
            return GAME_STATE_RANDOM
        end
    })
end