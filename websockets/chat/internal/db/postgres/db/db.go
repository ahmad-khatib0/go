package db

import (
	"strconv"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	db         *pgxpool.Pool
	cfg        *config.StorePostgresConfig
	adpVersion int
	utils      *utils.Utils
}

type DBArgs struct {
	DB    *pgxpool.Pool
	Cfg   *config.StorePostgresConfig
	Utils *utils.Utils
}

func NewDB(ua DBArgs) *DB {
	return &DB{db: ua.DB, cfg: ua.Cfg, utils: ua.Utils}
}

// IsOpen returns true if connection to database has been established.
//
// It does not check if connection is actually live.
func (d *DB) IsOpen() bool {
	return d.db != nil
}

// GetDbVersion returns current database version.
func (d *DB) GetDbVersion() (int, error) {
	if d.cfg.Version > 0 {
		return d.cfg.Version, nil
	}

	ctx, cancel := d.utils.GetContext(time.Duration(d.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	var vers string
	err := d.db.QueryRow(ctx, "SELECT value FROM kvmeta WHERE key = $1", "version").Scan(&vers)
	if err != nil {
		return -1, err
	}

	d.cfg.Version, _ = strconv.Atoi(vers)
	return d.cfg.Version, nil
}

func (d *DB) Close() error {
	if d.db != nil {
		d.db.Close()
		d.db = nil
		d.cfg.Version = -1
	}
	return nil
}

// GetName returns string that adapter uses to register itself with store.
func (d *DB) GetName() string {
	return "postgres"
}

// SetMaxResults configures how many results can be returned in a single DB call.
func (d *DB) SetMaxResults(val int) error {
	if val <= 0 {
		d.cfg.MaxResults = constants.DBDefaultMaxResults
	} else {
		d.cfg.MaxResults = val
	}
	return nil
}