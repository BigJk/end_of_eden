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

        print("[ACT_0 ST] history:", get_event_history())

        -- filter out events by id that have already been played
        possible = fun.iter(possible):filter(function(event)
            return event.id == "MERCHANT" or not table.contains(history, event.id)
        end):totable()

        print("[ACT_0 ST] possible:", fun.iter(possible):map(function(e) return e.id end):totable())

        -- fallback for now
        if #possible == 0 then
            possible = find_events_by_tags({ "_ACT_0_FIGHT" })
        end

        local choosen_id = random_int(0, #possible);
        print("[ACT_0 ST] choosen_id:", choosen_id)

        local choosen = possible[1 + choosen_id]
        if choosen ~= nil then
            print("[ACT_0 ST] choosen:", choosen.id)
            set_event(choosen.id)
        end


        return GAME_STATE_EVENT
    end
})
