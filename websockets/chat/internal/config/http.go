package config

type HttpConfig struct {
	// Cache-Control value for static content.
	CacheControl int `json:"cache_control" mapstructure:"cache_control"`
	// Take IP address of the client from HTTP header 'X-Forwarded-For'.
	// Useful when chat app is behind a proxy. If missing, fallback to default RemoteAddr.
	UseXForwardedFor bool `json:"use_x_forwarded_for" mapstructure:"use_x_forwarded_for"`
}
