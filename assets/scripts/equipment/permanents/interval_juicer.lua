register_artifact("INTERVA_JUICER", {
    name = "Interval Juicer",
    description = highlight("Heal 2") .. " at the beginning of combat",
    tags = { "_ACT_0" },
    price = 200,
    order = 0,
    callbacks = {
        on_player_turn = function(ctx)
            if ctx.owner == PLAYER_ID and ctx.round == 0 then
                heal(PLAYER_ID, PLAYER_ID, 2)
            end
            return nil
        end
    }
});
