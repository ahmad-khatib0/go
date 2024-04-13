// Package db contains the interfaces to be implemented by the database adapter
package db

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth"
	t "github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
)

type AdapterArgs struct {
	Conf   any
	Logger *logger.Logger
}

type DB interface {
	// Close closes the underlying database connection
	Close() error
	// IsOpen checks if the adapter is ready for use
	IsOpen() bool
	// GetDbVersion returns current database version.
	GetDbVersion() (int, error)
	// CheckDbVersion checks if the actual database version matches adapter version.
	CheckDbVersion() error
	// GetName returns the name of the adapter
	GetName() string
	// SetMaxResults configures how many results can be returned in a single DB call.
	SetMaxResults(val int) error
	// CreateDb creates the database optionally dropping an existing database first.
	CreateDb(reset bool) error
	// UpgradeDb upgrades database to the current adapter version.
	UpgradeDb() error
	// Version returns adapter version
	Version() int
	// DB connection stats object.
	Stats() any
}

type Users interface {
	// UserCreate creates user record
	// Create(user *t.User) error
	// Get returns record for a given user ID
	// Get(uid t.Uid) (*t.User, error)
	// GetAll returns user records for a given list of user IDs
	// GetAll(ids ...t.Uid) ([]t.User, error)
	// Delete deletes user record
	Delete(uid t.Uid, hard bool) error
	// Update updates user record
	Update(uid t.Uid, update map[string]any) error
	// UpdateTags adds, removes, or resets user's tags
	UpdateTags(uid t.Uid, add, remove, reset []string) ([]string, error)
	// GetByCred returns user ID for the given validated credential.
	GetByCred(method, value string) (t.Uid, error)
	// UnreadCount returns the total number of unread messages in all topics with
	// the R permission. If read fails, the counts are still returned with the original
	// user IDs but with the unread count undefined and non-nil error.
	UnreadCount(ids ...t.Uid) (map[t.Uid]int, error)
	// GetUnvalidated returns a list of no more than 'limit' uids who never logged in,
	// have no validated credentials and which haven't been updated since 'lastUpdatedBefore'.
	GetUnvalidated(lastUpdatedBefore time.Time, limit int) ([]t.Uid, error)
}

// Files upload records. The files are stored outside of the database.
type Files interface {
	// FileDeleteUnused deletes records where UseCount is zero. If olderThan is non-zero, deletes
	// unused records with UpdatedAt before olderThan.
	// Returns array of FileDef.Location of deleted filerecords so actual files can be deleted too.
	DeleteUnused(olderThan time.Time, limit int) error
	// StartUpload initializes a file upload.
	// StartUpload(fd *t.FileDef) error
	// FinishUpload marks file upload as completed, successfully or otherwise.
	// FinishUpload(fd *t.FileDef, success bool, size int64) (*t.FileDef, error)
	// Get fetches a record of a specific file
	// Get(fid string) (*t.FileDef, error)
	// LinkAttachments connects given topic or message to the file record IDs from the list.
	LinkAttachments(topic string, userId, msgId t.Uid, fids []string) error
}

type Credentials interface {
	// Upsert adds or updates a credential record. Returns true if record was inserted, false if updated.
	Upsert(cred *t.Credential) (bool, error)
	// GetActive returns the currently active credential record for the given method.
	GetActive(uid t.Uid, method string) (*t.Credential, error)
	// GetAll returns credential records for the given user and method, validated only or all.
	GetAll(uid t.Uid, method string, validatedOnly bool) ([]t.Credential, error)
	// Del deletes credentials for the given method/value. If method is empty, deletes all user's credentials.
	Del(uid t.Uid, method, value string) error
	// Confirm marks given credential as validated.
	Confirm(uid t.Uid, method string) error
	// Fail increments count of failed validation attepmts for the given credentials.
	Fail(uid t.Uid, method string) error
}

// Auth management for the basic authentication scheme
type Auth interface {
	// GetUniqueRecord returns user_id, auth level, secret, expire for a given unique value i.e. login.
	GetUniqueRecord(unique string) (t.Uid, auth.Level, []byte, time.Time, error)
	// GetRecord returns authentication record given user ID and method.
	GetRecord(user t.Uid, scheme string) (string, auth.Level, []byte, time.Time, error)
	// AddRecord creates new authentication record
	AddRecord(user t.Uid, scheme, unique string, authLvl auth.Level, secret []byte, expires time.Time) error
	// DelScheme deletes an existing authentication scheme for the user.
	DelScheme(user t.Uid, scheme string) error
	// DelAllRecords deletes all records of a given user.
	DelAllRecords(uid t.Uid) (int, error)
	// UpdRecord modifies an authentication record. Only non-default/non-zero values are updated.
	UpdRecord(user t.Uid, scheme, unique string, authLvl auth.Level, secret []byte, expires time.Time) error
}

