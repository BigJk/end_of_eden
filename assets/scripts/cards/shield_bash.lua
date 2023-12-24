register_card("SHIELD_BASH", {
    name = "Shield Bash",
    description = "Deal 4 (+2 for each upgrade) damage to the enemy and gain " .. highlight("block") ..
        " status effect equal to the damage dealt.",
    state = function(ctx)
        return "Deal " .. highlight(4 + ctx.level * 2) .. " damage to the enemy and gain " .. highlight("block") ..
                   " status effect equal to the damage dealt."
    end,
    max_level = 1,
    color = "#ff5722",
    need_target = true,
    point_cost = 1,
    price = 40,
    callbacks = {
        on_cast = function(ctx)
            local damage = deal_damage(ctx.caster, ctx.target, 4 + ctx.level * 2)
            give_status_effect("BLOCK", ctx.caster, damage)
            return nil
        end
    },
    test = function()
        dummy = add_actor_by_enemy("DUMMY")
        cards = get_cards(PLAYER_ID)

        -- Check if the card is in the player's hand
        if not cards[1] then
            return "Card not in hand"
        end

        card = get_card_instance(cards[1])
        if card.type_id ~= "SHIELD_BASH" then
            return "Card has wrong type: " .. card.type_id
        end

        cast_card(cards[1], dummy)

        if get_actor(dummy).hp ~= 96 then
            return "Expected 96 health, got " .. get_actor_health(dummy)
        end

        return assert_chain({ assert_status_effect_count(1), assert_status_effect("BLOCK", 4) })
    end
})
