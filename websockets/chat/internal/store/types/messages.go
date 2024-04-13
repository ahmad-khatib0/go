package types

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Message is a stored {data} message
type Message struct {
	ObjHeader `bson:",inline"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bson:",omitempty"`

	// ID of the hard-delete operation
	DelId int `json:"del_id,omitempty" bson:",omitempty"`
	// List of users who have marked this message as soft-deleted
	DeletedFor []SoftDelete `json:"deleted_for,omitempty" bson:",omitempty"`
	SeqId      int
	Topic      string
	// Sender's user ID as string (without 'usr' prefix), could be empty.
	From    string
	Head    MessageHeaders `json:"dead,omitempty" bson:",omitempty"`
	Content interface{}
}

// MessageHeaders is needed to attach Scan() to.
type MessageHeaders map[string]interface{}

// Scan implements sql.Scanner interface.
func (mh *MessageHeaders) Scan(val interface{}) error {
	return json.Unmarshal(val.([]byte), mh)
}

// Value implements sql's driver.Valuer interface.
func (mh MessageHeaders) Value() (driver.Value, error) {
	return json.Marshal(mh)
}

// SoftDelete is a single DB record of soft-deletetion.
type SoftDelete struct {
	User  string
	DelId int
}

// DelMessage is a log entry of a deleted message range.
type DelMessage struct {
	ObjHeader   `bson:",inline"`
	Topic       string
	DeletedFor  string
	DelId       int
	SeqIdRanges []Range
}

// Range is a range of message SeqIDs. Low end is inclusive (closed),
//
// high end is exclusive (open): [Low, Hi).
//
// If the range contains just one ID, Hi is set to 0
type Range struct {
	Low int
	Hi  int `json:"hi,omitempty" bson:",omitempty"`
}
