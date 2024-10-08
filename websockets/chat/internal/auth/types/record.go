package types

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// Rec is an authentication record.
type Rec struct {
	// User ID.
	Uid types.Uid `json:"uid,omitempty"`
	// Authentication level.
	AuthLevel Level `json:"authlvl,omitempty"`
	// Lifetime of this record.
	Lifetime Duration `json:"lifetime,omitempty"`
	// Bitmap of features. Currently 'validated'/'not validated' only (V and L).
	Features Feature `json:"features,omitempty"`
	// Tags generated by this authentication record.
	Tags []string `json:"tags,omitempty"`
	// User account state received or read by the authenticator.
	State types.ObjState
	// Credential 'method:value' associated with this record.
	Credential string `json:"cred,omitempty"`

	// Authenticator may request the server to create a new account.
	// These are the account parameters which can be used for creating the account.
	DefAcs  *types.DefaultAccess `json:"defacs,omitempty"`
	Public  interface{}          `json:"public,omitempty"`
	Private interface{}          `json:"private,omitempty"`
}
