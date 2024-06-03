package server

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat/internal/apikey"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/config"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	pt "github.com/ahmad-khatib0/go/websockets/chat/internal/push"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/stats"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/logger"
	"github.com/ahmad-khatib0/go/websockets/chat/pkg/utils"
	"google.golang.org/grpc"
)

// Delay before updating a User Agent
const uaTimerDelay = time.Second * 5

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

	// Credential validators.
	validators map[string]CredValidator
	// Credential validator config to pass to clients.
	validatorClientConfig map[string][]string
	// Validators required for each auth level.
	authValidators map[types.Level][]string

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

	currentVersion string
	buildstamp     string
	plugins        []Plugin
	store          *store.Store
	utils          *utils.Utils
	push           *pt.Push
	l              *logger.Logger
	stats          *stats.Stats
	apiKey         *apikey.ApiKey
}

// CredValidator holds additional config params for a credential validator.
type CredValidator struct {
	// AuthLevel(s) which require this validator.
	RequiredAuthLvl []types.Level
	AddToTags       bool
}

type ServerArgs struct {
	Cfg             *config.Config
	Log             *logger.Logger
	Stats           *stats.Stats
	Utils           *utils.Utils
	validators      map[string]CredValidator
	validatorCliCfg map[string][]string
	authValidators  map[types.Level][]string
	ImmutableTagNS  map[string]bool
	MaskedTagNS     map[string]bool
}

// Init() inits global config and logger, utils and ...
func Init(sa ServerArgs) {
	c := sa.Cfg
	globals.l = sa.Log
	globals.stats = sa.Stats
	globals.utils = sa.Utils
	globals.apiKeySalt = []byte(c.Secrets.ApiKeySalt)
	globals.validators = sa.validators
	globals.validatorClientConfig = sa.validatorCliCfg
	globals.immutableTagNS = sa.ImmutableTagNS
	globals.maskedTagNS = sa.MaskedTagNS
	globals.maxMessageSize = int64(c.WsConfig.MaxMessageSize)
	globals.maxSubscriberCount = c.WsConfig.MaxSubscriberCount
	globals.maxTagCount = c.WsConfig.MaxTagCount
	globals.permanentAccounts = c.App.PermanentAccount
	globals.useXForwardedFor = c.Http.UseXForwardedFor
	globals.defaultCountryCode = c.App.DefaultCountryCode
	globals.wsCompression = c.WsConfig.WSCompressionEnabled

	if globals.maxMessageSize <= 0 {
		globals.maxMessageSize = constants.DefaultMaxMessageSize
	}
	if globals.maxSubscriberCount <= 0 {
		globals.maxSubscriberCount = constants.DefaultMaxSubscriberCount
	}
	if globals.maxTagCount <= 0 {
		globals.maxTagCount = constants.DefaultMaxTagCount
	}
	if globals.maxTagCount <= 0 {
		globals.maxTagCount = constants.DefaultMaxTagCount
	}
	if globals.defaultCountryCode == "" {
		globals.defaultCountryCode = constants.DefaultCountryCode
	}

}

type Server struct{}

func NewServer(sa ServerArgs) *Server {
	globals.l = sa.Log
	globals.stats = sa.Stats

	return &Server{}
}
