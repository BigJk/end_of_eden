register_artifact(
    "TEST_ARTIFACT",
    {
        Name = "Test Artifact",
        Description = "This is a cool description",
        Price = 1337,
        Order = -10,
        Callbacks = {
            OnPickUp = function(type, guid, owner)
                -- hello world
                return nil
            end,
            OnCombatEnd = function(type, guid, owner)
                -- hello world
                return nil
            end
        }
    }
);

register_artifact(
    "DOUBLE_DAMAGE",
    {
        Name = "Stone Of Gigantic Strength",
        Description = "Double all damage dealt.",
        Price = 1000,
        Order = -10,
        Callbacks = {
            OnDamageCalc = function(type, guid, source, target, damage)
                return damage * 2
            end,
        }
    }
);

register_artifact(
        "LESSER_DAMAGE_HEAL",
        {
            Name = "Repulsion Stone",
            Description = "For each damage taken heal for 2",
            Price = 200,
            Order = 0,
            Callbacks = {
                OnDamageCalc = function(type, guid, source, target, damage)
                    heal()
                    return damage
                end,
            }
        }
);