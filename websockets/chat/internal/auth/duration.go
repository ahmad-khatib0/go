package auth

import (
	"encoding/json"
	"errors"
	"time"
)

// Duration is identical to time.Duration except it can be sanely unmarshallend from JSON.
type Duration time.Duration

// UnmarshalJSON handles the cases where duration is specified in JSON as a "5000s" string or just plain seconds.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}

	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value) * time.Second)
		return nil
	case string:
		d0, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(d0)
		return nil
	default:
		return errors.New("UnmarshalJSON invalid duration")
	}
}
