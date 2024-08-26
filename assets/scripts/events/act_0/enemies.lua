register_event("RUST_MITE", {
    name = "Tasty metals...",
    description = [[!!rust_mite.jpg

You are walking through the facility hoping to find a way out. After a few turns you hear a strange noise. You look around and come across a strange being.
It seems to be eating the metal from the walls. It looks at you and after a few seconds it rushes towards you.

**It seems to be hostile!**
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("RUST_MITE")
                if math.random() < 0.25 then
                    add_actor_by_enemy("RUST_MITE")
                end
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("CLEAN_BOT", {
    name = "Corpse. Clean. Engage.",
    description = [[!!clean_bot.jpg

While exploring the facility you hear a strange noise. Suddenly a strange robot appears from one of the corridors.
It seems to be cleaning up the area, but it's not working properly anymore and you can see small sparks coming out of it.
It looks at you and says "Corpse. Clean. Engage.".

**You're not sure what it means, but it doesn't seem to be friendly!**
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("CLEAN_BOT")
                if math.random() < 0.25 then
                    add_actor_by_enemy("CLEAN_BOT")
                end
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("CYBER_SPIDER", {
    name = "What is this thing at the ceiling?",
    description = [[!!cyber_spider.jpg

You come around a corner and see a strange creature hanging from the ceiling. It looks like a spider, but it's made out of metal.
It seems to be waiting for its prey to come closer and there is no way around it.
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("CYBER_SPIDER")
                if math.random() < 0.25 then
                    add_actor_by_enemy("CYBER_SPIDER")
                end
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("LASER_DRONE", {
    name = "A menacing drone appears...",
    description =
    [[As you explore the facility, you hear a high-pitched whirring sound. A drone equipped with a powerful laser cannon appears in front of you.

**It looks ready to attack!**
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("LASER_DRONE")
                if math.random() < 0.10 then
                    add_actor_by_enemy("LASER_DRONE")
                end
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("PLASMA_GOLEM", {
    name = "A glowing figure emerges...",
    description =
    [[As you delve deeper into the facility, you notice a bright glow emanating from a nearby chamber. A massive golem made of pure plasma energy steps into view.

**It looks ready to unleash its power!**
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("PLASMA_GOLEM")
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("CYBER_SLIME", {
    name = "A strange cybernetic slime appears...",
    description =
    [[As you explore the facility, you come across a strange cybernetic slime. It seems to be pulsating with energy and looks hostile.

**Prepare for a fight!**
    ]],
    tags = { "_ACT_0_FIGHT" },
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("CYBER_SLIME")
                return GAME_STATE_FIGHT
            end
        }
    }
})
