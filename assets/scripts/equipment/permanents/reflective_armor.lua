register_artifact("REFLECTIVE_ARMOR", {
    name = l("artifacts.REFLECTIVE_ARMOR.name", "Reflective Armor"),
    description = l("artifacts.REFLECTIVE_ARMOR.description", "Reflects 25% of the damage back to the attacker."),
    tags = { "ARMOR", "_ACT_0" },
    price = 300,
    order = 0,
    callbacks = {
        on_damage_calc = function(ctx)
            if ctx.target == ctx.owner then
                local reflected_damage = ctx.damage * 0.25
                deal_damage(ctx.target, ctx.source, reflected_damage, true)
            end
            return ctx.damage
        end
    }
})
