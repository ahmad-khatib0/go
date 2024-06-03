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

// RangeSorter is a helper type required by 'sort' package.
type RangeSorter []Range

// Len is the length of the range.
func (rs RangeSorter) Len() int {
	return len(rs)
}

// Swap swaps two items in a slice.
func (rs RangeSorter) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// Less is a comparator. Sort by Low ascending, then sort by Hi descending
func (rs RangeSorter) Less(i, j int) bool {
	if rs[i].Low < rs[j].Low {
		return true
	}
	if rs[i].Low == rs[j].Low {
		return rs[i].Hi >= rs[j].Hi
	}
	return false
}

// Normalize ranges - remove overlaps: [1..4],[2..4],[5..7] -> [1..7].
// The ranges are expected to be sorted.
// Ranges are inclusive-inclusive, i.e. [1..3] -> 1, 2, 3.
func (rs RangeSorter) Normalize() RangeSorter {
	if ll := rs.Len(); ll > 1 {
		prev := 0
		for i := 1; i < ll; i++ {
			if rs[prev].Low == rs[i].Low {
				// Earlier range is guaranteed to be wider or equal to the later range,
				// collapse two ranges into one (by doing nothing)
				continue
			}
			// Check for full or partial overlap
			if rs[prev].Hi > 0 && rs[prev].Hi+1 >= rs[i].Low {
				// Partial overlap
				if rs[prev].Hi < rs[i].Hi {
					rs[prev].Hi = rs[i].Hi
				}
				// Otherwise the next range is fully within the previous range, consume it by doing nothing.
				continue
			}
			// No overlap
			prev++
		}
		rs = rs[:prev+1]
	}

	return rs
}
