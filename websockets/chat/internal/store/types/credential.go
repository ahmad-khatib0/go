package types

// Credential hold data needed to validate and check validity of a credential like email or phone.
type Credential struct {
	ObjHeader `bson:",inline"`
	// Credential owner
	User string
	// Verification method (email, tel, captcha, etc)
	Method string
	// Credential value - `jdoe@example.com` or `+12345678901`
	Value string
	// Expected response
	Resp string
	// If credential was successfully confirmed
	Done bool
	// Retry count
	Retries int
}
