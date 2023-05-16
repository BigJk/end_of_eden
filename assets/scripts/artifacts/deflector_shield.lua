register_artifact("DEFLECTOR_SHIELD", {
    name = "Deflector Shield",
    description = "Gain 8 block at the start of combat.",
    price = 50,
    order = 0,
    callbacks = {
        on_player_turn = function(ctx)
            if ctx.round == 0 then
                give_status_effect("BLOCK", ctx.owner, 8)
            end
            return nil
        end
    }
});
