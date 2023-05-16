register_artifact("HOLY_GRAIL", {
    name = "Holy Grail",
    description = "At the start of each turn, heal for 2 HP for each card in your hand.",
    price = 150,
    order = 100, -- Evaluate late so that other draw artifacts have priority.
    callbacks = {
        on_player_turn = function(ctx)
            local num_cards = #get_cards(ctx.owner)
            local heal_amount = num_cards * 2
            heal(ctx.owner, ctx.owner, heal_amount)
            return nil
        end
    }
});