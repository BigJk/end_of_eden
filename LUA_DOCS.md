# Lua Documentation

- Lua 5.1 (+ 5.2 goto statement) supported
- The lua code tries to conform to the [luarocks style guide](https://github.com/luarocks/lua-style-guide).
- The [luafun](https://github.com/luafun/luafun) functional library is available by default to provide functions like ``map``, ``filter``, etc. which are very helpful. Check the [luafun docs](https://luafun.github.io/index.html) for more information.

## Known Problems

### luafun: ``iter``

Don't use ``iter`` from lua fun, instead use ``pairs`` or ``ipairs``. Otherwise the game will crash with some tables that come from the game.

```lua
-- Good
each(function(val)
    give_status_effect("VULNERABLE", val)
end, pairs(get_opponent_guids(ctx.owner))) -- use pairs

--- Bad
each(function(val)
    give_status_effect("VULNERABLE", val)
end, iter(get_opponent_guids(ctx.owner))) -- iter bad!

each(function(val)
    give_status_effect("VULNERABLE", val)
end, get_opponent_guids(ctx.owner)) -- neither iter nor pairs... bad!
```

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
add_actor_by_enemy             : function
give_artifact                  : function
remove_artifact                : function
give_status_effect             : function
remove_status_effect           : function
add_status_effect_stacks       : function
set_status_effect_stacks       : function
give_card                      : function
remove_card                    : function
cast_card                      : function
get_cards                      : function
deal_damage                    : function
deal_damage_multi              : function
heal                           : function
player_draw_card               : function
give_player_gold               : function
```

## Callbacks

### General

- Every callback function is always called with one arg called ``ctx``, which is a table that contains some contextual data
- ``type_id`` always contains the type id of the instance that the callback is executed on, so if a ``BLOCK`` card is ``CallbackOnCast``, then ``ctx.type_id == "BLOCK"``
- ``guid`` always contains the id of the instance, so the id to the instance of the card, actor, status_effect etc.
- Some callbacks have different ``ctx`` values depending on if a card, artifact or status effect is executed. For example the ``stacks`` value will only be present for status_effects.
- For lua all callback names are snake case, so ``CallbackOnActorDie`` is ``on_actor_die``.

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
CallbackOnRemove              : type_id guid owner 
CallbackOnStatusAdd           : type_id guid 
CallbackOnStatusRemove        : type_id guid owner 
CallbackOnStatusStack         : type_id guid owner stacks 
CallbackOnTurn                : type_id guid owner round stacks 
CallbackOnTurn                : type_id guid round 
```