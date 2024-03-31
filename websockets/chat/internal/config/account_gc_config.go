package config

// Stale unvalidated user account GC config.
type AccountGCConfig struct {
	Enabled         bool `json:"enabled"`
	GcPeriod        int  `json:"gc_period"`          // How often to run GC (seconds).
	GcBlockSize     int  `json:"gc_block_size"`      // Number of accounts to delete in one pass.
	GcMinAccountAge int  `json:"gc_min_account_age"` // Minimum hours since account was last modified.
}
