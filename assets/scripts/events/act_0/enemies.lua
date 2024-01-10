register_event("RUST_MITE", {
    name = "Tasty metals...",
    description = [[
        You come across a strange being. It seems to be eating the metal from the walls. It looks at you and after a few seconds it rushes towards you. It seems to be hostile.
    ]],
    tags = {"ACT_0"},
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("RUST_MITE")
                return GAME_STATE_FIGHT
            end
        }
    }
})

register_event("CLEAN_BOT", {
    name = "Corpse. Clean. Engage.",
    description = [[
        You come across a strange robot. It seems to be cleaning up the area. It looks at you and says "Corpse. Clean. Engage.". You're not sure what it means, but it doesn't seem to be friendly.
    ]],
    tags = {"ACT_0"},
    choices = {
        {
            description = "Fight!",
            callback = function()
                add_actor_by_enemy("CLEAN_BOT")
                return GAME_STATE_FIGHT
            end
        }
    }
})
