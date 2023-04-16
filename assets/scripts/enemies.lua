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
    "RUST_MITE",
    {
        Name = "Rust Mite",
        Description = "Loves to eat metal.",
        Look = "M",
        Color = "#e6e65a",
        InitialHP = 16,
        MaxHP = 16,
        Callbacks = {
            OnInit = function(ctx)
                give_card("BITE", ctx.guid)
            end,
            OnTurn = function(ctx)
                cast_random(ctx.guid, PLAYER_ID)
            end
        }
    }
)