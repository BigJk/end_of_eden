register_artifact("DEFLECTOR_SHIELD", {
    name = "Deflector Shield",
    description = "Gain 8 block at the start of combat.",
    price = 50,
    order = 0,
    callbacks = {
        on_player_turn = function(ctx)
            if ctx.round == 0 then
                give_status_effect("BLOCK", ctx.owner, 8)
            end
            return nil
        end
    },
    test = function(ctx)
        add_actor_by_enemy("DUMMY")

        status_effects = get_actor_status_effects(PLAYER_ID)

        -- check if length of status_effects is 1
        if #status_effects ~= 1 then
            return "Expected 1 status effect, got " .. #status_effects
        end

        -- check if the status effect is BLOCK
        block = get_status_effect_instance(status_effects[1])
        if block.type_id ~= "BLOCK" then
            return "Expected BLOCK status effect, got " .. tostring(block.type_id)
        end

        -- check if the block amount is 8
        if block.stacks ~= 8 then
            return "Expected 8 block, got " .. tostring(block.stacks)
        end
    end
});
