register_story_teller("STAGE_3", {
    active = function(ctx)
        if had_events_any(stage_1_init_events) and get_stages_cleared() > 20 then
            return 3
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 30 then
            -- BOSS
        end

        add_actor_by_enemy("DUMMY")

        return GAME_STATE_FIGHT
    end
})
