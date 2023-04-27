each(print, map(function(key, val)
    return val.description
end, registered.card))

--
-- Pre-Stage
--

stage_1_init_events = { "THE_CORE", "BIO_KINGDOM", "THE_WASTELAND" }

register_story_teller("PRE_STAGE", {
    active = function(ctx)
        if not had_events_any(stage_1_init_events) then
            return 1
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage >= 3 then
            -- If we didn't skip the pre-stage we get another artifact
            set_event(create_artifact_choice({ random_artifact(get_merchant_gold_max()), random_artifact(get_merchant_gold_max()) }, {
                description = [[As you explore the abandoned cryo facility, a feeling of dread washes over you. The facility is eerily quiet, with malfunctioning computers and flickering lights being the only signs of life. As you move through the winding corridors, you stumble upon a hidden door. It's almost as if the facility itself is trying to keep you from finding what lies beyond.

After some effort, you manage to open the door and find yourself in a small room. The room is dark, and you can barely make out a chest in the center of the room. As you approach it, the feeling of unease grows stronger. What secret artifact could be hidden inside this chest? Is it something that will aid you on your journey or something more sinister? You take a deep breath, steeling yourself for whatever you may find inside, and reach for the lid...]],
                on_end = function()
                    set_event(stage_1_init_events[math.random(#stage_1_init_events)])
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

stage_2 = {
    fights = {
        { "RUST_MITE", "RUST_MITE", "RUST_MITE" },
        { "SHADOW_ASSASSIN", "SHADOW_ASSASSIN" },
        { "SHADOW_ASSASSIN" }
    }
}

register_story_teller("STAGE_1", {
    active = function(ctx)
        if had_events_any(stage_1_init_events) then
            return 1
        end
        return 0
    end,
    decide = function(ctx)
        local stage = get_stages_cleared()

        if stage == 10 then
            -- BOSS
        end

        -- 10% chance to find a random artifact
        if math.random() < 0.1 then
            set_event(create_artifact_choice({ random_artifact(get_merchant_gold_max()), random_artifact(get_merchant_gold_max()) }))
        end

        local choice = stage_2.fights[math.random(#stage_2.fights)]
        for _, v in ipairs(choice) do
            add_actor_by_enemy(v)
        end

        return GAME_STATE_FIGHT
    end
})

--
-- Stage 2
--

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

        return nil
    end
})

--
-- Stage 3
--

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

        return nil
    end
})