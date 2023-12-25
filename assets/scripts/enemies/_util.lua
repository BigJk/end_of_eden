---cast_random casts a random card onto a target from the cards that the actor specified by guid owns.
---@param guid guid
---@param target guid
function cast_random(guid, target)
    local cards = get_cards(guid)
    if #cards == 0 then
        print("can't cast_random with zero cards available!")
    else
        cast_card(cards[math.random(#cards)], target)
    end
end

register_enemy("DUMMY", {
    name = "Dummy",
    description = "End me...",
    look = "DUM",
    color = "#deeb6a",
    initial_hp = 100,
    max_hp = 100,
    callbacks = {
        on_turn = function(ctx)
            return nil
        end
    }
})
