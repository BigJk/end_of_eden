print("Example Mod is loaded!")

register_enemy(
        "DOOR",
        {
            name = "Door",
            description = "End me...",
            look = "Door",
            color = "#cccccc",
            initial_hp = 100,
            max_hp = 100,
            callbacks = {
                on_turn = function(ctx)
                    return nil
                end
            }
        }
)

register_event("TEST_EVENT", {
    name = "Test Test",
    description = [[!!mod_test.png

Testing loading a image from mod folder.]],
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