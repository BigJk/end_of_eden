 register_story_teller("ACT_0", {
    active = function()
        if #get_event_history() < 5 then
            return 1
        end
        return 0
    end,
    decide = function()
        local possible = find_events_by_tags({"ACT_0"})
        local history = get_event_history()

        -- filter out events by id that have already been played
        possible = fun.iter(possible):filter(function(event)
            return not table.contains(history, event.id)
        end):totable()

        -- fallback for now
        if #possible == 0 then
            possible = find_events_by_tags({"ACT_0"})
        end

        set_event(possible[math.random(#possible)].id)

        return GAME_STATE_EVENT
    end
})
