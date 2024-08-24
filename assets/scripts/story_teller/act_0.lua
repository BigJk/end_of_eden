register_story_teller("_ACT_0", {
    active = function()
        -- if #get_event_history() <= 6 then
        --     return 1
        -- end
        --
        -- Keep active for now
        return 1
    end,
    decide = function()
        local history = get_event_history()
        local possible = {}

        -- every 3 events, play a non-combat event
        local events = #history - 1;
        if events ~= 0 and events % 2 == 0 then
            possible = find_events_by_tags({ "_ACT_0" })
        else
            possible = find_events_by_tags({ "_ACT_0_FIGHT" })
        end

        print(#get_event_history())

        -- filter out events by id that have already been played
        possible = fun.iter(possible):filter(function(event)
            return event == "MERCHANT" or not table.contains(history, event.id)
        end):totable()

        -- fallback for now
        if #possible == 0 then
            possible = find_events_by_tags({ "_ACT_0_FIGHT" })
        end
        set_event(possible[math.random(#possible)].id)

        -- if we cleared a stage, give the player a random artifact
        local last_stage_count = fetch("last_stage_count")
        local current_stage_count = get_stages_cleared()
        if last_stage_count ~= current_stage_count then
            local gets_random_artifact = math.random() < 0.25

            if gets_random_artifact then
                local player_artifacts = fun.iter(get_actor(PLAYER_ID).artifacts):map(function(id)
                    return get_artifact(id).id
                end):totable()
                local artifacts = find_artifacts_by_tags({ "_ACT_0" })
                if #artifacts > 0 then
                    local artifact = choose_weighted_by_price(artifacts)
                    if not table.contains(player_artifacts, artifact) then
                        give_artifact(PLAYER_ID, artifact)
                    end
                end
            end
        end

        return GAME_STATE_EVENT
    end
})
