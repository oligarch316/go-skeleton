package ctype

import (
	"encoding/json"
	"sort"

	"go.uber.org/zap/zapcore"
)

// StringSet TODO
type StringSet map[string]struct{}

// NewStringSet TODO
func NewStringSet(members ...string) StringSet {
	res := make(StringSet, len(members))
	res.Add(members...)
	return res
}

// Contains TODO
func (ss StringSet) Contains(member string) bool {
	_, res := ss[member]
	return res
}

// Add TODO
func (ss StringSet) Add(members ...string) {
	for _, member := range members {
		ss[member] = struct{}{}
	}
}

// MarshalJSON TODO
func (ss StringSet) MarshalJSON() ([]byte, error) {
	var (
		list = make([]string, len(ss))
		idx  int
	)

	for member := range ss {
		list[idx] = member
		idx++
	}

	sort.Strings(list)
	return json.Marshal(list)
}

// UnmarshalJSON TODO
func (ss *StringSet) UnmarshalJSON(data []byte) error {
	var tmp []string
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	*ss = make(StringSet, len(tmp))
	ss.Add(tmp...)
	return nil
}

// MarshalLogArray TODO.
func (ss StringSet) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	for member := range ss {
		enc.AppendString(member)
	}
	return nil
}
