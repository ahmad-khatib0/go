package auth

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Auth struct {
	db *pgxpool.Pool
}

// AddRecord implements db.Auth.
func (a *Auth) AddRecord(user types.Uid, scheme string, unique string, authLvl auth.Level, secret []byte, expires time.Time) error {
	panic("unimplemented")
}

// DelAllRecords implements db.Auth.
func (a *Auth) DelAllRecords(uid types.Uid) (int, error) {
	panic("unimplemented")
}

// DelScheme implements db.Auth.
func (a *Auth) DelScheme(user types.Uid, scheme string) error {
	panic("unimplemented")
}

// GetRecord implements db.Auth.
func (a *Auth) GetRecord(user types.Uid, scheme string) (string, auth.Level, []byte, time.Time, error) {
	panic("unimplemented")
}

// GetUniqueRecord implements db.Auth.
func (a *Auth) GetUniqueRecord(unique string) (types.Uid, auth.Level, []byte, time.Time, error) {
	panic("unimplemented")
}

// UpdRecord implements db.Auth.
func (a *Auth) UpdRecord(user types.Uid, scheme string, unique string, authLvl auth.Level, secret []byte, expires time.Time) error {
	panic("unimplemented")
}

type AuthArgs struct {
	DB *pgxpool.Pool
}

func NewAuth(ua AuthArgs) *Auth {
	return &Auth{db: ua.DB}
}
