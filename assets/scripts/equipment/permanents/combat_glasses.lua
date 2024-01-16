register_artifact("COMBAT_GLASSES", {
    name = "Combat Glasses",
    description = "Whenever you play a " .. highlight("Ranged (R)") .. " card, deal " .. highlight("1 additional damage"),
    tags = { "_ACT_0" },
    price = 100,
    order = 0,
    callbacks = {
        on_damage_calc = function(ctx)
            local card = get_card(ctx.card)
            if card ~= nil then
                if table.contains(card.tags, "R") and ctx.source == ctx.owner and ctx.target ~= ctx.owner then
                    return ctx.damage + 1
                end
            end
            return ctx.damage
        end
    },
    test = function()
        give_card(HAND_WEAPONS[3].id, PLAYER_ID)
        return assert_cast_damage(HAND_WEAPONS[3].id, HAND_WEAPONS[3].base_damage + 1)
    end
});
