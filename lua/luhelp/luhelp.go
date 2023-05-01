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
	"strings"
)

const LuaTag = "lua"

func baseToLua(val any) lua.LValue {
	switch val := val.(type) {
	case float64:
		return lua.LNumber(val)
	case float32:
		return lua.LNumber(val)
	case int:
		return lua.LNumber(val)
	case int64:
		return lua.LNumber(val)
	case string:
		return lua.LString(val)
	case bool:
		return lua.LBool(val)
	}
	return lua.LNil
}

var intType = reflect.TypeOf(int(0))

func ToLua(state *lua.LState, val any) lua.LValue {
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
		return ToLua(state, valValue.Elem().Interface())
	case reflect.Struct:
		s := structs.New(val)
		s.TagName = LuaTag
		return ToLua(state, s.Map())
	case reflect.Map:
		resultTable := state.NewTable()
		valValue := reflect.ValueOf(val)
		keys := valValue.MapKeys()
		for i := range keys {
			if keys[i].Kind() == reflect.String {
				key := strings.ToLower(stringy.New(keys[i].Interface().(string)).SnakeCase().Get())
				luaVal := ToLua(state, valValue.MapIndex(keys[i]).Interface())
				resultTable.RawSetString(key, luaVal)
			} else if keys[i].CanConvert(intType) {
				resultTable.RawSetInt(keys[i].Convert(intType).Interface().(int), ToLua(state, valValue.MapIndex(keys[i]).Interface()))
			}
		}
		return resultTable
	case reflect.Slice:
		resultTable := state.NewTable()
		valValue := reflect.ValueOf(val)
		for i := 0; i < valValue.Len(); i++ {
			resultTable.Append(ToLua(state, valValue.Index(i).Interface()))
		}
		return resultTable
	default:
		return baseToLua(val)
	}
}

var noProtect = os.Getenv("EOE_NO_PROTECT") == "1"

// BindToLua will create a OwnedCallback from a lua function and state.
func BindToLua(state *lua.LState, value lua.LValue) OwnedCallback {
	return func(args ...any) (any, error) {
		// Call our lua function
		if err := state.CallByParam(lua.P{
			Fn:      value,
			NRet:    1,
			Protect: !noProtect,
		}, lo.Map(args, func(item any, index int) lua.LValue {
			return ToLua(state, item)
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
