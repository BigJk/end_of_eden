# Lua Documentation

- Lua 5.1 (+ 5.2 goto statement) supported
- The lua code tries to conform to the [luarocks style guide](https://github.com/luarocks/lua-style-guide).
- The [luafun](https://github.com/luafun/luafun) functional library is available by default to provide functions like `map`, `filter`, etc. which are very helpful. Check the [luafun docs](https://luafun.github.io/index.html) for more information.
- If you are new to lua: [Learn Lua in 15 Minutes](https://tylerneylon.com/a/learn-lua/)
- For usage examples check the game scripts in `./assets/scripts`

## Game API

There is a multitude of functions and variables globally available to access and mutate the current game state. An automated documentation is generated:

- [Lua Game API](./LUA_API_DOCS.md)

## Modding

A mod is nothing more than a `meta.json` and a bunch of lua files. The `meta.json` defines the basic information about the mod like name, author etc. and the lua files contain the content. You can enable and disable mods via the mods menu found in the main menu.

- All lua files found in the mod folder are loaded when the mod is run.
- Check the [Lua Game API](./LUA_API_DOCS.md), especially the **Content Registry** for information.
- Check `/mods/example_mod` for the bare minimum.
- Check `/assets/scripts` for usage examples.

### `meta.json` example

```json
{
  "name": "Example Mod",
  "author": "BigJk",
  "description": "Serve as example",
  "version": "0.0.1",
  "url": ""
}
```

### Mod Loading Order

In the modding menu you can define the loading order of the mods. The mods are loaded from top to bottom. This is important if mods overwrite content of the base game, like changing the `START` event. If multiple mods change it the last loaded mod will be the last to overwrite it and so its event is used. Keep this in mind when organizing mods.

## Callbacks

### General

- Every callback function is always called with one arg called `ctx`, which is a table that contains some contextual data and is expected to return `nil` if no other data is returned.
- `type_id` always contains the type id of the instance that the callback is executed on, so if a `BLOCK` card is `CallbackOnCast`, then `ctx.type_id == "BLOCK"`
- `guid` always contains the id of the instance, so the id to the instance of the card, actor, status_effect etc.
- Some callbacks have different `ctx` values depending on if a card, artifact or status effect is executed. For example the `stacks` value will only be present for status_effects.
- For lua all callback names are snake case, so `CallbackOnActorDie` is `on_actor_die`.

### Example

```lua
register_artifact(
    "GIGANTIC_STRENGTH",
    {
        name = "Stone Of Gigantic Strength",
        description = "Double all damage dealt.",
        price = 250,
        order = 0,
        -- Callbacks
        callbacks = {
            on_damage_calc = function(ctx)
                -- If the source of the damage is the owner of the artifact
                -- then double the damage.
                if ctx.source == ctx.owner then
                    return ctx.damage * 2
                end
                return nil
            end,
        }
        -- Callbacks
    }
);
```

### Quick Overview

```
Callback Type                 : ctx values
CallbackOnActorDie            : type_id guid source target owner damage
CallbackOnActorDie            : type_id guid source target owner stacks damage
CallbackOnCast                : type_id guid caster target level
CallbackOnDamage              : type_id guid source target owner damage
CallbackOnDamage              : type_id guid source target owner stacks damage
CallbackOnDamageCalc          : type_id guid source target owner damage
CallbackOnDamageCalc          : type_id guid source target owner stacks damage
CallbackOnHealCalc            : type_id guid source target owner heal
CallbackOnHealCalc            : type_id guid source target owner stacks heal
CallbackOnInit                : type_id guid
CallbackOnPickUp              : type_id guid owner
CallbackOnPlayerTurn          : type_id guid owner round
CallbackOnPlayerTurn          : type_id guid owner round stacks
CallbackOnPlayerTurn          : type_id guid round
CallbackOnRemove              : type_id guid owner
CallbackOnStatusAdd           : type_id guid
CallbackOnStatusRemove        : type_id guid owner
CallbackOnStatusStack         : type_id guid owner stacks
CallbackOnTurn                : type_id guid owner round stacks
CallbackOnTurn                : type_id guid round
```

