register_enemy("PLASMA_GOLEM", {
    name = "Plasma Golem",
    description = "A golem made of pure plasma energy.",
    look = [[
  /\
 /  \
/_xx_\]],
    color = "#ff69b4",
    initial_hp = 12,
    max_hp = 12,
    gold = 80,
    intend = function(ctx)
        if ctx.round % 2 == 0 then
            return "Charge up for a powerful plasma attack"
        else
            return "Deal " .. highlight(8) .. " damage"
        end
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round % 2 == 0 then
            else
                deal_damage(ctx.guid, PLAYER_ID, 8)
            end
            return nil
        end
    }
})
