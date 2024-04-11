package types

// StoreError satisfies Error interface but allows constant values for
// direct comparison.
type StoreError string

const (
	// ErrInternal means DB or other internal failure.
	ErrInternal = StoreError("internal")
	// ErrMalformed means the secret cannot be parsed or otherwise wrong.
	ErrMalformed = StoreError("malformed")
	// ErrFailed means authentication failed (wrong login or password, etc).
	ErrFailed = StoreError("failed")
	// ErrDuplicate means duplicate credential, i.e. non-unique login.
	ErrDuplicate = StoreError("duplicate value")
	// ErrUnsupported means an operation is not supported.
	ErrUnsupported = StoreError("unsupported")
	// ErrExpired means the secret has expired.
	ErrExpired = StoreError("expired")
	// ErrPolicy means policy violation, e.g. password too weak.
	ErrPolicy = StoreError("policy")
	// ErrCredentials means credentials like email or captcha must be validated.
	ErrCredentials = StoreError("credentials")
	// ErrUserNotFound means the user was not found.
	ErrUserNotFound = StoreError("user not found")
	// ErrTopicNotFound means the topic was not found.
	ErrTopicNotFound = StoreError("topic not found")
	// ErrNotFound means the object other then user or topic was not found.
	ErrNotFound = StoreError("not found")
	// ErrPermissionDenied means the operation is not permitted.
	ErrPermissionDenied = StoreError("denied")
	// ErrInvalidResponse means the client's response does not match server's expectation.
	ErrInvalidResponse = StoreError("invalid response")
	// ErrRedirected means the subscription request was redirected to another topic.
	ErrRedirected = StoreError("redirected")
)

// Error is required by error interface.
func (s StoreError) Error() string {
	return string(s)
}
