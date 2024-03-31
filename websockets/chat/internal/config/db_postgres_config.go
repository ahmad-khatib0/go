package config

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DBPostgresConfig struct {
	DB         *pgxpool.Pool   `json:"db"`
	PoolConfig *pgxpool.Config `json:"pool_config"`
	Dsn        string          `json:"dsn"`
	DbName     string          `json:"db_name"`
	// Maximum number of records to return
	MaxResults int `json:"max_results"`
	// Maximum number of message records to return
	MaxMessageResults int `json:"max_message_results"`
	Version           int `json:"version"`

	// Single query timeout.
	SqlTimeout time.Duration `json:"sql_timeout"`
	// DB transaction timeout.
	TxTimeout time.Duration `json:"tx_timeout"`
}
