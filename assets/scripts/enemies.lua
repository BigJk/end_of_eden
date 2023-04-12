-- cast_random
--
-- casts a random card onto a target from the cards that the actor specified by guid owns.
function cast_random(guid, target)
    local cards = get_cards(guid)
    if #cards == 0 then
        debug_log("can't cast_random with zero cards available!")
    else
        cast_card(cards[math.random(#cards)], target)
    end
end

register_enemy(
    "DOOR",
    {
        Name = "Door",
        Description = "It's in your way...",
        Look = "D",
        Color = "#cccccc",
        InitialHP = 10,
        MaxHP = 10,
        Callbacks = { }
    }
)

register_enemy(
    "MUTATED_HAMSTER",
    {
        Name = "Mutated Hamster",
        Description = "Small but furious...",
        Look = "H",
        Color = "#e6e65a",
        InitialHP = 4,
        MaxHP = 4,
        Callbacks = {
            OnInit = function(type, guid)
                give_card("BITE", guid)
            end,
            OnTurn = function(type, guid)
                cast_random(guid, PLAYER_ID)
            end
        }
    }
)