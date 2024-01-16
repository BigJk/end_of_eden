register_artifact("PORTABLE_BUFFER", {
    name = "PRTBL Buffer",
    description = "Start each turn with 1 " .. highlight("Block"),
    tags = { "_ACT_0" },
    price = 100,
    order = 0,
    callbacks = {
        on_player_turn = function(ctx)
            give_status_effect("BLOCK", PLAYER_ID, 1)
            return nil
        end
    }
});
