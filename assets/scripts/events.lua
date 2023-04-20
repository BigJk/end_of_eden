register_event(
    "START",
    {
        name = "Waking up...",
        description = [[

```
 █     █░ ▄▄▄       ██ ▄█▀▓█████     █    ██  ██▓███
▓█░ █ ░█░▒████▄     ██▄█▒ ▓█   ▀     ██  ▓██▒▓██░  ██▒
▒█░ █ ░█ ▒██  ▀█▄  ▓███▄░ ▒███      ▓██  ▒██░▓██░ ██▓▒
░█░ █ ░█ ░██▄▄▄▄██ ▓██ █▄ ▒▓█  ▄    ▓▓█  ░██░▒██▄█▓▒ ▒
░░██▒██▓  ▓█   ▓██▒▒██▒ █▄░▒████▒   ▒▒█████▓ ▒██▒ ░  ░ ██▓  ██▓  ██▓
░ ▓░▒ ▒   ▒▒   ▓▒█░▒ ▒▒ ▓▒░░ ▒░ ░   ░▒▓▒ ▒ ▒ ▒▓▒░ ░  ░ ▒▓▒  ▒▓▒  ▒▓▒
  ▒ ░ ░    ▒   ▒▒ ░░ ░▒ ▒░ ░ ░  ░   ░░▒░ ░ ░ ░▒ ░      ░▒   ░▒   ░▒
  ░   ░    ░   ▒   ░ ░░ ░    ░       ░░░ ░ ░ ░░        ░    ░    ░
    ░          ░  ░░  ░      ░  ░      ░                ░    ░    ░
                                                        ░    ░    ░
```

You wake up in a dimly lit room, the faint glow of a red emergency light casting an eerie hue over the surroundings. The air is musty and stale, the metallic scent of the cryo-chamber still lingering in your nostrils. You feel groggy and disoriented, your mind struggling to process what's happening.

As you try to sit up, you notice that your body is stiff and unresponsive. It takes a few moments for your muscles to warm up and regain their strength. Looking around, you see that the walls are made of a dull gray metal, covered in scratches and scuff marks. There's a faint humming sound coming from somewhere, indicating that the facility is still operational.

You try to remember how you ended up here, but your memories are hazy and fragmented. The last thing you recall is a blinding flash of light and a deafening boom. You must have been caught in one of the nuclear explosions that devastated the world.

As you struggle to gather your bearings, you notice a blinking panel on the wall, with the words *"Cryo Sleep Malfunction"* displayed in bold letters. It seems that the system has finally detected the error that caused your prolonged slumber and triggered your awakening.

**Shortly after you realize that you are not alone...**]],
        choices = {
                {
                        description = "Try to escape the facility before it finds you...",
                        callback = function()
                                -- Try to escape
                                if math.random() < 0.5 then
                                        set_event("FIRST_OUTSIDE")
                                        return GAME_STATE_EVENT
                                end

                                -- Let OnEnd handle the state change
                                return nil
                        end
                },
                {
                        description = "Gather your strength and attack it!",
                        callback = function() return nil end
                }
        },
        on_enter = function()
                -- Give the player it's start cards
                give_card("MELEE_HIT", PLAYER_ID)
                give_card("MELEE_HIT", PLAYER_ID)
                give_card("MELEE_HIT", PLAYER_ID)
                give_card("RUPTURE", PLAYER_ID)
                give_card("BLOCK", PLAYER_ID)
                give_artifact(get_random_artifact_type(150), PLAYER_ID)
        end,
        on_end = function()
                return GAME_STATE_RANDOM
        end,
    }
)

register_event(
        "FIRST_OUTSIDE",
        {
                name = "Outside",
                description = [[You finally find a way leading to the outside, a narrow tunnel that winds its way through the thick layer of earth and rock above the facility. The tunnel is cramped and claustrophobic, and you have to crawl on your hands and knees for what feels like hours.

As you near the end of the tunnel, you feel a surge of excitement mixed with fear. What will you find on the other side? Will there be other survivors, or only mutated creatures and plants?

Finally, you emerge into the open air, blinking in the bright sunlight. The landscape that stretches out before you is both familiar and alien, a mix of twisted and mutated plant life, towering rock formations, and ruined remnants of the old world.

You take a deep breath of the fresh air, feeling the warmth of the sun on your face. You know that the journey ahead will be long and perilous, but you're determined to explore this new world and uncover its secrets. **The adventure has only just begun.**]],
                choices = {},
                on_end = function()
                        return GAME_STATE_RANDOM
                end,
        }
)

register_event(
        "CHOICE",
        {
                name = "A choice",
                description = [[Some challenges are behind you. Choose wisely...]],
                choices = {
                        {
                                description = "Meet the Merchant",
                                callback = function()
                                        return GAME_STATE_MERCHANT
                                end
                        },
                        {
                                description = "Remove a Card",
                                callback = function()
                                        set_event("CARD_REMOVE")
                                        return GAME_STATE_EVENT
                                end
                        },
                        {
                                description = "Gain a Card",
                                callback = function()
                                        set_event("CARD_GAIN")
                                        return GAME_STATE_EVENT
                                end
                        },
                        {
                                description = "Rest & Heal",
                                callback = function()
                                        set_event("REST")
                                        return GAME_STATE_EVENT
                                end
                        },
                        {
                                description = "Pass",
                                callback = function()
                                        return GAME_STATE_RANDOM
                                end
                        }
                },
                on_end = function(choice) return nil end,
        }
)

register_event(
        "MERCHANT",
        {
                name = "A strange figure",
                description = [[The merchant is a tall, lanky figure draped in a long, tattered coat made of plant fibers and animal hides. Their face is hidden behind a mask made of twisted roots and vines, giving them an unsettling, almost alien appearance.

Despite their strange appearance, the merchant is a shrewd negotiator and a skilled trader. They carry with them a collection of bizarre and exotic items, including plant-based weapons, animal pelts, and strange, glowing artifacts that seem to pulse with an otherworldly energy.

The merchant is always looking for a good deal, and they're not above haggling with potential customers...]],
                choices = {
                        {
                                description = "Trade",
                                callback = function()
                                        return GAME_STATE_MERCHANT
                                end
                        },
                        {
                                description = "Pass",
                                callback = function()
                                        return GAME_STATE_RANDOM
                                end
                        }
                },
                on_end = function(choice) return nil end,
        }
)