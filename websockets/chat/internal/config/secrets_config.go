package config

type SecretsConfig struct {
	// Salt used in signing API keys
	ApiKeySalt []byte `json:"api_key_salt"`
}
