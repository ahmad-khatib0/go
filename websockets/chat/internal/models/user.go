package models

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// UserCacheReq contains data which mutates one or more user cache entries.
type UserCacheReq struct {
	// Name of the node sending this request in case of cluster. Not set otherwise.
	Node string

	// UserID is set when count of unread messages is updated for a single user or
	// when the user is being deleted.
	UserID types.Uid

	// UserIdList  is set when subscription count is updated for users of a topic.
	UserIdList []types.Uid

	// Unread count (UserId is set)
	Unread int

	// In case of set UserId: treat Unread count as an increment as opposite to the final value.
	// In case of set UserIdList: intement (Inc == true) or decrement subscription count by one.
	Inc bool
	// User is being deleted, remove user from cache.
	Gone bool

	// Optional push notification
	// PushRcpt *push.Receipt
}

type UserCacheEntry struct {
	Unread int
	Topics int
}

// Preserved update entry kept while we read the unread counter from the DB.
type BufferedUpdate struct {
	Val int
	Inc bool
}

type IoResult struct {
	Counts map[types.Uid]int
	Err    error
}

// Represents pending push notification receipt.
type PendingReceipt struct {
	// Number of unread counters currently being read from the DB.
	PendingIOs int
	// The index is needed by update and is maintained by the heap.Interface methods.
	Index int
	// Underlying receipt.
	// rcpt *push.Receipt
}

// Pending pushes organized as a priority queue (priority = number of pending IOs).
// It allows to quickly discover receipts ready for sending (num pending IOs is 0).
type PendingReceiptsQueue []*PendingReceipt
