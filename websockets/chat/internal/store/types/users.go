package types

import "time"

// User is a representation of a DB-stored user record.
type User struct {
	ObjHeader
	State   ObjState
	StateAt *time.Time

	// Default access to user for P2P topics (used as default modeGiven)
	Access DefaultAccess

	// Values for 'me' topic:

	// Last time when the user joined 'me' topic, by User Agent
	LastSeen *time.Time
	// User agent provided when accessing the topic last time
	UserAgent string

	Public  interface{}
	Trusted interface{}

	// Unique indexed tags (email, phone) for finding this user. Stored on the
	// 'users' as well as indexed in 'tagunique'
	Tags StringSlice

	// Info on known devices, used for push notifications
	Devices map[string]*DeviceDef
	// Same for mongodb scheme. Ignore in other db backends if its not suitable.
	DeviceArray []*DeviceDef
}
