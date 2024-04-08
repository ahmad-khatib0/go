package auth

// Level is the type for authentication levels.
type Level int

// Authentication levels
const (
	// LevelNone is undefined/not authenticated
	LevelNone Level = iota * 10
	// LevelAnon is anonymous user/light authentication
	LevelAnon
	// LevelAuth is fully authenticated user
	LevelAuth
	// LevelRoot is a superuser (currently unused)
	LevelRoot
)

// AuthHandler is the interface which auth providers must implement.
type AuthHandler interface {
	// Init initializes the handler taking config and logical name as parameters.
	Init(conf interface{}, name string) error

	// IsInitialized returns true if the handler is initialized.
	IsInitialized() bool

	// GetRealName returns the hardcoded name of the authenticator.
	GetRealName() string

	// RestrictedTags returns the tag namespaces (prefixes) which are restricted by this authenticator.
	RestrictedTags() ([]string, error)

	// GetAuthConfig() gets the config for an authenticator
	GetAuthConfig() interface{}
}
