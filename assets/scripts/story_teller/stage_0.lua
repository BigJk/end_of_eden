stage_1_init_events = { "THE_CORE", "BIO_KINGDOM", "THE_WASTELAND" }

register_story_teller("STAGE_0", {
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