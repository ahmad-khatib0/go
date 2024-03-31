package config

type WebrtcConfig struct {
	// Enable video/voice calls.
	Enabled bool `json:"enabled"`
	// Timeout in seconds before a call is dropped if not answered.
	CallEstablishmentTimeout int                     `json:"call_establishment_timeout"`
	IceServers               []WebRTCConfigIceServer `json:"ice_servers"`
	// Alternative config as an external file.
	IceServersFile string `json:"ice_servers_file"`
}

type WebRTCConfigIceServer struct {
	Username       string   `json:"username"`
	Credential     string   `json:"credential"`
	CredentialType string   `json:"credential_type"`
	Urls           []string `json:"urls"`
}
