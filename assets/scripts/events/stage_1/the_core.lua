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