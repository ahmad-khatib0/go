package subscriptions

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Subscriptions struct {
	db *pgxpool.Pool
}

// Delete implements db.Subscriptions.
func (s *Subscriptions) Delete(topic string, user types.Uid) error {
	panic("unimplemented")
}

// Update implements db.Subscriptions.
func (s *Subscriptions) Update(topic string, user types.Uid, update map[string]any) error {
	panic("unimplemented")
}

type SubscriptionsArgs struct {
	DB *pgxpool.Pool
}

func NewSubscriptions(ua SubscriptionsArgs) *Subscriptions {
	return &Subscriptions{db: ua.DB}
}
