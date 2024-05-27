package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// Delete deletes a subscription
func (s *Store) SubsDelete(topic string, user types.Uid) error {
	return s.adp.Subscriptions().Delete(topic, user)
}

func (s *Store) SubsGetSubs(id types.Uid) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForUser(id)
}
