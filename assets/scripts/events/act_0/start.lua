register_event("START", {
    name = "Waking up...",
    description = [[!!cryo_start.jpg

You wake up in a dimly lit room, the faint glow of a red emergency light casting an eerie hue over the surroundings. The air is musty and stale, the metallic scent of the cryo-chamber still lingering in your nostrils. You feel groggy and disoriented, your mind struggling to process what's happening.

As you try to sit up, you notice that your body is stiff and unresponsive. It takes a few moments for your muscles to warm up and regain their strength. Looking around, you see that the walls are made of a dull gray metal, covered in scratches and scuff marks. There's a faint humming sound coming from somewhere, indicating that the facility is still operational.

You try to remember how you ended up here, but your memories are hazy and fragmented. The last thing you recall is a blinding flash of light and a deafening boom. You must have been caught in one of the nuclear explosions that devastated the world.

As you struggle to gather your bearings, you notice a blinking panel on the wall, with the words *"Cryo Sleep Malfunction"* displayed in bold letters. It seems that the system has finally detected the error that caused your prolonged slumber and triggered your awakening.

**Shortly after you realize that you are not alone...**]],
    choices = {
        {
            description = "Try to find a weapon. " ..
                highlight('Find melee weapon') .. " " .. highlight_warn("Take 4 damage"),
            callback = function()
                deal_damage(PLAYER_ID, PLAYER_ID, 4, true)
                give_artifact(
                    choose_weighted_by_price(find_artifacts_by_tags({ "HND", "M" })), PLAYER_ID
                )

                return nil
            end
        },
        {
            description = "Gather your strength and attack it!",
            callback = function()
                give_card("MELEE_HIT", PLAYER_ID)
                give_card("MELEE_HIT", PLAYER_ID)
                give_card("MELEE_HIT", PLAYER_ID)

                return nil
            end
        }
    },
    on_enter = function()
        play_music("energetic_orthogonal_expansions")
    end,
    on_end = function()
        actor_set_max_hp(PLAYER_ID, 10)
        actor_set_hp(PLAYER_ID, 10)

        give_card("BLOCK", PLAYER_ID)
        give_card("BLOCK", PLAYER_ID)

        return GAME_STATE_RANDOM
    end
})
