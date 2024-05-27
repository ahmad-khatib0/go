package server

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// presParams defines parameters for creating a presence notification.
type presParams struct {
	userAgent string
	seqID     int
	delID     int
	delSeq    []MsgDelRange

	// Uid who performed the action
	actor string
	// Subject of the action
	target string
	dWant  string
	dGiven string
}

type presFilters struct {
	// Send messages only to users with this access mode being non-zero.
	filterIn types.AccessMode
	// Exclude users with this access mode being non-zero.
	filterOut types.AccessMode
	// Send messages to the sessions of this single user defined by ID as a string 'usrABC'.
	singleUser string
	// Do not send messages to sessions of this user defined by ID as a string 'usrABC'.
	excludeUser string
}
