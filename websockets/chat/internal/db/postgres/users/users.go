package users

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Users struct {
	db *pgxpool.Pool
}

// Delete implements db.Users.
func (u *Users) Delete(uid types.Uid, hard bool) error {
	panic("unimplemented")
}

// GetByCred implements db.Users.
func (u *Users) GetByCred(method string, value string) (types.Uid, error) {
	panic("unimplemented")
}

// GetUnvalidated implements db.Users.
func (u *Users) GetUnvalidated(lastUpdatedBefore time.Time, limit int) ([]types.Uid, error) {
	panic("unimplemented")
}

// UnreadCount implements db.Users.
func (u *Users) UnreadCount(ids ...types.Uid) (map[types.Uid]int, error) {
	panic("unimplemented")
}

// Update implements db.Users.
func (u *Users) Update(uid types.Uid, update map[string]any) error {
	panic("unimplemented")
}

// UpdateTags implements db.Users.
func (u *Users) UpdateTags(uid types.Uid, add []string, remove []string, reset []string) ([]string, error) {
	panic("unimplemented")
}

type UsersArgs struct {
	DB *pgxpool.Pool
}

func NewUsers(ua UsersArgs) *Users {
	return &Users{db: ua.DB}
}
