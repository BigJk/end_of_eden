stage_1_the_wasteland = { fights = { { "SAND_STALKER" }, { "SAND_STALKER", "SAND_STALKER" } } }

register_story_teller("STAGE_1", {
    active = function(ctx)
        if had_event("THE_WASTELAND") then
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

        local choice = stage_1_the_wasteland.fights[math.random(#stage_1_the_wasteland.fights)]
        for _, v in ipairs(choice) do
            add_actor_by_enemy(v)
        end

        return GAME_STATE_FIGHT
    end
})
