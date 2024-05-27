package constants

import "time"

const (
	// idleSessionTimeout defines duration of being idle before terminating a session.
	IdleSessionTimeout = time.Second * 55
	// idleMasterTopicTimeout defines now long to keep master topic alive after the last session detached.
	IdleMasterTopicTimeout = time.Second * 4
	// Same as above but shut down the proxy topic sooner. Otherwise master topic would be kept alive for too long.
	IdleProxyTopicTimeout = time.Second * 2

	// defaultMaxMessageSize is the default maximum message size
	DefaultMaxMessageSize = 1 << 19 // 512K

	// defaultMaxSubscriberCount is the default maximum number of group topic subscribers.
	// Also set in adapter.
	DefaultMaxSubscriberCount = 256

	// defaultMaxTagCount is the default maximum number of indexable tags
	DefaultMaxTagCount = 16

	// minTagLength is the shortest acceptable length of a tag in runes. Shorter tags are discarded.
	MinTagLength = 2
	// maxTagLength is the maximum length of a tag in runes. Longer tags are trimmed.
	MaxTagLength = 96

	// Delay before updating a User Agent
	UaTimerDelay = time.Second * 5

	// maxDeleteCount is the maximum allowed number of messages to delete in one call.
	DefaultMaxDeleteCount = 1024

	// Base URL path for serving the streaming API.
	DefaultApiPath = "/"

	// Mount point where static content is served, http://host-name<defaultStaticMount>
	DefaultStaticMount = "/"

	// Local path to static content
	DefaultStaticPath = "static"

	// Default country code to fall back to if the "default_country_code" field
	// isn't specified in the config.
	DefaultCountryCode = "US"

	// Default timeout to drop an unanswered call, seconds.
	DefaultCallEstablishmentTimeout = 30
)
