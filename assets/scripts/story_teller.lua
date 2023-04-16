--
-- Pre-Stage
--

register_story_teller("PRE_STAGE", {
    Active = function(ctx)
        if not had_event("FIRST_OUTSIDE") then
            return 1
        end
        return 0
    end,
    Decide = function(ctx)
        local stage = get_stages_cleared()

        if stage > 3 then
            set_event("FIRST_OUTSIDE")
            return GAME_STATE_EVENT
        end

        add_actor_by_enemy("RUST_MITE")
        add_actor_by_enemy("RUST_MITE")
        add_actor_by_enemy("RUST_MITE")

        return GAME_STATE_FIGHT
    end
})

--
-- Stage 1
--

register_story_teller("STAGE_1", {
    Active = function(ctx)
        if had_event("FIRST_OUTSIDE") then
            return 1
        end
        return 0
    end,
    Decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 10 then
            -- BOSS
        end

        return nil
    end
})

--
-- Stage 2
--

register_story_teller("STAGE_2", {
    Active = function(ctx)
        if had_event("FIRST_OUTSIDE") and get_stages_cleared() > 10 then
            return 2
        end
        return 0
    end,
    Decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 20 then
            -- BOSS
        end

        return nil
    end
})