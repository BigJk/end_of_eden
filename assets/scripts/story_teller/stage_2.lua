register_story_teller("STAGE_2", {
    active = function(ctx)
        if had_events_any(stage_1_init_events) and get_stages_cleared() > 10 then
            return 2
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 20 then
            -- BOSS
        end

        return GAME_STATE_FIGHT
    end
})
