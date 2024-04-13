package types

import (
	"errors"
	"time"
)

// Various access mode constants.

// AccessMode is a definition of access mode bits.
type AccessMode uint

const (
	ModeJoin    AccessMode = 1 << iota // user can join, i.e. {sub} (J:1)
	ModeRead                           // user can receive broadcasts ({data}, {info}) (R:2)
	ModeWrite                          // user can Write, i.e. {pub} (W:4)
	ModePres                           // user can receive presence updates (P:8)
	ModeApprove                        // user can approve new members or evict existing members (A:0x10, 16)
	ModeShare                          // user can invite new members (S:0x20, 32)
	ModeDelete                         // user can hard-delete messages (D:0x40, 64)
	ModeOwner                          // user is the owner (O:0x80, 128) - full access
	ModeUnset                          // Non-zero value to indicate unknown or undefined mode (:0x100, 256), to make it different from ModeNone.

	// No access, requests to gain access are processed normally (N:0)
	ModeNone AccessMode = 0
	// Normal user's access to a topic ("JRWPS", 47, 0x2F).
	ModeCPublic AccessMode = ModeJoin | ModeRead | ModeWrite | ModePres | ModeShare
	// User's subscription to 'me' and 'fnd' ("JPS", 41, 0x29).
	ModeCSelf AccessMode = ModeJoin | ModePres | ModeShare
	// Owner's subscription to a generic topic ("JRWPASDO", 255, 0xFF).
	ModeCFull AccessMode = ModeJoin | ModeRead | ModeWrite | ModePres | ModeApprove | ModeShare | ModeDelete | ModeOwner
	// Default P2P access mode ("JRWPA", 31, 0x1F).
	ModeCP2P AccessMode = ModeJoin | ModeRead | ModeWrite | ModePres | ModeApprove
	// Default Auth access mode for a user ("JRWPAS", 63, 0x3F).
	ModeCAuth AccessMode = ModeCP2P | ModeCPublic
	// Read-only access to topic ("JR", 3).
	ModeCReadOnly = ModeJoin | ModeRead
	// Access to 'sys' topic by a root user ("JRWPD", 79, 0x4F).
	ModeCSys = ModeJoin | ModeRead | ModeWrite | ModePres | ModeDelete
	// Channel publisher: person authorized to publish content; no J: by invitation only ("RWPD", 78, 0x4E).
	ModeCChnWriter = ModeRead | ModeWrite | ModePres | ModeShare
	// Reader's access mode to a channel (JRP, 11, 0xB).
	ModeCChnReader = ModeJoin | ModeRead | ModePres

	// Admin: user who can modify access mode ("OA", dec: 144, hex: 0x90).
	ModeCAdmin = ModeOwner | ModeApprove
	// Sharer: flags which define user who can be notified of access mode changes ("OAS", dec: 176, hex: 0xB0).
	ModeCSharer = ModeCAdmin | ModeShare

	// Invalid mode to indicate an error.
	ModeInvalid AccessMode = 0x100000

	// All possible valid bits (excluding ModeInvalid and ModeUnset) = 0xFF, 255.
	ModeBitmask AccessMode = ModeJoin | ModeRead | ModeWrite | ModePres | ModeApprove | ModeShare | ModeDelete | ModeOwner
)

// DefaultAccess is a per-topic default access modes
type DefaultAccess struct {
	Auth AccessMode
	Anon AccessMode
}

// String returns string representation of AccessMode.
func (m AccessMode) String() string {
	res, err := m.MarshalText()
	if err != nil {
		return ""
	}
	return string(res)
}

// MarshalJSON converts AccessMode to a quoted string.
func (m AccessMode) MarshalJSON() ([]byte, error) {
	res, err := m.MarshalText()
	if err != nil {
		return nil, err
	}

	return append(append([]byte{'"'}, res...), '"'), nil
}

// UnmarshalJSON reads AccessMode from a quoted string.
func (m *AccessMode) UnmarshalJSON(b []byte) error {
	if b[0] != '"' || b[len(b)-1] != '"' {
		return errors.New("syntax error")
	}

	return m.UnmarshalText(b[1 : len(b)-1])
}

// MarshalText converts AccessMode to ASCII byte slice.
func (m AccessMode) MarshalText() ([]byte, error) {
	if m == ModeNone {
		return []byte{'N'}, nil
	}

	if m == ModeInvalid {
		return nil, errors.New("AccessMode invalid")
	}

	res := []byte{}
	modes := []byte{'J', 'R', 'W', 'P', 'A', 'S', 'D', 'O'}

	for i, chr := range modes {
		if (m & (1 << uint(i))) != 0 {
			res = append(res, chr)
		}
	}
	return res, nil
}

// UnmarshalText parses access mode string as byte slice.
// Does not change the mode if the string is empty or invalid.
func (m *AccessMode) UnmarshalText(b []byte) error {
	m0, err := ParseAcs(b)
	if err != nil {
		return err
	}

	if m0 != ModeUnset {
		*m = (m0 & ModeBitmask)
	}
	return nil
}

// Scan is an implementation of sql.Scanner interface. It expects the
//
// value to be a byte slice representation of an ASCII string.
func (m *AccessMode) Scan(val interface{}) error {
	if bb, ok := val.([]byte); ok {
		return m.UnmarshalText(bb)
	}

	return errors.New("AccessMode: failed to scan data, it is not a byte slice")
}

// IsOwner checks if owner bit O is set.
func (m AccessMode) IsOwner() bool {
	return m&ModeOwner != 0
}

// LastModified returns the greater of either TouchedAt or UpdatedAt.
func (s *Subscription) LastModified() time.Time {
	if s.UpdatedAt.Before(s.touchedAt) {
		return s.touchedAt
	}
	return s.UpdatedAt
}

// SetDefaultAccess updates default access values.
func (s *Subscription) SetDefaultAccess(auth, anon AccessMode) {
	s.modeDefault = &DefaultAccess{auth, anon}
}

// SetLastSeenAndUA updates lastSeen time and userAgent.
func (s *Subscription) SetLastSeenAndUA(when *time.Time, ua string) {
	if when != nil && !when.IsZero() {
		s.lastSeenUA = &LastSeenUA{
			When:      *when,
			UserAgent: ua,
		}
	} else {
		s.lastSeenUA = nil
	}
}

// ParseAcs parses AccessMode from a byte array.
func ParseAcs(b []byte) (AccessMode, error) {
	m0 := ModeUnset

Loop:
	for i := 0; i < len(b); i++ {
		switch b[i] {
		case 'J', 'j':
			m0 |= ModeJoin
		case 'R', 'r':
			m0 |= ModeRead
		case 'W', 'w':
			m0 |= ModeWrite
		case 'A', 'a':
			m0 |= ModeApprove
		case 'S', 's':
			m0 |= ModeShare
		case 'D', 'd':
			m0 |= ModeDelete
		case 'P', 'p':
			m0 |= ModePres
		case 'O', 'o':
			m0 |= ModeOwner
		case 'N', 'n':
			if m0 != ModeUnset {
				return ModeUnset, errors.New("AccessMode: access N cannot be combined with any other")
			}
			m0 = ModeNone // N means explicitly no access, all bits cleared
			break Loop
		default:
			return ModeUnset, errors.New("AccessMode: invalid character '" + string(b[i]) + "'")
		}
	}
	return m0, nil
}
