package config

type AuthConfig struct {
	Basic *AuthConfigBasic `json:"basic"`
	Token *AuthConfigToken `json:"token"`
	Code  *AuthConfigCode  `json:"code"`
}

type AuthConfigBasic struct {
	// Add 'auth-name:username' to tags making user discoverable by username.
	AddToTags bool `json:"add_to_tags"`
	// The minimum length of a login in unicode runes, i.e. "登录" is length 2, not 6.
	// The maximum length is 32 and it cannot be changed.
	MinLoginLength int `json:"min_login_length"`
	// The minimum length of a password in unicode runes, "пароль" is length 6, not 12.
	// There is no limit on maximum length.
	MinPasswordLength int `json:"min_password_length"`
}

type AuthConfigToken struct {
	// Lifetime of a security token in seconds. 1209600 = 2 weeks.
	ExpireIn int `json:"expire_in"`
	// Serial number of the token. Can be used to invalidate all issued tokens at once.
	SerialNumber int `json:"serial_number"`
	// Secret key (HMAC salt) for signing the tokens Any 32 random bytes base64 encoded.
	Key string `json:"key"`
}

// AuthConfigCode Short code authenticator for resetting passwords.
type AuthConfigCode struct {
	// Lifetime of a security code in seconds. 900 seconds = 15 minutes.
	ExpireIn int `json:"expire_in"`
	// Number of times a user can try to enter the code.
	MaxRetries int `json:"max_retries"`
	// Length of the secret code.
	CodeLength int `json:"code_length"`
}
