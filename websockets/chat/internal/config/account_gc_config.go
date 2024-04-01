package config

// Stale unvalidated user account GC config.
type AccountGCConfig struct {
	Enabled bool `json:"enabled" mapstructure:"enabled"`
	// How often to run GC (seconds).
	GcPeriod int `json:"gc_period" mapstructure:"gc_period"`
	// Number of accounts to delete in one pass.
	GcBlockSize int `json:"gc_block_size" mapstructure:"gc_block_size"`
	// Minimum hours since account was last modified.
	GcMinAccountAge int `json:"gc_min_account_age" mapstructure:"gc_min_account_age"`
}
