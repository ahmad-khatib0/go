package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// StringSlice is defined so Scanner and Valuer can be attached to it.
type StringSlice []string

// Scan implements sql.Scanner interface.
func (ss *StringSlice) Scan(val interface{}) error {
	if val == nil {
		return nil
	}
	return json.Unmarshal(val.([]byte), ss)
}

// Value implements sql/driver.Valuer interface.
func (ss StringSlice) Value() (driver.Value, error) {
	return json.Marshal(ss)
}

// TimeNow returns current wall time in UTC rounded to milliseconds.
func TimeNow() time.Time {
	return time.Now().UTC().Round(time.Millisecond)
}
