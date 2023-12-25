register_event("TALKING_BEING", {
    name = "Talking Being",
    description = [[!!alien2.png

Suddenly, a massive vine with a gaping, tooth-filled maw emerges from the shadows. It towers over you, its presence imposing and otherworldly.

*"Hello, little one,"* the creature speaks in a deep, rumbling voice. *"I have been watching you. I see potential in you. I offer you a gift, something that will aid you on your journey."*

You take a step back, unsure if you can trust this strange being.

*"My blood,"* the creature says. *"It is not like any substance you have encountered before. It will grant you extraordinary abilities. But it demands a price. Some of your blood, in exchange for this gift."*

The creature assures you that there are dangers to wielding such power and that it will change you in ways you cannot yet imagine. But the offer is tempting. Will you accept and risk the unknown, or do you refuse and potentially miss out on a powerful ally?

**The decision is yours...**]],
    choices = {
        {
            description_fn = function()
                return "Offer blood... " .. text_italic("(deals " .. highlight(get_player().hp * 0.2) .. " damage)")
            end,
            callback = function(ctx)
                actor_add_hp(PLAYER_ID, -get_player().hp * 0.2)
                give_card("VINE_VOLLEY", PLAYER_ID)
                give_card("VINE_VOLLEY", PLAYER_ID)
                give_card("VINE_VOLLEY", PLAYER_ID)
                return nil
            end
        }, {
        description = "Leave...",
        callback = function()
            return nil
        end
    }
    },
    on_end = function()
        return GAME_STATE_RANDOM
    end
})

register_card("VINE_VOLLEY", {
    name = "Vine Volley",
    description = "Deal " .. highlight("3x" .. tostring(3)) .. " damage.",
    state = function(ctx)
        return nil
    end,
    max_level = 0,
    color = "#588157",
    need_target = true,
    point_cost = 1,
    price = 100,
    callbacks = {
        on_cast = function(ctx)
            deal_damage(ctx.caster, ctx.target, 3)
            deal_damage(ctx.caster, ctx.target, 3)
            deal_damage(ctx.caster, ctx.target, 3)
            return nil
        end
    }
})
