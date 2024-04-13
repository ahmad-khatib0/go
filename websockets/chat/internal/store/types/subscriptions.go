package types

import "time"

// Subscription to a topic
type Subscription struct {
	ObjHeader `bson:",inline"`
	// User who has relationship with the topic
	User string
	// Topic subscribed to
	Topic     string
	DeletedAt *time.Time `bson:",omitempty"`

	// Values persisted through subscription soft-deletion

	// ID of the latest Soft-delete operation
	DelId int
	// Last SeqId reported by user as received by at least one of his sessions
	RecvSeqId int
	// Last SeqID reported read by the user
	ReadSeqId int

	// Access mode requested by this user
	ModeWant AccessMode
	// Access mode granted to this user
	ModeGiven AccessMode
	// User's private data associated with the subscription to a topic
	Private interface{}

	// Deserialized ephemeral values

	// Deserialized public value from topic or user (depends on context)
	// In case of P2P topics this is the Public value of the other user.
	public interface{}
	// In case of P2P topics this is the Trusted value of the other user.
	trusted interface{}
	// deserialized SeqID from user or topic
	seqId int
	// Deserialized TouchedAt from topic
	touchedAt time.Time
	// Timestamp & user agent of when the user was last online.
	lastSeenUA *LastSeenUA

	// P2P only. ID of the other user
	with string
	// P2P only. Default access: this is the mode given by the other user to this user
	modeDefault *DefaultAccess

	// Topic's or user's state.
	state ObjState

	// This is not a fully initialized subscription object
	dummy bool
}

// SetWith sets other user for P2P subscriptions.
func (s *Subscription) SetWith(with string) {
	s.with = with
}

// GetWith returns the other user for P2P subscriptions.
func (s *Subscription) GetWith() string {
	return s.with
}

// GetTouchedAt returns touchedAt.
func (s *Subscription) GetTouchedAt() time.Time {
	return s.touchedAt
}

// SetState assigns topic's or user's state.
func (s *Subscription) SetState(state ObjState) {
	s.state = state
}

// SetTouchedAt sets the value of touchedAt.
func (s *Subscription) SetTouchedAt(touchedAt time.Time) {
	if touchedAt.After(s.touchedAt) {
		s.touchedAt = touchedAt
	}
}

// SetSeqId sets seqId field.
func (s *Subscription) SetSeqId(id int) {
	s.seqId = id
}

// GetPublic reads value of `public`.
func (s *Subscription) GetPublic() interface{} {
	return s.public
}

// SetPublic assigns a value to `public`, otherwise not accessible from outside the package.
func (s *Subscription) SetPublic(pub interface{}) {
	s.public = pub
}

// SetTrusted assigns a value to `trusted`, otherwise not accessible from outside the package.
func (s *Subscription) SetTrusted(tstd interface{}) {
	s.trusted = tstd
}

// GetTrusted reads value of `trusted`.
func (s *Subscription) GetTrusted() interface{} {
	return s.trusted
}

// GetLastSeen returns lastSeen.
func (s *Subscription) GetLastSeen() *time.Time {
	if s.lastSeenUA != nil {
		return &s.lastSeenUA.When
	}
	return nil
}

// GetUserAgent returns userAgent.
func (s *Subscription) GetUserAgent() string {
	if s.lastSeenUA != nil {
		return s.lastSeenUA.UserAgent
	}
	return ""
}
