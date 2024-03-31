package config

type TlsConfig struct {
	Enabled bool `json:"enabled"`
	// Listen for connections on this address:port and redirect them to HTTPS port.
	RedirectHTTP string `json:"http_redirect"`
	// Enable Strict-Transport-Security by setting max_age > 0
	StrictMaxAge int `json:"strict_max_age"`
	// ACME autocert config, e.g. letsencrypt.org
	Autocert *TlsAutocertConfig `json:"autocert"`
	// If Autocert is not defined, provide file names of static certificate and key
	CertFile string `json:"cert_file"`
	KeyFile  string `json:"key_file"`
}

type TlsAutocertConfig struct {
	// Domains to support by autocert
	Domains []string `json:"domains"`
	// Name of directory where auto-certificates are cached, e.g. /etc/letsencrypt/live/your-domain-here
	CertCache string `json:"cache"`
	// Contact email for letsencrypt
	Email string `json:"email"`
}
