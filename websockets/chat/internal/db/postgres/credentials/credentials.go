package credentials

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Credentials struct {
	db *pgxpool.Pool
}

// Confirm implements db.Credentials.
func (c *Credentials) Confirm(uid types.Uid, method string) error {
	panic("unimplemented")
}

// Del implements db.Credentials.
func (c *Credentials) Del(uid types.Uid, method string, value string) error {
	panic("unimplemented")
}

// Fail implements db.Credentials.
func (c *Credentials) Fail(uid types.Uid, method string) error {
	panic("unimplemented")
}

type CredentialsArgs struct {
	DB *pgxpool.Pool
}

func NewCredentials(ua CredentialsArgs) *Credentials {
	return &Credentials{db: ua.DB}
}
