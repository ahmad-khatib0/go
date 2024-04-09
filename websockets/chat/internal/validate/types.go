package validate

// Validator handles validation of user's credentials, like email or phone.
type Validator interface {
	// Init initializes the validator.
	Init(valName string, jsonconf any) error
}
