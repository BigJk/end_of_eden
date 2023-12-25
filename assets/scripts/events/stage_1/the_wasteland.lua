register_event("THE_WASTELAND", {
    name = "The Wasteland",
    description = [[!!dark_city1.png

You finally find a way leading to the outside, and with a deep breath, you step out into the unforgiving wasteland.

The scorching sun beats down on you as the sand whips against your skin, a reminder of the horrors that have befallen the world. In the distance, the remains of once-great cities jut up from the ground like jagged teeth, now nothing more than crumbling ruins. The air is thick with the acrid smell of decay and the oppressive silence is only broken by the occasional howl of some mutated creature. As you take your first steps into this new world, you realize that survival will not be easy, and that the journey ahead will be fraught with danger at every turn...]],
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
