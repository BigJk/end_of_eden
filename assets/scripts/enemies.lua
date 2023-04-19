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
        "DUMMY",
        {
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
        }
)

register_enemy(
    "RUST_MITE",
    {
        name = "Rust Mite",
        description = "Loves to eat metal.",
        look = "/v\\",
        color = "#e6e65a",
        initial_hp = 16,
        max_hp = 16,
        gold = 10,
        callbacks = {
            on_turn = function(ctx)
                if ctx.round % 4 == 0 then
                    give_status_effect("RITUAL", ctx.guid)
                else
                    deal_damage(ctx.guid, PLAYER_ID, 6)
                end

                return nil
            end
        }
    }
)

register_enemy(
        "CLEAN_BOT",
        {
            name = "Cleaning Bot",
            description = "It never stopped cleaning...",
            look = "BOT",
            color = "#32a891",
            initial_hp = 22,
            max_hp = 22,
            gold = 15,
            callbacks = {
                on_turn = function(ctx)
                    local self = get_actor(ctx.guid)

                    if self.hp < 8 then
                        give_status_effect("BLOCK", ctx.guid, 4)
                    else
                        deal_damage(ctx.guid, PLAYER_ID, 7)
                    end

                    return nil
                end
            }
        }
)