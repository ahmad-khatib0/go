package config

type SecretsConfig struct {
	// Salt used in signing API keys
	ApiKeySalt string `json:"api_key_salt" mapstructure:"api_key_salt"`
}
