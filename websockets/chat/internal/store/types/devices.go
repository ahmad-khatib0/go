package types

import "time"

// DeviceDef is the data provided by connected device. Used primarily for
// push notifications.
type DeviceDef struct {
	// Device registration ID
	DeviceId string
	// Device platform (iOS, Android, Web)
	Platform string
	// Last logged in
	LastSeen time.Time
	// Device language, ISO code
	Lang string
}
