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
        initial_hp = 22,
        max_hp = 22,
        gold = 10,
        intend = function(ctx)
            if ctx.round % 4 == 0 then
                return "Gather strength"
            end

            return "Deal " .. highlight(6) .. " damage"
        end,
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
        look = [[ \_/
(* *)
 )#(]],
        color = "#32a891",
        initial_hp = 25,
        max_hp = 25,
        gold = 15,
        intend = function(ctx)
            local self = get_actor(ctx.guid)
            if self.hp <= 8 then
                return "Block " .. highlight(4)
            end

            return "Deal " .. highlight(7) .. " damage"
        end,
        callbacks = {
            on_player_turn = function(ctx)
                local self = get_actor(ctx.guid)

                if self.hp <= 8 then
                    give_status_effect("BLOCK", ctx.guid, 4)
                end
            end,
            on_turn = function(ctx)
                local self = get_actor(ctx.guid)

                if self.hp > 8 then
                    deal_damage(ctx.guid, PLAYER_ID, 7)
                end

                return nil
            end
        }
    }
)

register_enemy(
    "SHADOW_ASSASSIN",
    {
        name = "Shadow Assassin",
        description = "A master of stealth and deception.",
        look = "???",
        color = "#6c5b7b",
        initial_hp = 20,
        max_hp = 20,
        gold = 30,
        intend = function(ctx)
            local bleeds = fun.iter(pairs(get_actor_status_effects(PLAYER_ID)))
                              :map(get_status_effect_instance)
                              :filter(function(val) return val.type_id == "BLEED" end)
                              :totable()

            if #bleeds > 0 then
                return "Deal " .. highlight(10) .. " damage"
            elseif ctx.round % 3 == 0 then
                return "Inflict bleed"
            else
                return "Deal " .. highlight(5) .. " damage"
            end

            return nil
        end,
        callbacks = {
            on_turn = function(ctx)
                -- Count bleed stacks
                local bleeds = fun.iter(pairs(get_actor_status_effects(PLAYER_ID)))
                                  :map(get_status_effect_instance)
                                  :filter(function(val) return val.type_id == "BLEED" end)
                                  :totable()

                if #bleeds > 0 then -- If bleeding do more damage
                    deal_damage(ctx.guid, PLAYER_ID, 10)
                elseif ctx.round % 3 == 0 then -- Try to bleed every 2 rounds with 3 dmg
                    if deal_damage(ctx.guid, PLAYER_ID, 3) > 0 then
                        give_status_effect("BLEED", PLAYER_ID, 2)
                    end
                else -- Just hit with 5 damage
                    deal_damage(ctx.guid, PLAYER_ID, 5)
                end

                return nil
            end
        }
    }
)