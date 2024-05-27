package server

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	pt "github.com/ahmad-khatib0/go/websockets/chat/internal/push"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"google.golang.org/grpc"
)

var globals struct {
	// Topics cache and processing.
	hub *Hub
	// Indicator that shutdown is in progress
	shuttingDown bool
	// Sessions cache.
	sessionStore *SessionStore
	// Cluster data.
	cluster *Cluster
	// gRPC server.
	grpcServer *grpc.Server
	// Plugins.
	// plugins []Plugin

	// Users cache communication channel.
	// usersUpdate chan *UserCacheReq

	// Credential validator config to pass to clients.
	validatorClientConfig map[string][]string
	// Validators required for each auth level.
	authValidators map[types.Level][]string
	// Credential validators.
	validators map[string]credValidator

	// Salt used for signing API key.
	apiKeySalt []byte
	// Tag namespaces (prefixes) which are immutable to the client.
	immutableTagNS map[string]bool
	// Tag namespaces which are immutable on User and partially mutable on Topic:
	// user can only mutate tags he owns.
	maskedTagNS map[string]bool

	// Add Strict-Transport-Security to headers, the value signifies age.
	// Empty string "" turns it off
	tlsStrictMaxAge string
	// Listen for connections on this address:port and redirect them to HTTPS port.
	tlsRedirectHTTP string
	// Maximum message size allowed from peer.
	maxMessageSize int64
	// Maximum number of group topic subscribers.
	maxSubscriberCount int
	// Maximum number of indexable tags.
	maxTagCount int
	// If true, ordinary users cannot delete their accounts.
	permanentAccounts bool

	// Maximum allowed upload size.
	maxFileUploadSize int64
	// Periodicity of a garbage collector for abandoned media uploads.
	mediaGcPeriod time.Duration

	// Prioritize X-Forwarded-For header as the source of IP address of the client.
	useXForwardedFor bool

	// Country code to assign to sessions by default.
	defaultCountryCode string

	// Time before the call is dropped if not answered.
	callEstablishmentTimeout int

	// ICE servers config (video calling)
	iceServers []config.WebRtcConfigIceServer

	// Websocket per-message compression negotiation is enabled.
	wsCompression bool

	// Users cache communication channel.
	usersUpdate chan *UserCacheReq

	// URL of the main endpoint.
	// TODO: implement file-serving API for gRPC and remove this feature.
	servingAt string

	store   *store.Store
	plugins []Plugin
	push    *pt.Push
	l       *logger.Logger
	stats   *stats.Stats
}

type Server struct{}

type ServerArgs struct {
	log   *logger.Logger
	stats *stats.Stats
}

// CredValidator holds additional config params for a credential validator.
type credValidator struct {
	// AuthLevel(s) which require this validator.
	requiredAuthLvl []types.Level
	addToTags       bool
}

func NewServer(sa ServerArgs) *Server {
	globals.l = sa.log
	globals.stats = sa.stats

	return &Server{}
}
