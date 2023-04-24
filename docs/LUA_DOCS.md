# Lua Documentation

- Lua 5.1 (+ 5.2 goto statement) supported
- The lua code tries to conform to the [luarocks style guide](https://github.com/luarocks/lua-style-guide).
- The [luafun](https://github.com/luafun/luafun) functional library is available by default to provide functions like ``map``, ``filter``, etc. which are very helpful. Check the [luafun docs](https://luafun.github.io/index.html) for more information.
- If you are new to lua: [Learn Lua in 15 Minutes](https://tylerneylon.com/a/learn-lua/)

## Global Function & Values

There is a multitude of functions and variables globally available to access and mutate the current game state.

### Quick Overview

```
name                           : type
PLAYER_ID                      : value
GAME_STATE_FIGHT               : value
GAME_STATE_EVENT               : value
GAME_STATE_MERCHANT            : value
GAME_STATE_RANDOM              : value
DECAY_ONE                      : value
DECAY_ALL                      : value
DECAY_NONE                     : value
guid                           : function
text_bold                      : function
text_italic                    : function
text_underline                 : function
text_color                     : function
text_bg                        : function
log_i                          : function
log_w                          : function
log_d                          : function
log_s                          : function
debug_value                    : function
debug_log                      : function
play_audio                     : function
set_event                      : function
set_game_state                 : function
set_fight_description          : function
get_fight_round                : function
get_stages_cleared             : function
get_fight                      : function
get_event_history              : function
had_event                      : function
get_player                     : function
get_actor                      : function
get_opponent_by_index          : function
get_opponent_count             : function
get_opponent_guids             : function
remove_actor                   : function
actor_add_max_hp               : function
add_actor_by_enemy             : function
give_artifact                  : function
remove_artifact                : function
get_random_artifact_type       : function
get_artifact                   : function
get_artifact_instance          : function
give_status_effect             : function
remove_status_effect           : function
add_status_effect_stacks       : function
set_status_effect_stacks       : function
get_actor_status_effects       : function
get_status_effect              : function
get_status_effect_instance     : function
give_card                      : function
remove_card                    : function
cast_card                      : function
get_cards                      : function
get_card                       : function
get_card_instance              : function
deal_damage                    : function
deal_damage_multi              : function
heal                           : function
player_draw_card               : function
give_player_gold               : function
get_merchant                   : function
add_merchant_card              : function
add_merchant_artifact          : function
get_merchant_gold_max          : function
gen_face                       : function
random_card                    : function
random_artifact                : function
```

## Callbacks

### General

- Every callback function is always called with one arg called ``ctx``, which is a table that contains some contextual data and is expected to return ``nil`` if no other data is returned.
- ``type_id`` always contains the type id of the instance that the callback is executed on, so if a ``BLOCK`` card is ``CallbackOnCast``, then ``ctx.type_id == "BLOCK"``
- ``guid`` always contains the id of the instance, so the id to the instance of the card, actor, status_effect etc.
- Some callbacks have different ``ctx`` values depending on if a card, artifact or status effect is executed. For example the ``stacks`` value will only be present for status_effects.
- For lua all callback names are snake case, so ``CallbackOnActorDie`` is ``on_actor_die``.

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
                if ctx.target == ctx.owner then
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