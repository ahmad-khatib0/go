package devices

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Devices struct {
	db *pgxpool.Pool
}

// Delete implements db.Devices.
func (d *Devices) Delete(uid types.Uid, deviceID string) error {
	panic("unimplemented")
}

type DevicesArgs struct {
	DB *pgxpool.Pool
}

func NewDevices(ua DevicesArgs) *Devices {
	return &Devices{db: ua.DB}
}
