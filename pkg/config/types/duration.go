package ctype

import (
	"encoding/json"
	"time"
)

// Duration TODO
type Duration struct{ time.Duration }

// MarshalJSON TODO
func (d Duration) MarshalJSON() ([]byte, error) { return json.Marshal(d.String()) }

// UnmarshalJSON TODO
func (d *Duration) UnmarshalJSON(data []byte) error {
	var tmp string

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return err
	}

	d.Duration, err = time.ParseDuration(tmp)
	return err
}
