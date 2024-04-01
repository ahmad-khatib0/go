package config

type WebrtcConfig struct {
	// Enable video/voice calls.
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// Timeout in seconds before a call is dropped if not answered.
	CallEstablishmentTimeout int                     `json:"call_establishment_timeout" mapstructure:"call_establishment_timeout"`
	IceServers               []WebRtcConfigIceServer `json:"ice_servers" mapstructure:"ice_servers"`
	// Alternative config as an external file.
	IceServersFile string `json:"ice_servers_file" mapstructure:"ice_servers_file"`
}

type WebRtcConfigIceServer struct {
	Username       string            `json:"username" mapstructure:"username"`
	Credential     string            `json:"credential" mapstructure:"credential"`
	CredentialType string            `json:"credential_type" mapstructure:"credential_type"`
	Urls           []string          `json:"urls" mapstructure:"urls"`
	Config         WebRtcVideoConfig `json:"config" mapstructure:"config"`
}

type WebRtcVideoConfig struct {
	Enabled     bool   `json:"enabled" mapstructure:"enabled"`
	EndpointUrl string `json:"endpoint_url" mapstructure:"endpoint_url"`
	ApiKey      string `json:"api_key" mapstructure:"api_key"`
	ApiSecret   string `json:"api_secret" mapstructure:"api_secret"`
	MaxDuration int    `json:"max_duration" mapstructure:"max_duration"`
}
