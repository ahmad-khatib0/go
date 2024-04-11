package types

import (
	"database/sql/driver"
	"encoding/json"
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
