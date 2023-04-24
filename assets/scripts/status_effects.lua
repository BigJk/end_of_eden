--
-- Negative Status Effects
--

register_status_effect("WEAKEN", {
    name = "Weaken",
    description = "Weakens damage for each stack",
    look = "W",
    foreground = "#ed985f",
    state = function()
        return "-" .. tostring(ctx.stacks * 2) .. " damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_damage_calc = function()
            if ctx.source == ctx.owner then
                return ctx.damage - ctx.stacks * 2
            end
            return ctx.damage
        end,
    }
})

register_status_effect("VULNERABLE", {
    name = "Vulnerable",
    description = "Increases received damage for each stack",
    look = "Vur",
    foreground = "#ffba08",
    state = function(ctx)
        return tostring(ctx.stacks * 25) .. "% more damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                return ctx.damage * (1.0 + 0.25 * ctx.stacks)
            end
            return ctx.damage
        end,
    }
})

register_status_effect("BURN", {
    name = "Burning",
    description = "The enemy burns and receives damage.",
    look = "Brn",
    foreground = "#d00000",
    state = function(ctx)
        return tostring(ctx.stacks * 4) .. " damage per turn"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, ctx.owner, ctx.stacks * 2, true)
            return nil
        end,
    }
})

--
-- Positive Status Effects
--

register_status_effect("STRENGTH", {
    name = "Strength",
    description = "Increases damage for each stack",
    look = "Str",
    foreground = "#d00000",
    state = function(ctx)
        return tostring(ctx.stacks) .. " damage"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage + ctx.stacks
            end
            return ctx.damage
        end,
    }
})

register_status_effect("BLOCK", {
    name = "Block",
    description = "Decreases incoming damage for each stack",
    look = "Blk",
    foreground = "#219ebc",
    state = function(ctx)
        return tostring(ctx.stacks) .. " damage reduced"
    end,
    can_stack = true,
    decay = DECAY_ALL,
    rounds = 1,
    order = 100,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                add_status_effect_stacks(ctx.guid, -ctx.damage)
                return ctx.damage - ctx.stacks
            end
            return ctx.damage
        end,
    }
})

register_status_effect("RITUAL", {
    name = "Ritual",
    description = "Gain strength each round",
    look = "Rit",
    foreground = "#bb3e03",
    state = function(ctx)
        return tostring(ctx.stacks * 25) .. "% more damage"
    end,
    can_stack = true,
    decay = DECAY_NONE,
    rounds = 0,
    callbacks = {
        on_turn = function(ctx)
            local guid = give_status_effect("STRENGTH", ctx.owner)
            set_status_effect_stacks(guid, 3 + ctx.stacks)
        end,
    }
})

register_status_effect("FEAR", {
    name = "Fear",
    description = "Can't act.",
    look = "Fear",
    foreground = "#bb3e03",
    state = function(ctx)
        return "Can't act for " .. highlight(tostring(ctx.stacks)) .. " turns"
    end,
    can_stack = true,
    decay = DECAY_ONE,
    rounds = 1,
    callbacks = {
        on_turn = function(ctx)
            return true
        end,
    }
})