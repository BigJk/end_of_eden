register_artifact("ARM_MOUNTED_GUN", {
    name = "Arm Mounted Gun",
    description = "Weapon that is mounted on your arm. It is very powerful.",
    tags = { "ARM" },
    price = 250,
    order = 0,
    callbacks = {
        on_pick_up = function(ctx)
            clear_artifacts_by_tag("ARM", { ctx.guid })
            clear_cards_by_tag("ARM")
            give_card("ARM_MOUNTED_GUN", PLAYER_ID)
            return nil
        end
    }
});

register_card("ARM_MOUNTED_GUN", {
    name = l("cards.ARM_MOUNTED_GUN.name", "Arm Mounted Gun"),
    description = l("cards.ARM_MOUNTED_GUN.description", "Exhaust. Use your arm mounted gun to deal 15 (+3 for each upgrade) damage."),
    state = function(ctx)
        return string.format(l("cards.ARM_MOUNTED_GUN.state", "Use your arm mounted gun to deal %s damage."), highlight(7 + ctx.level * 3))
    end,
    tags = { "ATK", "R", "T", "ARM" },
    max_level = 1,
    color = COLOR_GRAY,
    need_target = true,
    does_exhaust = true,
    point_cost = 3,
    price = -1,
    callbacks = {
        on_cast = function(ctx)
            deal_damage_card(ctx.caster, ctx.guid, ctx.target, 7 + ctx.level * 3)
            return nil
        end
    },
    test = function ()
        return assert_cast_damage("ARM_MOUNTED_GUN", 7)
    end
})