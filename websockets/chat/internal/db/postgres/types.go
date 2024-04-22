package postgres

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/credentials"
	idb "github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/db"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/devices"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/files"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/messages"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/persistentcache"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/search"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/subscriptions"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/topics"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/postgres/users"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/db/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	// If DB request timeout is specified,
	// we allocate txTimeoutMultiplier times more time for transactions.
	txTimeoutMultiplier = 1.5
)

type postgres struct {
	auth            *auth.Auth
	credentials     *credentials.Credentials
	dB              *idb.DB
	devices         *devices.Devices
	files           *files.Files
	messages        *messages.Messages
	persistentCache *persistentcache.PersistentCache
	search          *search.Search
	subscriptions   *subscriptions.Subscriptions
	topics          *topics.Topics
	users           *users.Users

	db         *pgxpool.Pool
	logger     *logger.Logger
	poolConfig *pgxpool.Config
	version    int
}

// Auth implements db.Adapter.
func (p *postgres) Auth() types.Auth {
	return p.auth
}

// Credentials implements db.Adapter.
func (p *postgres) Credentials() types.Credentials {
	return p.credentials
}

// DB implements db.Adapter.
func (p *postgres) DB() types.DB {
	return p.dB
}

// Devices implements db.Adapter.
func (p *postgres) Devices() types.Devices {
	return p.devices
}

// Files implements db.Adapter.
func (p *postgres) Files() types.Files {
	return p.files
}

// Messages implements db.Adapter.
func (p *postgres) Messages() types.Messages {
	return p.messages
}

// PersistentCache implements db.Adapter.
func (p *postgres) PersistentCache() types.PersistentCache {
	return p.persistentCache
}

// Search implements db.Adapter.
func (p *postgres) Search() types.Search {
	return p.search
}

// Subscriptions implements db.Adapter.
func (p *postgres) Subscriptions() types.Subscriptions {
	return p.subscriptions
}

// Topics implements db.Adapter.
func (p *postgres) Topics() types.Topics {
	return p.topics
}

// Users implements db.Adapter.
func (p *postgres) Users() types.Users {
	return p.users
}
