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

type perUserData struct {
	private interface{}
	want    AccessMode
	given   AccessMode
}

// LastSeenUA is a timestamp and a user agent of when the user was last seen.
type LastSeenUA struct {
	// When is the timestamp when the user was last online.
	When time.Time
	// UserAgent is the client UA of the last online access.
	UserAgent string
}

// QueryOpt is options of a query, [since, before] - both ends inclusive (closed)
type QueryOpt struct {
	// Subscription query
	User            Uid
	Topic           string
	IfModifiedSince *time.Time
	// ID-based query parameters: Messages
	Since  int
	Before int
	// Common parameter
	Limit int
}
