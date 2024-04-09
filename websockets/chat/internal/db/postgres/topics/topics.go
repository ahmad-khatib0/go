package topics

import (
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Topics struct {
	db *pgxpool.Pool
}

// ChannelsForUser implements db.Topics.
func (t *Topics) ChannelsForUser(uid types.Uid) ([]string, error) {
	panic("unimplemented")
}

// Delete implements db.Topics.
func (t *Topics) Delete(topic string, isChan bool, hard bool) error {
	panic("unimplemented")
}

// OwnerChange implements db.Topics.
func (t *Topics) OwnerChange(topic string, newOwner types.Uid) error {
	panic("unimplemented")
}

// Owns implements db.Topics.
func (t *Topics) Owns(uid types.Uid) ([]string, error) {
	panic("unimplemented")
}

// Update implements db.Topics.
func (t *Topics) Update(topic string, update map[string]any) error {
	panic("unimplemented")
}

type TopicsArgs struct {
	DB *pgxpool.Pool
}

func NewTopics(ua TopicsArgs) *Topics {
	return &Topics{db: ua.DB}
}
