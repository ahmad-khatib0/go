package devices

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Devices struct {
	db     *pgxpool.Pool
	utils  *utils.Utils
	cfg    *config.StorePostgresConfig
	shared *shared.Shared
}

type DevicesArgs struct {
	DB     *pgxpool.Pool
	Utils  *utils.Utils
	Cfg    *config.StorePostgresConfig
	Shared *shared.Shared
}

func NewDevices(da DevicesArgs) *Devices {
	return &Devices{db: da.DB, utils: da.Utils, cfg: da.Cfg, shared: da.Shared}
}

// Upsert creates or updates a device record
func (d *Devices) Upsert(uid types.Uid, def *types.DeviceDef) error {
	hash := d.shared.DeviceHasher(def.DeviceID)
	ctx, cancel := d.utils.GetContext(time.Duration(d.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	// Ensure uniqueness of the device ID: delete all records of the device ID
	_, err = tx.Exec(ctx, "DELETE FROM devices WHERE hash = $1 ", hash)
	if err != nil {
		return err
	}

	// Actually add/update DeviceId for the new user

	stmt := `
	  INSERT INTO devices(user_id, hash, device_id, platform, last_seen, lang) 
		VALUES($1, $2, $3, $4, $5, $6)
	`
	_, err = tx.Exec(
		ctx,
		stmt,
		store.DecodeUid(uid),
		hash,
		def.DeviceID,
		def.Platform,
		def.LastSeen,
		def.Lang,
	)

	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// GetAll returns all devices for a given set of users
func (d *Devices) GetAll(uids ...types.Uid) (map[types.Uid][]types.DeviceDef, int, error) {
	var unupg []any
	for _, uid := range uids {
		unupg = append(unupg, store.DecodeUid(uid))
	}

	stmt := `
	  SELECT user_id, device_id, platform, last_seen, lang 
		FROM devices WHERE user_id IN (?)
	`

	query, unupg := d.shared.ExpandQuery(stmt, unupg)
	ctx, cancel := d.utils.GetContext(time.Duration(d.cfg.SqlTimeout))
	if cancel != nil {
		defer cancel()
	}

	rows, err := d.db.Query(ctx, query, unupg...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var device struct {
		Userid   int64
		Deviceid string
		Platform string
		Lastseen time.Time
		Lang     string
	}

	result := make(map[types.Uid][]types.DeviceDef)
	count := 0
	for rows.Next() {
		if err = rows.Scan(
			&device.Userid,
			&device.Deviceid,
			&device.Platform,
			&device.Lastseen,
			&device.Lang,
		); err != nil {
			break
		}

		uid := store.EncodeUid(device.Userid)
		udev := result[uid]
		udev = append(udev, types.DeviceDef{
			DeviceID: device.Deviceid,
			Platform: device.Platform,
			LastSeen: device.Lastseen,
			Lang:     device.Lang,
		})

		result[uid] = udev
		count++
	}

	if err == nil {
		err = rows.Err()
	}

	return result, count, err
}

// Delete deletes a device record
func (d *Devices) Delete(uid types.Uid, deviceID string) error {
	ctx, cancel := d.utils.GetContext(time.Duration(d.cfg.TxTimeout))
	if cancel != nil {
		defer cancel()
	}

	tx, err := d.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	err = d.shared.DeviceDelete(ctx, tx, uid, deviceID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