type Topics interface {
	// Create creates a topic
	Create(topic *t.Topic) error
	// CreateP2P creates a p2p topic
	CreateP2P(initiator, invited *t.Subscription) error
	// Get loads a single topic by name, if it exists. If the topic does not exist the call returns (nil, nil)
	Get(topic string) (*t.Topic, error)
	// TopicsForUser() loads subscriptions for a given user. Reads public value.
	//
	// When the 'opts.IfModifiedSince' query is not nil the subscriptions with UpdatedAt > opts.IfModifiedSince
	//
	// are returned, where UpdatedAt can be either a subscription, a topic, or a user update timestamp.
	//
	// This is needed in order to support paginagion of subscriptions: get subscriptions page by page
	//
	// from the oldest updates to most recent:
	//
	// 1. Client already has subscriptions with the latest update timestamp X.
	//
	// 2. Client asks for N updated subscriptions since X. The server returns N with updates between X and Y.
	//
	// 3. Client goes to step 1 with X := Y.
	TopicsForUser(uid t.Uid, keepDeleted bool, opts *t.QueryOpt) ([]t.Subscription, error)
	// UsersForTopic loads users' subscriptions for a given topic. Public is loaded.
	//
	// The difference between UsersForTopic vs SubsForTopic is that the former loads user.Public,
	//
	// but the latter does not.
	UsersForTopic(topic string, keepDeleted bool, opts *t.QueryOpt) ([]t.Subscription, error)
	// OwnTopics() loads a slice of topic names where the user is the owner.
	OwnTopics(uid t.Uid) ([]string, error)
	// ChannelsForUser loads a slice of topic names where the user is a channel reader and notifications (P) are enabled.
	ChannelsForUser(uid t.Uid) ([]string, error)
	// Share creates topic subscriptions
	Share(subs []*t.Subscription) error
	// Delete deletes topic, subscription, messages
	Delete(topic string, isChan, hard bool) error
	// Update updates topic record.
	Update(topic string, update map[string]any) error
	// UpdateOnMessage increments Topic's or User's SeqId value and updates TouchedAt timestamp.
	UpdateOnMessage(topic string, msg *t.Message) error
	// UpdateTopicOwner updates topic's owner
	UpdateTopicOwner(topic string, newOwner t.Uid) error
}

type Subscriptions interface {
	// Get reads a subscription of a user to a topic
	// Get(topic string, user t.Uid, keepDeleted bool) (*t.Subscription, error)
	// SubsForUser loads all subscriptions of a given user. Does NOT load Public or Private values,
	// does not load deleted subscriptions.
	// SubsForUser(user t.Uid) ([]t.Subscription, error)
	// SubsForTopic gets a list of subscriptions to a given topic.. Does NOT load Public value.
	// SubsForTopic(topic string, keepDeleted bool, opts *t.QueryOpt) ([]t.Subscription, error)
	// Update updates parts of a subscription object. Pass nil for fields which don't need to be updated
	Update(topic string, user t.Uid, update map[string]any) error
	// Delete deletes a single subscription
	Delete(topic string, user t.Uid) error
}

type Search interface {
	// FindUsers searches for new contacts given a list of tags.
	// FindUsers(user t.Uid, req [][]string, opt []string, activeOnly bool) ([]t.Subscription, error)
	// FindTopics searches for group topics given a list of tags.
	// FindTopics(req [][]string, opt []string, activeOnly bool) ([]t.Subscription, error)
}

type Messages interface {
	// Save saves message to database
	// Save(msg *t.Message) error
	// GetAll returns messages matching the query
	// GetAll(topic string, forUser t.Uid, opts *t.QueryOpt) ([]t.Message, error)
	// DeleteList marks messages as deleted.
	// Soft- or Hard- is defined by forUser value: forUSer.IsZero == true is hard.
	// DeleteList(topic string, toDel *t.DelMessage) error
	// GetDeleted returns a list of deleted message Ids.
	// GetDeleted(topic string, forUser t.Uid, opts *t.QueryOpt) ([]t.DelMessage, error)
}

// Devices (for push notifications)
type Devices interface {
	// Upsert creates or updates a device record
	// Upsert(uid t.Uid, dev *t.DeviceDef) error
	// GetAll returns all devices for a given set of users
	// GetAll(uid ...t.Uid) (map[t.Uid][]t.DeviceDef, int, error)
	// Delete deletes a device record
	Delete(uid t.Uid, deviceID string) error
}

// Persistent cache management.
type PersistentCache interface {
	// Get reads a persistent cache entry.
	Get(key string) (string, error)
	// Upsert creates or updates a persistent cache entry.
	Upsert(key string, value string, failOnDuplicate bool) error
	// Delete deletes a single persistent cache entry.
	Delete(key string) error
	// Expire expires older entries with the specified key prefix.
	Expire(keyPrefix string, olderThan time.Time) error
}

type Adapter interface {
	// Open opens the db connection and configure the releated fields for the adapter
	Open(aa AdapterArgs) (Adapter, error)
	DB() DB
	Users() Users
	Files() Files
	Credentials() Credentials
	Auth() Auth
	Topics() Topics
	Subscriptions() Subscriptions
	Search() Search
	Messages() Messages
	Devices() Devices
	PersistentCache() PersistentCache
}
