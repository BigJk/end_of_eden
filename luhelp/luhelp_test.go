package luhelp

import (
	"github.com/stretchr/testify/assert"
	lua "github.com/yuin/gopher-lua"
	"testing"
)

func TestLuHelp(t *testing.T) {
	state := lua.NewState()
	mapper := NewMapper(state)

	t.Run("Struct", func(t *testing.T) {
		type testStructInner struct {
			A int     `lua:"a"`
			B string  `lua:"b"`
			C float64 `lua:"c"`
		}

		type testStruct struct {
			Foo      string            `lua:"foo"`
			Bar      string            `lua:"bar"`
			Data     map[string]any    `lua:"data"`
			Inner    testStructInner   `lua:"inner"`
			InnerPtr *testStructInner  `lua:"innerPtr"`
			NilPtr   *testStructInner  `lua:"nilPtr"`
			Slice    []testStructInner `lua:"slice"`
			StrSlice []string          `lua:"strSlice"`
		}

		data := testStruct{
			Foo: "Hello",
			Bar: "World",
			Data: map[string]any{
				"Hello": "World",
				"1":     2.0,
			},
			Inner: testStructInner{
				A: 3,
				B: "2",
				C: 1,
			},
			InnerPtr: &testStructInner{
				A: 231,
				B: "23123",
				C: 22,
			},
			Slice: []testStructInner{
				{
					A: 1,
					B: "2",
					C: 3,
				},
				{
					A: 4,
					B: "5",
					C: 6,
				},
			},
			StrSlice: []string{"1", "2", "hello world"},
		}

		var passed testStruct
		luaVal := ToLua(data)
		assert.NoError(t, mapper.Map(luaVal.(*lua.LTable), &passed))
		assert.Equal(t, data, passed)
	})
}
