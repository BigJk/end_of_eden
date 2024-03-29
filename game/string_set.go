package game

import (
	"bytes"
	"encoding/gob"
	"github.com/samber/lo"
	"sort"
)

// StringSet represents a string set that can be serialized by Gob.
//
// Note: As the GobDecode needs to overwrite its receiver we need to have the map
// behind a struct pointer.
type StringSet struct {
	values map[string]struct{}
}

func NewStringSet() *StringSet {
	return &StringSet{
		values: map[string]struct{}{},
	}
}

func (s *StringSet) Has(val string) bool {
	_, ok := s.values[val]
	return ok
}

func (s *StringSet) Add(val string) {
	s.values[val] = struct{}{}
}

func (s *StringSet) Remove(val string) {
	delete(s.values, val)
}

func (s *StringSet) Append(vals ...string) {
	for _, val := range vals {
		s.Add(val)
	}
}

func (s *StringSet) Clear() {
	for key := range s.values {
		delete(s.values, key)
	}
}

func (s *StringSet) ToSlice() []string {
	keys := lo.Keys(s.values)
	sort.Strings(keys)
	return keys
}

func (s *StringSet) Clone() *StringSet {
	return &StringSet{values: CopyMap(s.values)}
}

func (s *StringSet) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s.ToSlice())
	return buf.Bytes(), err
}

func (s *StringSet) GobDecode(data []byte) error {
	*s = StringSet{values: map[string]struct{}{}}

	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var keys []string
	if err := dec.Decode(&keys); err != nil {
		return err
	}

	s.Append(keys...)

	return nil
}

// CopyMap copies a map. If the value is a pointer, the pointer is copied, not the value.
func CopyMap[K comparable, V any](m map[K]V) map[K]V {
	result := make(map[K]V)
	for k, v := range m {
		result[k] = v
	}
	return result
}
