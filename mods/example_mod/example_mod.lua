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