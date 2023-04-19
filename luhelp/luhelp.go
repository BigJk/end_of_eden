package luhelp

import (
	"fmt"
	"github.com/fatih/structs"
	"github.com/gobeam/stringy"
	"github.com/samber/lo"
	"github.com/sanity-io/litter"
	lua "github.com/yuin/gopher-lua"
	"os"
	"reflect"
)

const LuaTag = "lua"

func baseToLua(val any) lua.LValue {
	switch val.(type) {
	case float64:
		return lua.LNumber(val.(float64))
	case float32:
		return lua.LNumber(val.(float32))
	case int:
		return lua.LNumber(val.(int))
	case int64:
		return lua.LNumber(val.(int64))
	case string:
		return lua.LString(val.(string))
	case bool:
		return lua.LBool(val.(bool))
	}
	return lua.LNil
}

var intType = reflect.TypeOf(int(0))

func ToLua(val any) lua.LValue {
	if val == nil {
		return lua.LNil
	}

	valType := reflect.TypeOf(val)

	switch valType.Kind() {
	case reflect.Pointer:
		valValue := reflect.ValueOf(val)
		if valValue.IsNil() {
			return lua.LNil
		}
		return ToLua(valValue.Elem().Interface())
	case reflect.Struct:
		s := structs.New(val)
		s.TagName = LuaTag
		return ToLua(s.Map())
	case reflect.Map:
		resultTable := &lua.LTable{}
		valValue := reflect.ValueOf(val)
		keys := valValue.MapKeys()
		for i := range keys {
			if keys[i].Kind() == reflect.String {
				key := stringy.New(keys[i].Interface().(string)).SnakeCase().Get()
				resultTable.RawSetString(key, ToLua(valValue.MapIndex(keys[i]).Interface()))
			} else if keys[i].CanConvert(intType) {
				resultTable.RawSetInt(keys[i].Convert(intType).Interface().(int), ToLua(valValue.MapIndex(keys[i]).Interface()))
			}
		}
		return resultTable
	case reflect.Slice:
		resultTable := &lua.LTable{}
		valValue := reflect.ValueOf(val)
		for i := 0; i < valValue.Len(); i++ {
			resultTable.Append(ToLua(valValue.Index(i).Interface()))
		}
		return resultTable
	default:
		return baseToLua(val)
	}
}

var noProtect = os.Getenv("PG_NO_PROTECT") == "1"

// BindToLua will create a OwnedCallback from a lua function and state.
func BindToLua(state *lua.LState, value lua.LValue) OwnedCallback {
	return func(args ...any) (any, error) {
		// Call our lua function
		if err := state.CallByParam(lua.P{
			Fn:      value,
			NRet:    1,
			Protect: !noProtect,
		}, lo.Map(args, func(item any, index int) lua.LValue {
			return ToLua(item)
		})...); err != nil {
			return nil, err
		}

		// Fetch return value
		ret := state.Get(-1)
		state.Pop(1)

		// Parse to accepted return values
		switch ret.Type() {
		case lua.LTString:
			return lua.LVAsString(ret), nil
		case lua.LTNumber:
			return float64(lua.LVAsNumber(ret)), nil
		case lua.LTBool:
			return lua.LVAsBool(ret), nil
		case lua.LTTable:
			mapper := NewMapper(state)
			maxn := value.(*lua.LTable).MaxN()
			if maxn == 0 {
				var data map[string]any
				if err := mapper.Map(ret.(*lua.LTable), &data); err != nil {
					return nil, err
				}
				return data, nil
			}

			data := make([]any, 0)
			if err := mapper.Map(ret.(*lua.LTable), &data); err != nil {
				return nil, err
			}
			return data, nil
		}

		// Don't error for now
		return nil, nil
	}
}

func ToString(val lua.LValue, mapper *Mapper) string {
	switch val.Type() {
	case lua.LTString:
		return lua.LVAsString(val)
	case lua.LTNumber:
		return fmt.Sprint(float64(lua.LVAsNumber(val)))
	case lua.LTBool:
		return fmt.Sprint(lua.LVAsBool(val))
	case lua.LTTable:
		maxn := val.(*lua.LTable).MaxN()
		if maxn == 0 {
			var data map[string]interface{}
			if err := mapper.Map(val.(*lua.LTable), &data); err != nil {
				return "Error: " + err.Error()
			}
			return litter.Sdump(data)
		}

		ret := make([]any, 0)
		if err := mapper.Map(val.(*lua.LTable), &ret); err != nil {
			return "Error: " + err.Error()
		}
		return litter.Sdump(ret)
	case lua.LTUserData:
		return fmt.Sprint(val.(*lua.LUserData).Value)
	case lua.LTNil:
		return "nil"
	}

	return "<" + val.Type().String() + ">"
}
