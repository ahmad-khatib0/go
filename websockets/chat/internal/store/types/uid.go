package types

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
)

// Uid is a database-specific record id, suitable to be used as a primary key.
type Uid uint64

// ZeroUid is a constant representing uninitialized Uid.
const ZeroUid = 0

// NullValue is a Unicode DEL character which indicated that the value is being deleted.
const NullValue = "\u2421"

// Lengths of various Uid representations.
const (
	uidBase64Unpadded = 11
	p2pBase64Unpadded = 22
)

// IsZero checks if Uid is uninitialized.
func (u Uid) IsZero() bool {
	return u == ZeroUid
}

// UserId converts Uid to string prefixed with 'usr', like usrXXXXX.
func (u Uid) UserId() string {
	return u.PrefixId("usr")
}

func (u Uid) PrefixId(prefix string) string {
	if u.IsZero() {
		return ""
	}
	return prefix + u.String()
}

// String converts Uid to base64 string.
func (u Uid) String() string {
	buf, _ := u.MarshalText()
	return string(buf)
}

// MarshalText converts Uid to string represented as byte slice.
func (u *Uid) MarshalText() ([]byte, error) {
	if *u == ZeroUid {
		return []byte{}, nil
	}

	src := make([]byte, 0)
	dst := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).EncodedLen(8))
	binary.LittleEndian.PutUint64(src, uint64(*u))
	base64.URLEncoding.WithPadding(base64.NoPadding).Encode(dst, src)

	return dst, nil
}

// UnmarshalText reads Uid from string represented as byte slice.
func (u *Uid) UnmarshalText(src []byte) error {
	if len(src) != uidBase64Unpadded {
		return errors.New("Uid.UnmarshalText: invalid length")
	}

	dec := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).DecodedLen(uidBase64Unpadded))
	count, err := base64.URLEncoding.WithPadding(base64.NoPadding).Decode(dec, src)
	if count < 8 {
		if err != nil {
			return errors.New("Uid.UnmarshalText: failed to decode " + err.Error())
		}
		return errors.New("Uid.UnmarshalText: failed to decode")
	}

	*u = Uid(binary.LittleEndian.Uint64(dec))
	return nil
}

// MarshalBinary converts Uid to byte slice.
func (uid Uid) MarshalBinary() ([]byte, error) {
	dst := make([]byte, 8)
	binary.LittleEndian.PutUint64(dst, uint64(uid))
	return dst, nil
}

// P2PName takes two Uids and generates a P2P topic name.
func (uid Uid) P2PName(u2 Uid) string {
	if !uid.IsZero() && !u2.IsZero() {
		b1, _ := uid.MarshalBinary()
		b2, _ := u2.MarshalBinary()

		if uid < u2 {
			b1 = append(b1, b2...)
		} else if uid > u2 {
			b1 = append(b2, b1...)
		} else {
			return "" // Explicitly disable P2P with self
		}

		return "p2p" + base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(b1)
	}
	return ""
}

// ParseUid parses string NOT prefixed with anything.
func ParseUid(s string) Uid {
	var uid Uid
	uid.UnmarshalText([]byte(s))
	return uid
}