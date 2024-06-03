package config

type AppConfig struct {
	PermanentAccount  bool   `json:"permanent_account" mapstructure:"permanent_account"`
	Version           string `json:"version" mapstructure:"version"`
	BuildStampCommand string `json:"buildstamp_command" mapstructure:"buildstamp_command"`
	// 2-letter country code (ISO 3166-1 alpha-2) to assign to sessions by default
	// when the country isn't specified by the client explicitly and
	// it's impossible to infer it.
	DefaultCountryCode string `json:"default_country_code" mapstructure:"default_country_code"`
	// Masked tags: tags immutable on User (mask), mutable on Topic only within the mask.
	MaskedTagsNS []string `json:"masked_tags_ns" mapstructure:"masked_tags_ns"`
}
