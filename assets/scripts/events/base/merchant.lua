register_event("MERCHANT", {
    name = "A strange figure",
    description =
    [[!!merchant.jpg
    
The merchant is a tall, lanky figure draped in a long, tattered coat made of plant fibers and animal hides. Their face is hidden behind a mask made of twisted roots and vines, giving them an unsettling, almost alien appearance.

Despite their strange appearance, the merchant is a shrewd negotiator and a skilled trader. They carry with them a collection of bizarre and exotic items, including plant-based weapons, animal pelts, and strange, glowing artifacts that seem to pulse with an otherworldly energy.

The merchant is always looking for a good deal, and they're not above haggling with potential customers...]],
    tags = {"_ACT_0", "_ACT_1", "_ACT_2", "_ACT_3"},
    choices = {
        {
            description = "Trade",
            callback = function()
                return GAME_STATE_MERCHANT
            end
        }, {
        description = "Pass",
        callback = function()
            return GAME_STATE_RANDOM
        end
    }
    },
    on_end = function(ctx)
        return nil
    end
})
