package types

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"strings"
	"time"
)

// TopicCat is an enum of topic categories.
type TopicCat int

const (
	// TopicCatMe is a value denoting 'me' topic.
	TopicCatMe TopicCat = iota
	// TopicCatFnd is a value denoting 'fnd' topic.
	TopicCatFnd
	// TopicCatP2P is a value denoting 'p2p topic.
	TopicCatP2P
	// TopicCatGrp is a value denoting group topic.
	TopicCatGrp
	// TopicCatSys is a constant indicating a system topic.
	TopicCatSys
)

// Topic stored in database. Topic's name is Id
type Topic struct {
	ObjHeader `bson:",inline"`

	// State of the topic: normal (ok), suspended, deleted
	State   ObjState
	StateAt *time.Time `json:"StateAt,omitempty" bson:",omitempty"`

	// Timestamp when the last message has passed through the topic
	TouchedAt time.Time

	// Indicates that the topic is a channel.
	UseBt bool

	// Topic owner. Could be zero
	Owner string

	// Default access to topic
	Access DefaultAccess

	// Server-issued sequential ID
	SeqId int
	// If messages were deleted, sequential id of the last operation to delete them
	DelId int

	Public  interface{}
	Trusted interface{}

	// Indexed tags for finding this topic.
	Tags StringSlice

	// Deserialized ephemeral params
	perUser map[Uid]*perUserData // deserialized from Subscription
}

// TopicsParseP2P extracts uids from the name of a p2p topic.
func TopicsParseP2P(p2p string) (uid1, uid2 Uid, err error) {
	if strings.HasPrefix(p2p, "p2p") {
		src := []byte(p2p)[3:] // bytes that comme after the p2p  (p2p:some-uid)

		if len(src) != p2pBase64Unpadded {
			err = errors.New("invalid length <TopicsParseP2P>")
			return
		}

		var count int
		dec := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).DecodedLen(p2pBase64Unpadded))
		count, err = base64.URLEncoding.WithPadding(base64.NoPadding).Decode(dec, src)

		if count < 16 {
			if err != nil {
				err = errors.New("failed to decode <TopicsParseP2P>" + err.Error())
			} else {
				err = errors.New("invalid decoded length <TopicsParseP2P>")
			}
			return
		}

		uid1 = Uid(binary.LittleEndian.Uint64(dec))
		uid2 = Uid(binary.LittleEndian.Uint64(dec[8:]))
	} else {
		err = errors.New("missing or invalid prefix <TopicsParseP2P>")
	}

	return
}

// GetTopicCat given topic name returns topic category.
func GetTopicCat(name string) TopicCat {
	switch name[:3] {
	case "usr":
		return TopicCatMe
	case "p2p":
		return TopicCatP2P
	case "grp", "chn":
		return TopicCatGrp
	case "fnd":
		return TopicCatFnd
	case "sys":
		return TopicCatSys
	default:
		panic("invalid topic type for name '" + name + "'")
	}
}

// TopicsChnToGrp gets group topic name from channel name.
//
// If it's a non-channel group topic, the name is returned unchanged.
func TopicsChnToGrp(chn string) string {
	if strings.HasPrefix(chn, "chn") {
		return strings.Replace(chn, "chn", "grp", 1)
	}

	// Return unchanged if it's a group already.
	if strings.HasPrefix(chn, "grp") {
		return chn
	}
	return ""
}

// TopicsGrpToChn converts group topic name to corresponding channel name.
func TopicsGrpToChn(grp string) string {
	if strings.HasPrefix(grp, "grp") {
		return strings.Replace(grp, "grp", "chn", 1)
	}

	// Return unchanged if it's a channel already.
	if strings.HasPrefix(grp, "chn") {
		return grp
	}
	return ""
}

// IsChannel checks if the given topic name is a reference to a channel.
//
// The "nch" should not be considered a channel reference because it can only
//
// be used by the topic owner at the time of group topic creation.
func IsChannel(name string) bool {
	return strings.HasPrefix(name, "chn")
}

// ChnToGrp gets group topic name from channel name.
// If it's a non-channel group topic, the name is returned unchanged.
func ChnToGrp(chn string) string {
	if strings.HasPrefix(chn, "chn") {
		return strings.Replace(chn, "chn", "grp", 1)
	}
	// Return unchanged if it's a group already.
	if strings.HasPrefix(chn, "grp") {
		return chn
	}
	return ""
}

// P2PNameForUser takes a user ID and a full name of a P2P topic and generates the name of the
// P2P topic as it should be seen by the given user.
func P2PNameForUser(uid Uid, p2p string) (string, error) {
	uid1, uid2, err := ParseP2P(p2p)
	if err != nil {
		return "", err
	}
	if uid.Compare(uid1) == 0 {
		return uid2.UserId(), nil
	}
	return uid1.UserId(), nil
}

// ParseP2P extracts uids from the name of a p2p topic.
func ParseP2P(p2p string) (uid1, uid2 Uid, err error) {
	if strings.HasPrefix(p2p, "p2p") {
		src := []byte(p2p)[3:]
		if len(src) != p2pBase64Unpadded {
			err = errors.New("ParseP2P: invalid length")
			return
		}
		dec := make([]byte, base64.URLEncoding.WithPadding(base64.NoPadding).DecodedLen(p2pBase64Unpadded))
		var count int
		count, err = base64.URLEncoding.WithPadding(base64.NoPadding).Decode(dec, src)
		if count < 16 {
			if err != nil {
				err = errors.New("ParseP2P: failed to decode " + err.Error())
			} else {
				err = errors.New("ParseP2P: invalid decoded length")
			}
			return
		}
		uid1 = Uid(binary.LittleEndian.Uint64(dec))
		uid2 = Uid(binary.LittleEndian.Uint64(dec[8:]))
	} else {
		err = errors.New("ParseP2P: missing or invalid prefix")
	}
	return
}
