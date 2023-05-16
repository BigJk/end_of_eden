register_artifact("SHORT_RADIANCE", {
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
        end
    }
});
