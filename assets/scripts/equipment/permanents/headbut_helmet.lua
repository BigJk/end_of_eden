register_artifact("HEADBUT_HELMET", {
    name = "Headbut Helmet",
    description = "Gain 1 Knock Out card.",
    tags = { "_ACT_0" },
    price = 200,
    order = 0,
    callbacks = {
        on_pick_up = function(ctx)
            give_card("KNOCK_OUT", PLAYER_ID)
            return nil
        end
    },
    test = function()
        local cards = get_cards(PLAYER_ID)
        for _, card_guid in ipairs(cards) do
            local card = get_card_instance(card_guid)
            if card.type_id == "KNOCK_OUT" then
                return nil
            end
        end
        return "Expected to find KNOCK_OUT card, but did not."
    end
})
