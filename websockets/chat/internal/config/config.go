package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App        AppConfig     `json:"app" mapstructure:"app"`
	Http       HttpConfig    `mapstructure:"http"`
	Paths      PathsConfig   `json:"paths" mapstructure:"paths"`
	WsConfig   WSConfig      `json:"ws_config" mapstructure:"ws_config"`
	GrpcConfig GrpcConfig    `json:"grpc_config" mapstructure:"grpc_config"`
	Secrets    SecretsConfig `json:"secrets" mapstructure:"secrets"`
	Media      *MediaConfig  `json:"media" mapstructure:"media"`
	PProf      PprofConf     `json:"pprof" mapstructure:"pprof"`

	// Configs for subsystems
	Cluster   ClusterConfig    `json:"cluster" mapstructure:"cluster"`
	Plugins   []PluginConfig   `json:"plugins" mapstructure:"plugins"`
	Store     StoreConfig      `json:"store" mapstructure:"store"`
	Push      PushConfig       `json:"push" mapstructure:"push"`
	Tls       TlsConfig        `json:"tls" mapstructure:"tls"`
	Auth      AuthConfig       `json:"auth" mapstructure:"auth"`
	Validator ValidatorConfig  `json:"validator" mapstructure:"validator"`
	AccountGC *AccountGCConfig `json:"account_gc" mapstructure:"account_gc"`
	Webrtc    *WebrtcConfig    `json:"webrtc" mapstructure:"webrtc"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to init configurations %w", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to init configurations %w", err)
	}

	return &cfg, nil
}
