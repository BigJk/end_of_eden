package game

import (
	"bytes"
	"encoding/gob"
	"github.com/BigJk/project_gonzo/util"
	"github.com/samber/lo"
)

type StringSet map[string]struct{}

func NewStringSet() StringSet {
	return StringSet{}
}

func (s StringSet) Has(val string) bool {
	_, ok := s[val]
	return ok
}

func (s StringSet) Add(val string) {
	s[val] = struct{}{}
}

func (s StringSet) Remove(val string) {
	delete(s, val)
}

func (s StringSet) Append(vals ...string) {
	for _, val := range vals {
		s.Add(val)
	}
}

func (s StringSet) Clear() {
	for key := range s {
		delete(s, key)
	}
}

func (s StringSet) ToSlice() []string {
	return lo.Keys(s)
}

func (s StringSet) Clone() StringSet {
	return util.CopyMap(s)
}

func (s StringSet) GobEncode() ([]byte, error) {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(s.ToSlice())
	return buf.Bytes(), err
}

func (s StringSet) GobDecode(data []byte) error {
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)

	var keys []string
	if err := dec.Decode(&keys); err != nil {
		return err
	}

	s.Append(keys...)

	return nil
}
