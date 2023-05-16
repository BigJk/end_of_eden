register_artifact("BAG_OF_HOLDING", {
    name = "Bag of Holding",
    description = "Start with a additional card at the beginning of combat.",
    price = 50,
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
