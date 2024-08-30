delete_base_game("event")

register_enemy("TUTORIAL_DUMMY_1", {
    name = l("enemies.TUTORIAL_DUMMY_1.name", "Dummy"),
    description = l("enemies.TUTORIAL_DUMMY_1.description", "A dummy enemy for the tutorial"),
    look = "D",
    color = "#e6e65a",
    initial_hp = 4,
    max_hp = 4,
    gold = 0,
    intend = function(ctx)
        return "Deal " .. highlight(simulate_deal_damage(ctx.guid, PLAYER_ID, 1)) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            deal_damage(ctx.guid, PLAYER_ID, 1)
            return nil
        end
    }
})

register_enemy("TUTORIAL_DUMMY_2", {
    name = l("enemies.TUTORIAL_DUMMY_2.name", "Dummy"),
    description = l("enemies.TUTORIAL_DUMMY_2.description", "A dummy enemy for the tutorial"),
    look = "D",
    color = "#e6e65a",
    initial_hp = 3,
    max_hp = 3,
    gold = 0,
    intend = function(ctx)
        return "Apply " .. highlight("Weakness")
    end,
    callbacks = {
        on_turn = function(ctx)
            give_status_effect("WEAKNESS", PLAYER_ID)
            return nil
        end
    }
})

register_status_effect("WEAKNESS", {
    name = "Weakness",
    description = "Decreases damage dealt by 1",
    look = "W",
    foreground = COLOR_RED,
    state = function(ctx)
        return "Deals " .. highlight(1) .. " less damage"
    end,
    rounds = 2,
    decay = DECAY_ONE,
    can_stack = false,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.source == ctx.owner then
                return ctx.damage - 1
            end
            return ctx.damage
        end
    }
})

register_card("MELEE_HIT", {
    name = l("cards.MELEE_HIT.name", "Melee Hit"),
    description = l("cards.MELEE_HIT.description", "Use your bare hands to deal 2 (+1 for each upgrade) damage."),
    state = function(ctx)
        return string.format(l("cards.MELEE_HIT.state", "Use your bare hands to deal %s damage."),
            highlight(2 + ctx.level))
    end,
    tags = { "ATK", "M", "HND" },
    max_level = 1,
    color = COLOR_GRAY,
    need_target = true,
    point_cost = 1,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            deal_damage_card(ctx.caster, ctx.guid, ctx.target, 2 + ctx.level)
            return nil
        end
    },
    test = function()
        return assert_cast_damage("MELEE_HIT", 2)
    end
})

register_event("START", {
    name = "Welcome!",
    description = [[Welcome to *End of Eden*!

This game is a roguelite deckbuilder where you explore a post-apocalyptic world, collect cards and artifacts and fight enemies. **Try to stay alive as long as possible!**

**Lets start with some keyboard shortcuts**

- *ESC* - Open the menu where you can see your cards, artifacts, ... or abort choices
- *SPACE* - End your turn
- *ARROW LEFT / ARROW RIGHT* - Select card or enemy to hit
- *ENTER* - Confirm your choice
- *X* - If you hover over a enemy and press X, you can see more infos about a enemy
- *S* - Open player status

You can also use the **mouse** to select cards, enemies and click buttons.

**Cards**

You have a deck of cards that you can use to attack, defend or apply status effects. You can see your cards in the bottom of the screen. All cards cost action points that reset on each turn. Use them wisely!

**Combat**

If you press Continue you will fight a dummy enemy. See if you are able to kill it!

]],
    choices = {
        {
            description = "Continue",
            callback = function()
                return nil
            end
        },
    },
    on_enter = function()
    end,
    on_end = function(ctx)
        add_actor_by_enemy("TUTORIAL_DUMMY_1")
        give_player_gold(500)
        give_card("MELEE_HIT", PLAYER_ID)
        set_event("TUTORIAL_1")
        return GAME_STATE_FIGHT
    end
})

register_event("TUTORIAL_1", {
    name = "Status Effects",
    description = [[*Awesome! You have defeated the dummy enemy!*

Now you will face a enemy that will apply a *status effect* on you. Status effects can be positive or negative and can be applied by cards, enemies or other sources. Status effects that are applied to you are shown on the bottom of the screen. You can click on them or press *S* to see more information.

If you press Continue you will fight some dummy enemies. See if you are able to kill them!
]],
    choices = {
        {
            description = "Continue",
            callback = function()
                return nil
            end
        },
    },
    on_enter = function()
    end,
    on_end = function(ctx)
        add_actor_by_enemy("TUTORIAL_DUMMY_1")
        add_actor_by_enemy("TUTORIAL_DUMMY_2")
        set_event("TUTORIAL_2")
        give_card("BLOCK", PLAYER_ID)
        return GAME_STATE_FIGHT
    end
})

register_event("TUTORIAL_2", {
    name = "The Merchant",
    description = [[*Awesome! You have defeated the dummy enemies!*

Every now and then you will encounter a merchant. The merchant will offer you cards and artifacts that you can buy with gold. You can also remove or upgrade cards. Gold is earned by defeating enemies.

If you press Continue you will meet the merchant. *Try to buy or upgrade something!*
]],
    choices = {
        {
            description = "Continue",
            callback = function()
                return nil
            end
        },
    },
    on_enter = function()
    end,
    on_end = function(ctx)
        set_event("TUTORIAL_3")
        return GAME_STATE_MERCHANT
    end
})

register_event("TUTORIAL_3", {
    name = "Finished!",
    description = [[*Awesome! You have bought some stuff!*

This is the end of the tutorial. You can now continue to explore the world and fight enemies. Good luck!
]],
    choices = {
        {
            description = "Continue",
            callback = function()
                return nil
            end
        },
    },
    on_enter = function()
    end,
    on_end = function(ctx)
        add_actor_by_enemy("TUTORIAL_DUMMY_1")
        add_actor_by_enemy("TUTORIAL_DUMMY_2")
        add_actor_by_enemy("TUTORIAL_DUMMY_1")
        set_event("TUTORIAL_4")
        return GAME_STATE_FIGHT
    end
})


register_event("TUTORIAL_4", {
    name = "Be gone!",
    description = [[*It is time to go...*]],
    choices = {
        {
            description = "Continue",
            callback = function()
                return nil
            end
        },
    },
    on_enter = function()
    end,
    on_end = function(ctx)
        deal_damage("TUTORIAL", PLAYER_ID, 1000, true)
        return GAME_STATE_GAMEOVER
    end
})
