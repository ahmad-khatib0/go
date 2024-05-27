package store

import "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"

// Create creates a topic and owner's subscription to it.
func (s *Store) TopCreate(topic *types.Topic, owner types.Uid, private interface{}) error {

	topic.InitTimes()
	topic.TouchedAt = topic.CreatedAt
	topic.Owner = owner.String()

	err := s.adp.Topics().Create(topic)
	if err != nil {
		return err
	}

	if !owner.IsZero() {
		err = s.SubsCreate(&types.Subscription{
			ObjHeader: types.ObjHeader{CreatedAt: topic.CreatedAt},
			User:      owner.String(),
			Topic:     topic.Id,
			ModeGiven: types.ModeCFull,
			ModeWant:  topic.GetAccess(owner),
			Private:   private})
	}

	return err
}

// CreateP2P creates a P2P topic by generating two user's subsciptions to each other.
func (s *Store) TopCreateP2P(initiator, invited *types.Subscription) error {
	initiator.InitTimes()
	initiator.SetTouchedAt(initiator.CreatedAt)
	invited.InitTimes()
	invited.SetTouchedAt(invited.CreatedAt)

	return s.adp.Topics().CreateP2P(initiator, invited)
}

// Get a single topic with a list of relevant users de-normalized into it
func (s *Store) TopGet(topic string) (*types.Topic, error) {
	return s.adp.Topics().Get(topic)
}

// GetUsers loads subscriptions for topic plus loads user.Public+Trusted.
// Deleted subscriptions are not loaded.
func (s *Store) TopGetUsers(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Topics().UsersForTopic(topic, false, opts)
}

// GetUsersAny loads subscriptions for topic plus loads user.Public+Trusted.
//
// It's the same as GetUsers, except it loads deleted subscriptions too.
func (s *Store) TopGetUsersAny(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Topics().UsersForTopic(topic, true, opts)
}

// Update is a generic topic update.
func (s *Store) TopUpdate(topic string, update map[string]interface{}) error {
	if _, ok := update["UpdatedAt"]; !ok {
		update["UpdatedAt"] = types.TimeNow()
	}
	return s.adp.Topics().Update(topic, update)
}

// OwnerChange replaces the old topic owner with the new owner.
func (s *Store) TopOwnerChange(topic string, newOwner types.Uid) error {
	return s.adp.Topics().UpdateTopicOwner(topic, newOwner)
}

// Delete deletes topic, messages, attachments, and subscriptions.
func (s *Store) TopDelete(topic string, isChan, hard bool) error {
	return s.adp.Topics().Delete(topic, isChan, hard)
}

// TopGetTopicSubs loads a list of subscriptions to the given topic, user.Public+Trusted
//
// and deleted subscriptions are not loaded. Suspended subscriptions are loaded.
func (s *Store) TopGetTopicSubs(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForTopic(topic, false, opts)
}

// TopGetTopicSubsAny loads a list of subscriptions to the given topic including
//
// deleted subscription. user.Public/Trusted are not loaded
func (s *Store) TopGetTopicSubsAny(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForTopic(topic, true, opts)
}

// GetSubs loads a list of subscriptions to the given topic, user.Public+Trusted and deleted
// subscriptions are not loaded. Suspended subscriptions are loaded.
func (s *Store) TopGetSubs(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForTopic(topic, false, opts)
}

// GetSubsAny loads a list of subscriptions to the given topic including deleted subscription.
// user.Public/Trusted are not loaded
func (s *Store) TopGetSubsAny(topic string, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForTopic(topic, true, opts)
}
