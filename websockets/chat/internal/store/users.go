package store

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

// UpdateLastSeen updates LastSeen and UserAgent.
func (s *Store) UsersUpdateLastSeen(uid types.Uid, userAgent string, when time.Time) error {
	return s.adp.Users().Update(uid, map[string]interface{}{"LastSeen": when, "UserAgent": userAgent})
}

// Create inserts User object into a database, updates creation time and assigns UID
func (s *Store) UsersCreate(user *types.User, private interface{}) (*types.User, error) {
	user.SetUid(s.UidGen.Get())
	user.InitTimes()

	err := s.adp.Users().Create(user)
	if err != nil {
		return nil, err
	}

	// Create user's subscription to 'me' && 'fnd'. These topics are ephemeral,
	// the topic object need not to be inserted.
	err = s.SubsCreate(
		&types.Subscription{
			ObjHeader: types.ObjHeader{CreatedAt: user.CreatedAt},
			User:      user.Id,
			Topic:     user.Uid().UserId(),
			ModeWant:  types.ModeCSelf,
			ModeGiven: types.ModeCSelf,
			Private:   private,
		},
		&types.Subscription{
			ObjHeader: types.ObjHeader{CreatedAt: user.CreatedAt},
			User:      user.Id,
			Topic:     user.Uid().FndName(),
			ModeWant:  types.ModeCSelf,
			ModeGiven: types.ModeCSelf,
			Private:   nil,
		})
	if err != nil {
		// Best effort to delete incomplete user record. Orphaned user records are not a problem.
		// They just take up space.
		s.adp.Users().Delete(user.Uid(), true)
		return nil, err
	}

	return user, nil
}

// GetAll returns a slice of user objects for the given user ids
func (s *Store) UsersGetAll(uid ...types.Uid) ([]types.User, error) {
	return s.adp.Users().GetAll(uid...)
}

// GetByCred returns user ID for the given validated credential.
func (s *Store) UsersGetByCred(method, value string) (types.Uid, error) {
	return s.adp.Users().GetByCred(method, value)
}

// Delete deletes user records.
func (s *Store) UsersDelete(id types.Uid, hard bool) error {
	return s.adp.Users().Delete(id, hard)
}

// Update is a general-purpose update of user data.
func (s *Store) UsersUpdate(uid types.Uid, update map[string]interface{}) error {
	if _, ok := update["UpdatedAt"]; !ok {
		update["UpdatedAt"] = types.TimeNow()
	}
	return s.adp.Users().Update(uid, update)
}

// UpdateTags either adds, removes, or resets tags to the given slices.
func (s *Store) UsersUpdateTags(uid types.Uid, add, remove, reset []string) ([]string, error) {
	return s.adp.Users().UpdateTags(uid, add, remove, reset)
}

// UpdateState changes user's state and state of some topics associated with the user.
func (s *Store) UsersUpdateState(uid types.Uid, state types.ObjState) error {
	update := map[string]interface{}{
		"State":   state,
		"StateAt": types.TimeNow()}
	return s.adp.Users().Update(uid, update)
}

// Get returns a user object for the given user id
func (s *Store) UsersGet(uid types.Uid) (*types.User, error) {
	return s.adp.Users().Get(uid)
}

// GetSubs loads *all* subscriptions for the given user.
// Does not load Public/Trusted or Private, does not load deleted subscriptions.
func (s *Store) UsersGetSubs(id types.Uid) ([]types.Subscription, error) {
	return s.adp.Subscriptions().SubsForUser(id)
}

// FindSubs find a list of users and topics for the given tags. Results are formatted as subscriptions.
//
// `required` specifies an AND of ORs for required terms:
//
// at least one element of every sublist in `required` must be present in the object's tags list.
//
// `optional` specifies a list of optional terms.
func (s *Store) UsersFindSubs(id types.Uid, required [][]string, optional []string, activeOnly bool) ([]types.Subscription, error) {
	usubs, err := s.adp.Search().FindUsers(id, required, optional, activeOnly)
	if err != nil {
		return nil, err
	}

	tsubs, err := s.adp.Search().FindTopics(required, optional, activeOnly)
	if err != nil {
		return nil, err
	}

	allSubs := append(usubs, tsubs...)
	for i := range allSubs {
		// Indicate that the returned access modes are not 'N', but rather undefined.
		allSubs[i].ModeGiven = types.ModeUnset
		allSubs[i].ModeWant = types.ModeUnset
	}

	return allSubs, nil
}

// GetTopics load a list of user's subscriptions with Public+Trusted fields copied to subscription
func (s *Store) UsersGetTopics(id types.Uid, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Topics().TopicsForUser(id, false, opts)
}

// GetTopicsAny load a list of user's subscriptions with Public+Trusted
//
// fields copied to subscription. Deleted topics are returned too.
func (s *Store) UsersGetTopicsAny(id types.Uid, opts *types.QueryOpt) ([]types.Subscription, error) {
	return s.adp.Topics().TopicsForUser(id, true, opts)
}

// GetOwnTopics returns a slice of group topic names where the user is the owner.
func (s *Store) UsersGetOwnTopics(id types.Uid) ([]string, error) {
	return s.adp.Topics().OwnTopics(id)
}

// GetChannels returns a slice of group topic names where the user is a channel reader.
func (s *Store) UsersGetChannels(id types.Uid) ([]string, error) {
	return s.adp.Topics().ChannelsForUser(id)
}

// GetUnreadCount returs users' total count of unread messages in all topics with the R permissions.
func (s *Store) UsersGetUnreadCount(ids ...types.Uid) (map[types.Uid]int, error) {
	return s.adp.Users().UnreadCount(ids...)
}

// GetUnvalidated returns a list of stale user ids which have unvalidated credentials,
//
// their auth levels and a comma-separated list of these credential names.
func (s *Store) UsersGetUnvalidated(lastUpdatedBefore time.Time, limit int) ([]types.Uid, error) {
	return s.adp.Users().GetUnvalidated(lastUpdatedBefore, limit)
}

// GetAllCreds returns credentials of the given user, all or validated only.
func (s *Store) UsersGetAllCreds(id types.Uid, method string, validatedOnly bool) ([]types.Credential, error) {
	return s.adp.Credentials().GetAll(id, method, validatedOnly)
}
