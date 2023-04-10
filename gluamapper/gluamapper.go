// gluamapper provides an easy way to map GopherLua tables to Go structs.
package gluamapper

/*

	The MIT License (MIT)

	Copyright (c) 2015 Yusuke Inuzuka

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in all
	copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
	SOFTWARE.

*/

import (
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/yuin/gopher-lua"
	"regexp"
	"strings"
)

// Option is a configuration that is used to create a new mapper.
type Option struct {
	// Function to convert a lua table key to Go's one. This defaults to "ToUpperCamelCase".
	NameFunc func(string) string

	// Returns error if unused keys exist.
	ErrorUnused bool

	// A struct tag name for lua table keys . This defaults to "gluamapper"
	TagName string

	DecodeHook any

	FnHook func(value lua.LValue) any
}

// Mapper maps a lua table to a Go struct pointer.
type Mapper struct {
	Option Option
}

// NewMapper returns a new mapper.
func NewMapper(opt Option) *Mapper {
	if opt.NameFunc == nil {
		opt.NameFunc = ToUpperCamelCase
	}
	if opt.TagName == "" {
		opt.TagName = "gluamapper"
	}
	if opt.FnHook == nil {
		opt.FnHook = func(value lua.LValue) any {
			return value
		}
	}
	return &Mapper{opt}
}

// Map maps the lua table to the given struct pointer.
func (mapper *Mapper) Map(tbl *lua.LTable, st any) error {
	opt := mapper.Option
	mp, ok := ToGoValue(tbl, opt).(map[any]any)
	if !ok {
		return errors.New("arguments #1 must be a table, but got an array")
	}
	config := &mapstructure.DecoderConfig{
		DecodeHook:       opt.DecodeHook,
		WeaklyTypedInput: true,
		Result:           st,
		TagName:          opt.TagName,
		ErrorUnused:      opt.ErrorUnused,
	}
	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}
	return decoder.Decode(mp)
}

// Map maps the lua table to the given struct pointer with default options.
func Map(tbl *lua.LTable, st any) error {
	return NewMapper(Option{}).Map(tbl, st)
}

// Id is an Option.NameFunc that returns given string as-is.
func Id(s string) string {
	return s
}

var camelre = regexp.MustCompile(`_([a-z])`)

// ToUpperCamelCase is an Option.NameFunc that converts strings from snake case to upper camel case.
func ToUpperCamelCase(s string) string {
	return strings.ToUpper(string(s[0])) + camelre.ReplaceAllStringFunc(s[1:len(s)], func(s string) string { return strings.ToUpper(s[1:len(s)]) })
}

// ToGoValue converts the given LValue to a Go object.
func ToGoValue(lv lua.LValue, opt Option) any {
	if lv.Type() == lua.LTFunction {
		return opt.FnHook(lv)
	}

	switch v := lv.(type) {
	case *lua.LNilType:
		return nil
	case lua.LBool:
		return bool(v)
	case lua.LString:
		return string(v)
	case lua.LNumber:
		return float64(v)
	case *lua.LTable:
		maxn := v.MaxN()
		if maxn == 0 { // table
			ret := make(map[any]any)
			v.ForEach(func(key, value lua.LValue) {
				keystr := fmt.Sprint(ToGoValue(key, opt))
				ret[opt.NameFunc(keystr)] = ToGoValue(value, opt)
			})
			return ret
		} else { // array
			ret := make([]any, 0, maxn)
			for i := 1; i <= maxn; i++ {
				ret = append(ret, ToGoValue(v.RawGetInt(i), opt))
			}
			return ret
		}
	default:
		return v
	}
}
