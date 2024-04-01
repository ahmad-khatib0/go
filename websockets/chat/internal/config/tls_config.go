package config

type TlsConfig struct {
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// Listen for connections on this address:port and redirect them to HTTPS port.
	HttpRedirect string `json:"http_redirect" mapstructure:"http_redirect"`
	// Enable Strict-Transport-Security by setting max_age > 0
	StrictMaxAge int `json:"strict_max_age" mapstructure:"strict_max_age"`
	// ACME autocert config, e.g. letsencrypt.org
	Autocert *TlsAutocertConfig `json:"autocert" mapstructure:"autocert"`
	// If Autocert is not defined, provide file names of static certificate and key
	CertFile string `json:"cert_file" mapstructure:"cert_file"`
	KeyFile  string `json:"key_file" mapstructure:"key_file"`
}

type TlsAutocertConfig struct {
	// Domains to support by autocert
	Domains []string `json:"domains" mapstructure:"domains"`
	// Name of directory where auto-certificates are cached, e.g. /etc/letsencrypt/live/your-domain-here
	Cache string `json:"cache" mapstructure:"cache"`
	// Contact email for letsencrypt
	Email string `json:"email" mapstructure:"email"`
}
