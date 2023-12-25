register_event("BIO_KINGDOM", {
    name = "Bio Kingdom",
    description = [[!!plant_enviroment.png

You finally find a way leading to the outside, and step out of the cryo facility into a world you no longer recognize.

The air is thick with humidity and the sounds of the jungle are overwhelming. Strange, mutated plants tower over you, their vines twisting and tangling around each other in a macabre dance. The colors of the leaves and flowers are sickly, a greenish hue that reminds you of illness rather than life. The ruins of buildings are visible in the distance, swallowed up by the overgrowth. You can hear the chirping and buzzing of insects, but it's mixed with something else - something that sounds almost like whispers or moans. The "jungle" seems to be alive, but not in any way that you would have imagined.]],
    choices = {
        {
            description = "Go...",
            callback = function()
                set_event("MERCHANT")
                return GAME_STATE_EVENT
            end
        }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})
