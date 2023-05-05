--
-- Base Events
--
register_event("MERCHANT", {
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
        }, {
            description = "Pass",
            callback = function()
                return GAME_STATE_RANDOM
            end
        }
    },
    on_end = function(ctx)
        return nil
    end
})

register_event("START", {
    name = "Waking up...",
    description = [[!!cryo_start.png

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
        }, {
            description = "Gather your strength and attack it!",
            callback = function()
                return nil
            end
        }
    },
    on_enter = function()
        play_music("energetic_orthogonal_expansions")

        -- Give the player it's start cards
        give_card("MELEE_HIT", PLAYER_ID)
        give_card("MELEE_HIT", PLAYER_ID)
        give_card("MELEE_HIT", PLAYER_ID)
        give_card("MELEE_HIT", PLAYER_ID)
        give_card("MELEE_HIT", PLAYER_ID)

        give_card("RUPTURE", PLAYER_ID)

        give_card("BLOCK", PLAYER_ID)
        give_card("BLOCK", PLAYER_ID)
        give_card("BLOCK", PLAYER_ID)

        give_artifact(get_random_artifact_type(150), PLAYER_ID)
    end,
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

--
-- Stage 1 Entrance
--

register_event("FIRST_OUTSIDE", {
    name = "Outside",
    description = [[You finally find a way leading to the outside, a narrow tunnel that winds its way through the thick layer of earth and rock above the facility. The tunnel is cramped and claustrophobic, and you have to crawl on your hands and knees for what feels like hours.

As you near the end of the tunnel, you feel a surge of excitement mixed with fear. What will you find on the other side? Will there be other survivors, or only mutated creatures and plants?

Finally, you emerge into the open air, blinking in the bright sunlight. The landscape that stretches out before you is both familiar and alien, a mix of twisted and mutated plant life, towering rock formations, and ruined remnants of the old world.

You take a deep breath of the fresh air, feeling the warmth of the sun on your face. You know that the journey ahead will be long and perilous, but you're determined to explore this new world and uncover its secrets. **The adventure has only just begun.**]],
    choices = {
        {
            description = "Go...",
            callback = function()
                return nil
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_event("THE_WASTELAND", {
    name = "The Wasteland",
    description = [[!!dark_city1.png

You finally find a way leading to the outside, and with a deep breath, you step out into the unforgiving wasteland.

The scorching sun beats down on you as the sand whips against your skin, a reminder of the horrors that have befallen the world. In the distance, the remains of once-great cities jut up from the ground like jagged teeth, now nothing more than crumbling ruins. The air is thick with the acrid smell of decay and the oppressive silence is only broken by the occasional howl of some mutated creature. As you take your first steps into this new world, you realize that survival will not be easy, and that the journey ahead will be fraught with danger at every turn...]],
    choices = {
        {
            description = "Go...",
            callback = function()
                return nil
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_event("THE_CORE", {
    name = "The Wasteland",
    description = [[!!underground1.png
    
You finally find a way you thought would lead to the outside, only to discover that you're still inside the massive facility known as *"The Core."*

As you step out of the cryo facility, the eerie silence is broken by the sound of metal scraping against metal and distant whirring of malfunctioning machinery. The flickering lights and sparks from faulty wires cast a sickly glow on the cold metal walls. You realize that this place is not as deserted as you initially thought, and the unsettling feeling in your gut only grows stronger as you make your way through the dimly lit corridors, surrounded by the echoes of your own footsteps and the sound of flickering computer screens.]],
    choices = {
        {
            description = "Go...",
            callback = function()
                return nil
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_event("BIO_KINGDOM", {
    name = "Bio Kingdom",
    description = [[!!plant_enviroment.png

You finally find a way leading to the outside, and step out of the cryo facility into a world you no longer recognize.

The air is thick with humidity and the sounds of the jungle are overwhelming. Strange, mutated plants tower over you, their vines twisting and tangling around each other in a macabre dance. The colors of the leaves and flowers are sickly, a greenish hue that reminds you of illness rather than life. The ruins of buildings are visible in the distance, swallowed up by the overgrowth. You can hear the chirping and buzzing of insects, but it's mixed with something else - something that sounds almost like whispers or moans. The "jungle" seems to be alive, but not in any way that you would have imagined.]],
    choices = {
        {
            description = "Go...",
            callback = function()
                return nil
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_event("", {
    name = "Talking Being",
    description = [[!!alien2.png

Suddenly, a massive vine with a gaping, tooth-filled maw emerges from the shadows. It towers over you, its presence imposing and otherworldly.

*"Hello, little one,"* the creature speaks in a deep, rumbling voice. *"I have been watching you. I see potential in you. I offer you a gift, something that will aid you on your journey."*

You take a step back, unsure if you can trust this strange being.

*"My blood,"* the creature says. *"It is not like any substance you have encountered before. It will grant you extraordinary abilities. But it demands a price. Some of your blood, in exchange for this gift."*

The creature assures you that there are dangers to wielding such power and that it will change you in ways you cannot yet imagine. But the offer is tempting. Will you accept and risk the unknown, or do you refuse and potentially miss out on a powerful ally?

**The decision is yours...**]],
    choices = {
        {
            description_fn = function()
                return "Offer blood... " .. text_italic("(deals " .. highlight(get_player().hp * 0.2) .. " damage)")
            end,
            callback = function(ctx)
                actor_add_hp(PLAYER_ID, -get_player().hp * 0.2)
                give_card("VINE_VOLLEY", PLAYER_ID)
                give_card("VINE_VOLLEY", PLAYER_ID)
                give_card("VINE_VOLLEY", PLAYER_ID)
                return nil
            end
        }, {
            description = "Leave...",
            callback = function()
                return nil
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})
