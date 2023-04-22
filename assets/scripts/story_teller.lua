each(print, map(function(key, val)
    return val.description
end, registered.card))

--
-- Pre-Stage
--

register_story_teller("PRE_STAGE", {
    active = function(ctx)
        if not had_event("FIRST_OUTSIDE") then
            return 1
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage >= 3 then
            -- If we didn't skip the pre-stage we get another artifact
            set_event(create_artifact_choice({ random_artifact(get_merchant_gold_max()), random_artifact(get_merchant_gold_max()) }, {
                on_end = function()
                    set_event("FIRST_OUTSIDE")
                    return GAME_STATE_EVENT
                end
            }))

            return GAME_STATE_EVENT
        end

        -- Fight against rust mites or clean bots
        local d = math.random(2)
        if d == 1 then
            add_actor_by_enemy("RUST_MITE")
        elseif d == 2 then
            add_actor_by_enemy("CLEAN_BOT")
        end

        return GAME_STATE_FIGHT
    end
})

--
-- Stage 1
--

register_story_teller("STAGE_1", {
    active = function(ctx)
        if had_event("FIRST_OUTSIDE") then
            return 1
        end
        return 0
    end,
    decide = function(ctx)
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
    active = function(ctx)
        if had_event("FIRST_OUTSIDE") and get_stages_cleared() > 10 then
            return 2
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 20 then
            -- BOSS
        end

        return nil
    end
})

--
-- Stage 3
--

register_story_teller("STAGE_3", {
    active = function(ctx)
        if had_event("FIRST_OUTSIDE") and get_stages_cleared() > 20 then
            return 3
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 30 then
            -- BOSS
        end

        return nil
    end
})