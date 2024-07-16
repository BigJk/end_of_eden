register_artifact("COMBAT_GLOVES", {
    name = "Combat Gloves",
    description = "Whenever you play a " .. highlight("Melee (M)") .. " card, deal " .. highlight("1 additional damage"),
    tags = { "_ACT_0" },
    price = 100,
    order = 0,
    callbacks = {
        on_damage_calc = function (ctx)
            local card = get_card(ctx.card)
            if card ~= nil then
                if table.contains(card.tags, "M") and ctx.source == ctx.owner and ctx.target ~= ctx.owner then
                    return ctx.damage + 1
                end
            end
            return ctx.damage
        end
    },
    test = function ()
        give_card("MELEE_HIT", PLAYER_ID)
        return assert_cast_damage("MELEE_HIT", 2)
    end
});
