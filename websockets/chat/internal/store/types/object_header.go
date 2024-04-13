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
