package config

type StorePostgresConfig struct {
	User            string `json:"user" mapstructure:"user"`
	Password        string `json:"password" mapstructure:"password"`
	Host            string `json:"host" mapstructure:"host"`
	Port            int    `json:"port" mapstructure:"port"`
	DbName          string `json:"db_name" mapstructure:"db_name"`
	Dsn             string `json:"dsn" mapstructure:"dsn"`
	MaxOpenConn     int    `json:"max_open_conn" mapstructure:"max_open_conn"`
	MaxIdleConn     int    `json:"max_idle_conn" mapstructure:"max_idle_conn"`
	MaxLifetimeConn int    `json:"max_lifetime_conn" mapstructure:"max_lifetime_conn"`

	// Maximum number of records to return
	MaxResults int `json:"max_results" mapstructure:"max_results"`
	// Maximum number of message records to return
	MaxMessageResults int `json:"max_message_results" mapstructure:"max_message_results"`
	Version           int `json:"version" mapstructure:"version"`

	// Single query timeout.
	SqlTimeout int `json:"sql_timeout" mapstructure:"sql_timeout"`
	// DB transaction timeout.
	TxTimeout int `json:"tx_timeout" mapstructure:"tx_timeout"`
}
