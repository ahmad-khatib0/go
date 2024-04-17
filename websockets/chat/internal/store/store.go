// Package store provides methods for registering and accessing database adapters.
package store

import (
	"fmt"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// NewStore() configure the selected db adapter (AdapterName), and opens the db connection
// TODO: register the avaiable auth methods from the auth pkg
func NewStore(a StoreArgs) (*Store, error) {
	var adp db.Adapter
	var err error

	if a.WorkerID < 0 || a.WorkerID > 1023 {
		return nil, fmt.Errorf("NewStore: invalid workerID")
	}

	uid, err := types.NewUID(a.WorkerID, []byte(a.Cfg.Store.UidKey))
	if err != nil {
		return nil, fmt.Errorf("NewStore: failed to init uid generator %w", err)
	}

	args := db.AdapterArgs{Conf: a.Cfg, Logger: a.Logger, UGen: uid}
	switch a.Cfg.Store.AdapterName {
	case "postgres":
		adp, err = postgres.NewPostgres(args)
	}

	return &Store{
		adp:    adp,
		cfg:    a.Cfg,
		logger: a.Logger,
	}, nil
}
