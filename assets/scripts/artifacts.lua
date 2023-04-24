register_artifact(
    "GIGANTIC_STRENGTH",
    {
        name = "Stone Of Gigantic Strength",
        description = "Double all damage dealt.",
        price = 250,
        order = 0,
        callbacks = {
            on_damage_calc = function(ctx)
                if ctx.target == ctx.owner then
                    return ctx.damage * 2
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "REPULSION_STONE",
    {
        name = "Repulsion Stone",
        description = "For each damage taken heal for 2",
        price = 100,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.target == ctx.owner then
                    heal(ctx.owner, 2)
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "RADIANT_SEED",
    {
        name = "Radiant Seed",
        description = "A small glowing seed.",
        price = 140,
        order = 0,
        callbacks = {
            on_pick_up = function(ctx)
                give_card("RADIANT_SEED", ctx.owner)
                return nil
            end,
        }
    }
);

register_artifact(
    "JUICY_FRUIT",
    {
        name = "Juicy Fruit",
        description = "Tastes good and boosts your HP.",
        price = 80,
        order = 0,
        callbacks = {
            on_pick_up = function(ctx)
                actor_add_max_hp(ctx.owner, 10)
                return nil
            end,
        }
    }
);

register_artifact(
    "DEFLECTOR_SHIELD",
    {
        name = "Deflector Shield",
        description = "Gain 8 block at the start of combat.",
        price = 50,
        order = 0,
        callbacks = {
            on_player_turn = function(ctx)
                if ctx.round == 0 then
                    give_status_effect("BLOCK", ctx.owner, 8)
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "SHORT_RADIANCE",
    {
        name = "Short Radiance",
        description = "Apply 1 vulnerable at the start of combat.",
        price = 50,
        order = 0,
        callbacks = {
            on_player_turn = function(ctx)
                if ctx.round == 0 then
                    each(function(val)
                        give_status_effect("VULNERABLE", val)
                    end, pairs(get_opponent_guids(ctx.owner)))
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "BAG_OF_HOLDING",
    {
        name = "Bag of Holding",
        description = "Start with a additional card at the beginning of combat.",
        price = 50,
        order = 0,
        callbacks = {
            on_player_turn = function(ctx)
                if ctx.owner == PLAYER_ID and ctx.round == 0 then
                    player_draw_card(1)
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "SPIKED_PLANT",
    {
        name = "Spiked Plant",
        description = "Deal 2 damage back to enemy attacks.",
        price = 50,
        order = 0,
        callbacks = {
            on_damage = function(ctx)
                if ctx.source ~= ctx.owner and ctx.owner == ctx.target then
                    deal_damage(ctx.owner, ctx.source, 2)
                end
                return nil
            end,
        }
    }
);

register_artifact(
    "GOLD_CONVERTER",
    {
        name = "Gold Converter",
        description = "Gain 10 extra gold for each killed enemy.",
        price = 50,
        order = 0,
        callbacks = {
            on_actor_die = function(ctx)
                if ctx.owner == PLAYER_ID and ctx.owner == ctx.source then
                    give_player_gold(10)
                end
                return nil
            end,
        }
    }
);