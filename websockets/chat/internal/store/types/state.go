package types

import (
	"database/sql/driver"
	"errors"
	"strings"
)

// ObjState represents information on objects state,
// such as an indication that User or Topic is suspended/soft-deleted.
type ObjState int

const (
	// StateOK indicates normal user or topic.
	StateOK ObjState = 0
	// StateSuspended indicates suspended user or topic.
	StateSuspended ObjState = 10
	// StateDeleted indicates soft-deleted user or topic.
	StateDeleted ObjState = 20
	// StateUndefined indicates state which has not been set explicitly.
	StateUndefined ObjState = 30
)

// NewObjState parses string into an ObjState.
func NewObjState(in string) (ObjState, error) {
	in = strings.ToLower(in)
	switch in {
	case "", "ok":
		return StateOK, nil
	case "susp":
		return StateSuspended, nil
	case "del":
		return StateDeleted, nil
	case "undef":
		return StateUndefined, nil
	}
	// This is the default.
	return StateOK, errors.New("failed to parse object state")

}

// String returns string representation of ObjState.
func (os ObjState) String() string {
	switch os {
	case StateOK:
		return "ok"
	case StateSuspended:
		return "susp"
	case StateDeleted:
		return "del"
	case StateUndefined:
		return "undef"
	}
	return ""
}

// MarshalJSON converts ObjState to a quoted string.
func (os ObjState) MarshalJSON() ([]byte, error) {
	return append(append([]byte{'"'}, []byte(os.String())...), '"'), nil
}

// UnmarshalJSON reads ObjState from a quoted string.
func (os *ObjState) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("syntax error")
	}

	state, err := NewObjState(string(b[1 : len(b)-1]))
	if err == nil {
		*os = state
	}
	return err
}

// Scan is an implementation of sql.Scanner interface. It expects the
// value to be a byte slice representation of an ASCII string.
func (os *ObjState) Scan(val interface{}) error {
	switch intval := val.(type) {
	case int64:
		*os = ObjState(intval)
		return nil
	}
	return errors.New("ObjState: data is not an int64 when Scanning")
}

// Value is an implementation of sql.driver.Valuer interface.
func (os ObjState) Value() (driver.Value, error) {
	return int64(os), nil
}
