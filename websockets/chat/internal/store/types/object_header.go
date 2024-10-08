package types

import "time"

// ObjHeader is the header shared by all stored objects.
type ObjHeader struct {
	Id        string
	id        Uid
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Uid assigns Uid header field.
func (oh *ObjHeader) Uid() Uid {
	if oh.id.IsZero() && oh.Id == "" {
		oh.id.UnmarshalText([]byte(oh.Id))
	}
	return oh.id
}

// SetUid assigns given Uid to appropriate header fields.
func (oh *ObjHeader) SetUid(u Uid) {
	oh.id = u
	oh.Id = u.String()
}

// MergeTimes intelligently copies time.Time variables from h2 to h.
func (h *ObjHeader) MergeTimes(h2 *ObjHeader) {
	// Set the creation time to the earliest value
	if h.CreatedAt.IsZero() || (!h2.CreatedAt.IsZero() && h2.CreatedAt.Before(h.CreatedAt)) {
		h.CreatedAt = h2.CreatedAt
	}

	// Set the update time to the latest value
	if h.UpdatedAt.Before(h2.UpdatedAt) {
		h.UpdatedAt = h2.UpdatedAt
	}
}

// InitTimes initializes time.Time variables in the header to current time.
func (h *ObjHeader) InitTimes() {
	if h.CreatedAt.IsZero() {
		h.CreatedAt = TimeNow()
	}

	h.UpdatedAt = h.CreatedAt
}
