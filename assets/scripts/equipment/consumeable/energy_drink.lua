local drinks = {
    {
        id = "ENERGY_DRINK",
        name = "ENRGY Drink X91",
        description = "Gain 1 action point.",
        price = 150,
        action_points = 1,
    },
    {
        id = "ENERGY_DRINK_2",
        name = "ENRGY Drink X92",
        description = "Gain 2 action points.",
        price = 250,
        action_points = 2,
    },
    {
        id = "ENERGY_DRINK_3",
        name = "ENRGY Drink X93",
        description = "Gain 3 action points.",
        price = 350,
        action_points = 3,
    },
}

for _, drink in ipairs(drinks) do
    register_card(drink.id, {
        name = l("cards." .. drink.id .. ".name", drink.name),
        description = string.format(
            l("cards." .. drink.id .. ".description","%s\n\n%s"),
            highlight("One-Time"),
            drink.description
        ),
        tags = { "UTIL", "_ACT_0" },
        max_level = 0,
        color = COLOR_ORANGE,
        need_target = false,
        does_consume = true,
        point_cost = 0,
        price = drink.price,
        callbacks = {
            on_cast = function(ctx)
                player_give_action_points(drink.action_points)
                return nil
            end
        },
        test = function ()
            return assert_chain({
                function() return assert_card_present(drink.id) end,
                function() return assert_cast_card(drink.id) end,
                function()
                    if get_fight().current_points ~= 3 + drink.action_points then
                        return "Expected " .. tostring(3 + drink.action_points) .. " points, got " .. get_fight().current_points
                    end
                end,
            })
        end
    })
end