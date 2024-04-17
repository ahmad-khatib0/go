package config

type StoreConfig struct {
	// 16-byte key for XTEA. Used to initialize types.UidGenerator.
	UidKey string `json:"uid_key" mapstructure:"uid_key"`
	// Maximum number of results to return from adapter.
	MaxResults  int                  `json:"max_results" mapstructure:"max_results"`
	AdapterName string               `json:"adapter_name" mapstructure:"adapter_name"`
	Postgres    *StorePostgresConfig `json:"postgres" mapstructure:"postgres"`
}
