package models

// AuthLevel is the type for authentication levels.
type AuthLevel int

// Authentication levels
const (
	// LevelNone is undefined/not authenticated
	LevelNone AuthLevel = iota * 10
	// LevelAnon is anonymous user/light authentication
	LevelAnon
	// LevelAuth is fully authenticated user
	LevelAuth
	// LevelRoot is a superuser (currently unused)
	LevelRoot
)
