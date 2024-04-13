package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/db"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/credentials"
	idb "github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/db"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/devices"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/files"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/messages"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/persistentcache"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/search"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/shared"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/subscriptions"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/topics"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/users"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Open opens the db connection and configure the releated fields for the adapter
func (p *postgres) Open(aa db.AdapterArgs) (db.Adapter, error) {
	if p.db != nil {
		return nil, errors.New("postgres db is alread connected")
	}

	c, ok := (aa.Conf).(config.StorePostgresConfig)
	if !ok {
		return nil, errors.New("postgres db failed to cast config to StorePostgresConfig")
	}

	dsn, err := parseConnString(&c)
	if err != nil {
		return nil, err
	}

	if c.MaxResults <= 0 {
		p.maxResults = constants.DBDefaultMaxResults
	}
	if c.MaxMessageResults <= 0 {
		p.maxMessageResults = constants.DBDefaultMaxMessageResults
	}
	if c.MaxOpenConn > 0 {
		p.poolConfig.MaxConns = int32(c.MaxOpenConn)
	}
	if c.MaxIdleConn > 0 {
		p.poolConfig.MinConns = int32(c.MaxIdleConn)
	}
	if c.MaxLifetimeConn > 0 {
		p.poolConfig.MaxConnLifetime = time.Duration(c.MaxLifetimeConn) * time.Second
	}
	if c.SqlTimeout > 0 {
		p.sqlTimeout = time.Duration(c.SqlTimeout) * time.Second
		p.txTimeout = time.Duration(float64(c.SqlTimeout)*txTimeoutMultiplier) * time.Second
	}

	if p.poolConfig, err = pgxpool.ParseConfig(dsn); err != nil {
		return nil, fmt.Errorf("postgres db failed to parse DSN config %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	p.db, err = pgxpool.NewWithConfig(ctx, p.poolConfig)
	if err != nil {
		return nil, fmt.Errorf("postgres db failed to parse DSN config %w", err)
	}

	ut := utils.NewUtils()
	sh := shared.NewShared(shared.SharedArgs{Utils: ut})

	p.dB = idb.NewDB(idb.DBArgs{DB: p.db, Cfg: &c, Utils: ut})
	p.users = users.NewUsers(users.UsersArgs{DB: p.db, Utils: ut, Cfg: &c, Shared: sh})
	p.credentials = credentials.NewCredentials(credentials.CredentialsArgs{DB: p.db, Utils: ut, Cfg: &c, Shared: sh})
	p.auth = auth.NewAuth(auth.AuthArgs{DB: p.db, Cfg: &c, Utils: ut, Shared: sh})
	p.topics = topics.NewTopics(topics.TopicsArgs{DB: p.db, Cfg: &c, Utils: ut, Shared: sh, Logger: aa.Logger})
	p.subscriptions = subscriptions.NewSubscriptions(
		subscriptions.SubscriptionsArgs{DB: p.db, Utils: ut, Cfg: &c, Shared: sh},
	)
	p.devices = devices.NewDevices(devices.DevicesArgs{DB: p.db})
	p.files = files.NewFiles(files.FilesArgs{DB: p.db})
	p.messages = messages.NewMessages(messages.MessagesArgs{DB: p.db})
	p.persistentCache = persistentcache.NewPersistentCache(persistentcache.PersistentCacheArgs{DB: p.db})
	p.search = search.NewSearch(search.SearchArgs{DB: p.db})

	// TODO: check here the missing db
	err = p.db.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("postgres db failed to ping database %w", err)
	}

	return p, nil
}

func parseConnString(c *config.StorePostgresConfig) (string, error) {
	if c.User == "" || c.Password == "" || c.Host == "" || c.Port == 0 || c.DbName == "" {
		return "", errors.New("postgres db invalid config value")
	}

	connStr := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DbName,
		c.SqlTimeout,
	)

	return connStr, nil

}
