register_artifact("SPEED_ENHANCER", {
    name = "Speed Enhancer",
    description = "Start with a additional card at the beginning of combat.",
    tags = { "_ACT_0" },
    price = 100,
    order = 0,
    callbacks = {
        on_player_turn = function(ctx)
            if ctx.owner == PLAYER_ID and ctx.round == 0 then
                player_draw_card(1)
            end
            return nil
        end
    }
});
