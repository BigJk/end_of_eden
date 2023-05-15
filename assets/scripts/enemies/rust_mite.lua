register_enemy("RUST_MITE", {
    name = "Rust Mite",
    description = "Loves to eat metal.",
    look = "/v\\",
    color = "#e6e65a",
    initial_hp = 22,
    max_hp = 22,
    gold = 10,
    intend = function(ctx)
        if ctx.round % 4 == 0 then
            return "Gather strength"
        end

        return "Deal " .. highlight(6) .. " damage"
    end,
    callbacks = {
        on_turn = function(ctx)
            if ctx.round % 4 == 0 then
                give_status_effect("RITUAL", ctx.guid)
            else
                deal_damage(ctx.guid, PLAYER_ID, 6)
            end

            return nil
        end
    }
})