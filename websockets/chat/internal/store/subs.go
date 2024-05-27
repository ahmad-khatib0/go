package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// Delete deletes a subscription
func (s *Store) SubsDelete(topic string, user types.Uid) error {
	return s.adp.Subscriptions().Delete(topic, user)
}

// SubsGetSubs loads *all* subscriptions for the given user.
//
// Does not load Public/Trusted or Private, does not load deleted subscriptions.
func (s *Store) SubsGetSubs(topic string, user types.Uid, keepDeleted bool) (*types.Subscription, error) {
	return s.adp.Subscriptions().Get(topic, user, keepDeleted)
}

// Create creates multiple subscriptions
func (s *Store) SubsCreate(subs ...*types.Subscription) error {
	for _, sub := range subs {
		sub.InitTimes()
	}
	return s.adp.Topics().Share(subs)
}

// Update values of topic's subscriptions.
func (s *Store) SubsUpdate(topic string, user types.Uid, update map[string]interface{}) error {
	update["UpdatedAt"] = types.TimeNow()
	return s.adp.Subscriptions().Update(topic, user, update)
}
