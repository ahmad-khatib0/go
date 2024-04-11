package types

import "time"

// ObjHeader is the header shared by all stored objects.
type ObjHeader struct {
	ID        string
	id        Uid
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Uid assigns Uid header field.
func (oh *ObjHeader) Uid() Uid {
	if oh.id.IsZero() && oh.ID == "" {
		oh.id.UnmarshalText([]byte(oh.ID))
	}
	return oh.id
}

// SetUid assigns given Uid to appropriate header fields.
func (oh *ObjHeader) SetUid(u Uid) {
	oh.id = u
	oh.ID = u.String()
}
