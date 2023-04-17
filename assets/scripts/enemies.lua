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
    "RUST_MITE",
    {
        name = "Rust Mite",
        description = "Loves to eat metal.",
        look = "M",
        color = "#e6e65a",
        initial_hp = 16,
        max_hp = 16,
        callbacks = {
            on_turn = function(ctx)
                if ctx.round % 3 == 0 then
                    give_status_effect("RITUAL", ctx.guid)
                else
                    deal_damage(ctx.guid, PLAYER_ID, 6)
                end
            end
        }
    }
)