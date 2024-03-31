package config

type Config struct {
	// Cache-Control value for static content.
	CacheControl int `json:"cache_control"`
	// Take IP address of the client from HTTP header 'X-Forwarded-For'.
	// Useful when chat app is behind a proxy. If missing, fallback to default RemoteAddr.
	UseXForwardedFor bool `json:"use_x_forwarded_for"`
	// 2-letter country code (ISO 3166-1 alpha-2) to assign to sessions by default
	// when the country isn't specified by the client explicitly and
	// it's impossible to infer it.
	DefaultCountryCode string `json:"default_country_code"`

	// Configs for subsystems
	Cluster   ClusterConfig   `json:"cluster"`
	Plugins   []PluginConfig  `json:"plugins"`
	Store     StoreConfig     `json:"store"`
	Push      PushConfig      `json:"push"`
	Tls       TlsConfig       `json:"tls"`
	Auth      AuthConfig      `json:"auth"`
	Validator ValidatorConfig `json:"validator"`
	AccountGC AccountGCConfig `json:"account_gc"`
	Webrtc    WebrtcConfig    `json:"webrtc"`
}

func LoadConfig() {

}
