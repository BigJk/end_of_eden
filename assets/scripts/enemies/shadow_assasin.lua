register_enemy("SHADOW_ASSASSIN", {
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
            :filter(function(val)
                return val.type_id == "BLEED"
            end):totable()

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
                :filter(function(
                    val)
                    return val.type_id == "BLEED"
                end):totable()

            if #bleeds > 0 then
                -- If bleeding do more damage
                deal_damage(ctx.guid, PLAYER_ID, 10)
            elseif ctx.round % 3 == 0 then
                -- Try to bleed every 2 rounds with 3 dmg
                if deal_damage(ctx.guid, PLAYER_ID, 3) > 0 then
                    give_status_effect("BLEED", PLAYER_ID, 2)
                end
            else
                -- Just hit with 5 damage
                deal_damage(ctx.guid, PLAYER_ID, 5)
            end

            return nil
        end
    }
})

register_status_effect("BLEED", {
    name = "Bleed",
    description = "Losing some red sauce.",
    look = "Bld",
    foreground = "#ff0000",
    state = function(ctx)
        return nil
    end,
    can_stack = false,
    decay = DECAY_ONE,
    rounds = 2,
    callbacks = {
        on_turn = function(ctx)
            return nil
        end
    }
})
