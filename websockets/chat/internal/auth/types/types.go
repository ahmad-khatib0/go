// Package auth provides interfaces and types required for implementing an authenticaor.
package types

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
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

	// AddRecord adds persistent authentication record to the database.
	// Returns: updated auth record, error
	AddRecord(rec *Rec, secret []byte, remoteAddr string) (*Rec, error)

	// UpdateRecord updates existing record with new credentials.
	// Returns updated auth record, error.
	UpdateRecord(rec *Rec, secret []byte, remoteAddr string) (*Rec, error)

	// Authenticate: given a user-provided authentication secret (such as "login:password"), either
	//
	// return user's record (ID, time when the secret expires, etc), or issue a challenge to
	//
	// continue the authentication process to the next step, or return an error code.
	//
	// The remoteAddr (i.e. the IP address of the client) can be used by custom authenticators for
	//
	// additional validation. The stock authenticators don't use it.
	//
	// store.Users.GetAuthRecord("scheme", "unique")
	//
	// Returns: user auth record, challenge, error.
	Authenticate(secret []byte, remoteAddr string) (*Rec, []byte, error)

	// AsTag converts search token into prefixed tag or an empty string if it
	// cannot be represented as a prefixed tag.
	AsTag(token string) string

	// IsUnique verifies if the provided secret can be considered unique by the auth scheme
	// E.g. if login is unique.
	IsUnique(secret []byte, remoteAddr string) (bool, error)

	// GenSecret generates a new secret, if appropriate.
	GenSecret(rec *Rec) ([]byte, time.Time, error)

	// DelRecords deletes (or disables) all authentication records for the given user.
	DelRecords(uid types.Uid) error

	// GetResetParams returns authenticator parameters passed to password reset handler
	// for the provided user id.
	// Returns: map of params.
	GetResetParams(uid types.Uid) (map[string]interface{}, error)
}
