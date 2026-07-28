register_event("MERCHANT", {
    name = l("events.MERCHANT.name", "A strange figure"),
    description =
    l("events.MERCHANT.description", [[!!merchant.jpg

The merchant is a tall, lanky figure draped in a long, tattered coat made of plant fibers and animal hides. Their face is hidden behind a mask made of twisted roots and vines, giving them an unsettling, almost alien appearance.

Despite their strange appearance, the merchant is a shrewd negotiator and a skilled trader. They carry with them a collection of bizarre and exotic items, including plant-based weapons, animal pelts, and strange, glowing artifacts that seem to pulse with an otherworldly energy.

The merchant is always looking for a good deal, and they're not above haggling with potential customers...]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.MERCHANT.choices.0.description", "Trade"),
            callback = function()
                return GAME_STATE_MERCHANT
            end
        }, {
        description = l("events.MERCHANT.choices.1.description", "Pass"),
        callback = function()
            return GAME_STATE_RANDOM
        end
    }
    },
    on_end = function(ctx)
        return nil
    end
})


register_event("RANDOM_ARTIFACT_ACT_0", {
    name = l("events.RANDOM_ARTIFACT_ACT_0.name", "Random Artifact"),
    description = l("events.RANDOM_ARTIFACT_ACT_0.description", [[!!artifact_chest.jpg

You found a chest with a strange symbol on it. The chest is protected by a strange barrier. You can either open it and take some damage or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.RANDOM_ARTIFACT_ACT_0.choices.0.prefix", "Random Artifact ") ..
                highlight_success(l("ui.gain_1_artifact", "Gain 1 Artifact")) .. " " .. highlight_warn(l("ui.take_2_damage", "Take 2 damage")),
            callback = function()
                local possible = find_artifacts_by_tags({ "_ACT_0" })
                local choosen = choose_weighted_by_price(possible)
                if choosen then
                    give_artifact(choosen, PLAYER_ID)
                    deal_damage(PLAYER_ID, PLAYER_ID, 2, true)
                end
                return nil
            end
        },
        {
            description = l("events.RANDOM_ARTIFACT_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})

register_event("RANDOM_CONSUMEABLE_ACT_0", {
    name = l("events.RANDOM_CONSUMEABLE_ACT_0.name", "Random Consumeable"),
    description = l("events.RANDOM_CONSUMEABLE_ACT_0.description", [[!!artifact_chest.jpg

You found a chest with a strange symbol on it. The chest is protected by a strange barrier. You can either open it and take some damage or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.RANDOM_CONSUMEABLE_ACT_0.choices.0.prefix", "Random Artifact ") ..
                highlight_success(l("ui.gain_1_consumeable", "Gain 1 Consumeable")) .. " " .. highlight_warn(l("ui.take_2_damage", "Take 2 damage")),
            callback = function()
                local possible = fun.iter(find_cards_by_tags({ "_ACT_0" }))
                    :filter(function(card)
                        return card.does_consume
                    end):totable()
                local choosen = choose_weighted_by_price(possible)
                if choosen then
                    give_card(choosen, PLAYER_ID)
                    deal_damage(PLAYER_ID, PLAYER_ID, 2, true)
                end
                return nil
            end
        },
        {
            description = l("events.RANDOM_CONSUMEABLE_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})

register_event("GAIN_GOLD_ACT_0", {
    name = l("events.GAIN_GOLD_ACT_0.name", "Old Gold Cache"),
    description = l("events.GAIN_GOLD_ACT_0.description", [[
You find an old chest filled with gold. You can either take it or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.GAIN_GOLD_ACT_0.choices.0.prefix", "Take it! ") .. highlight_success(l("ui.gain_20_gold", "Gain 20 Gold")),
            callback = function()
                give_player_gold(20)
                return nil
            end
        },
        {
            description = l("events.GAIN_GOLD_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})


register_event("GOLD_TO_HP_ACT_0", {
    name = l("events.GOLD_TO_HP_ACT_0.name", "Old Vending Machine"),
    description = l("events.GOLD_TO_HP_ACT_0.description", [[
You find an old vending machine, it seems to be still working. You can either pay 20 Gold to get 5 HP or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.GOLD_TO_HP_ACT_0.choices.0.prefix", "Pay ") .. highlight_warn(l("ui.pay_20_gold", "20 Gold")) .. " " .. highlight_success(l("ui.gain_5_hp", "Gain 5 HP")),
            callback = function()
                if get_actor(PLAYER_ID).gold < 20 then
                    return nil
                end
                give_player_gold(-20)
                heal(PLAYER_ID, PLAYER_ID, 5)
                return nil
            end
        },
        {
            description = l("events.GOLD_TO_HP_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})

register_event("MAX_LIFE_ACT_0", {
    name = l("events.MAX_LIFE_ACT_0.name", "Symbiotic Parasite"),
    description = l("events.MAX_LIFE_ACT_0.description", [[!!symbiotic_parasite.jpg

You find a strange creature, it seems to be a symbiotic parasite. It offers to increase your max HP by 5. You can either accept or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.MAX_LIFE_ACT_0.choices.0.prefix", "Accept it! ") .. highlight_success(l("ui.gain_5_max_hp", "Gain 5 Max HP")),
            callback = function()
                actor_add_max_hp(PLAYER_ID, 5)
                return nil
            end
        },
        {
            description = l("events.MAX_LIFE_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})

register_event("GAMBLE_1_ACT_0", {
    name = l("events.GAMBLE_1_ACT_0.name", "Electro Barrier"),
    description = l("events.GAMBLE_1_ACT_0.description", [[!!electro_barrier.jpg

You find a room with a strange device in the middle. It seems to be some kind of electro barrier protecting a storage container. You can either try to disable the barrier or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.GAMBLE_1_ACT_0.choices.0.prefix", "50% ") ..
                highlight_success(l("ui.gain_artifact_and_consumeable", "Gain Artifact & Consumeable")) .. " 50% " .. highlight_warn(l("ui.take_2_damage", "Take 2 damage")),
            callback = function()
                local possible_artifacts = find_artifacts_by_tags({ "_ACT_0" })
                local possible_consumeables = fun.iter(find_cards_by_tags({ "_ACT_0" }))
                    :filter(function(card)
                        return card.does_consume
                    end):totable()
                if random() < 0.5 then
                    local choosen = choose_weighted_by_price(possible_artifacts)
                    if choosen then
                        give_artifact(choosen, PLAYER_ID)
                    end
                    choosen = choose_weighted_by_price(possible_consumeables)
                    if choosen then
                        give_card(choosen, PLAYER_ID)
                    end
                else
                    deal_damage(PLAYER_ID, PLAYER_ID, 2, true)
                end
                return nil
            end
        },
        {
            description = l("events.GAMBLE_1_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})

register_event("UPRAGDE_CARD_ACT_0", {
    name = l("events.UPRAGDE_CARD_ACT_0.name", "Upgrade Station"),
    description = l("events.UPRAGDE_CARD_ACT_0.description", [[!!upgrade_station.jpg

You find a old automatic workstation. You are able to get it working again. You can either upgrade a random card or leave.
    ]]),
    tags = { "_ACT_0" },
    choices = {
        {
            description = l("events.UPRAGDE_CARD_ACT_0.choices.0.prefix", "Upgrade a card ") ..
                highlight_success(l("ui.upgrade_a_card", "Upgrade a card")) .. " " .. highlight_warn(l("ui.take_2_damage", "Take 2 damage")),
            callback = function()
                local cards = fun.iter(get_cards(PLAYER_ID))
                    :filter(function(guid)
                        local type = get_card(guid)
                        local instance = get_card_instance(guid)

                        return instance.level < type.max_level
                    end)
                    :totable()

                if #cards == 0 then
                    return nil
                end

                local choosen = cards[random_int(0, #cards)]
                upgrade_card(choosen)
                deal_damage(PLAYER_ID, PLAYER_ID, 2, true)

                return nil
            end
        },
        {
            description = l("events.UPRAGDE_CARD_ACT_0.choices.1.description", "Leave!"),
            callback = function()
                return nil
            end
        }
    }
})