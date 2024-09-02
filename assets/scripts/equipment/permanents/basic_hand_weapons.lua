local hand_warning =
"**Important:** If you already carry a artifact in your hand, you will have to drop it and related cards to pick up the new one."

HAND_WEAPONS = {
    {
        id = "CROWBAR",
        name = "Crowbar",
        image = "red_room.jpg",
        description = "A crowbar. It's a bit rusty, but it should still be useful!",
        base_damage = 2,
        base_cards = 3,
        tags = { "ATK", "M", "T", "HND" },
        event_tags = { "_ACT_0" },
        additional_cards = { "KNOCK_OUT" },
        price = 80
    },
    {
        id = "VIBRO_KNIFE",
        name = "VIBRO Knife",
        description = "A VIBRO knife. Uses ultrasonic vibrations to cut through almost anything.",
        base_damage = 3,
        base_cards = 3,
        tags = { "ATK", "M", "T", "HND" },
        event_tags = { "_ACT_0" },
        additional_cards = { "VIBRO_OVERCLOCK" },
        price = 180
    },
    {
        id = "LZR_PISTOL",
        name = "LZR Pistol",
        description = "A LZR pistol. Fires a concentrated beam of light.",
        base_damage = 4,
        base_cards = 3,
        tags = { "ATK", "R", "T", "HND" },
        event_tags = { "_ACT_1" },
        additional_cards = { "LZR_OVERCHARGE" },
        price = 280
    },
    {
        id = "HAR_II",
        name = "HAR-II",
        description = "A HAR-II. A heavy assault rifle with a high rate of fire.",
        base_damage = 5,
        base_cards = 3,
        tags = { "ATK", "R", "T", "HND" },
        event_tags = { "_ACT_1" },
        additional_cards = { "HAR_BURST", "TARGET_PAINTER" },
        price = 380
    }
}

HAND_WEAPONS_ARTIFACT_IDS = fun.iter(HAND_WEAPONS):map(function(w) return w.id end):totable()

for _, weapon in pairs(HAND_WEAPONS) do
    register_card(weapon.id, {
        name = l("cards." .. weapon.id .. ".name", weapon.name),
        description = l("cards." .. weapon.id .. ".description",
            string.format("Use to deal %s (+1 for each upgrade) damage.", weapon.base_damage)),
        state = function(ctx)
            return string.format(l("cards." .. weapon.id .. ".state", "Use to deal %s damage."),
                highlight(weapon.base_damage + ctx.level * 1))
        end,
        tags = weapon.tags,
        max_level = 3,
        color = COLOR_GRAY,
        need_target = true,
        point_cost = 1,
        price = 0,
        callbacks = {
            on_cast = function(ctx)
                deal_damage_card(ctx.caster, ctx.guid, ctx.target, weapon.base_damage + ctx.level * 3)
                return nil
            end
        },
        test = function()
            local dummy = add_actor_by_enemy("DUMMY")
            local cards = get_cards(PLAYER_ID)

            -- Check if the card is in the player's hand
            if not cards[1] then
                return "Card not in hand"
            end

            local card = get_card_instance(cards[1])
            if card.type_id ~= weapon.id then
                return "Card has wrong type: " .. card.type_id
            end

            cast_card(cards[1], dummy)

            if get_actor(dummy).hp ~= 100 - weapon.base_damage then
                return "Expected " .. tostring(100 - weapon.base_damage) .. " health, got " .. get_actor(dummy).hp
            end

            return nil
        end
    })

    register_artifact(weapon.id, {
        name = weapon.name,
        description = weapon.description .. " Can be used in your hand.",
        tags = weapon.tags,
        price = weapon.price,
        order = 0,
        callbacks = {
            on_pick_up = function(ctx)
                clear_artifacts_by_tag("HND", { ctx.guid })
                clear_cards_by_tag("HND")

                -- add basic cards
                for i = 1, weapon.base_cards do
                    give_card(weapon.id, PLAYER_ID)
                end

                -- add additional cards
                if weapon.additional_cards then
                    for _, card in pairs(weapon.additional_cards) do
                        give_card(card, PLAYER_ID)
                    end
                end

                return nil
            end
        }
    })

    add_found_artifact_event(weapon.id, weapon.image, string.format("%s\n\n%s", weapon.description, hand_warning),
        registered.card[weapon.id].description, weapon.event_tags)
end

---hand_weapon_event returns a random hand weapon event weighted by price.
---@return string
function hand_weapon_event()
    local ids = fun.iter(HAND_WEAPONS):map(function(w) return w.id end):totable()
    local prices = fun.iter(HAND_WEAPONS):map(function(w) return 500 - w.price end):totable()
    return choose_weighted(ids, prices)
end
