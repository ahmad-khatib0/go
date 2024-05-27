package server

import (
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	"google.golang.org/grpc"
)

const (
	plgHi = 1 << iota
	plgAcc
	plgLogin
	plgSub
	plgLeave
	plgPub
	plgGet
	plgSet
	plgDel
	plgNote
	plgData
	plgMeta
	plgPres
	plgInfo

	plgClientMask = plgHi | plgAcc | plgLogin | plgSub | plgLeave | plgPub | plgGet | plgSet | plgDel | plgNote
	plgServerMask = plgData | plgMeta | plgPres | plgInfo
)

const (
	plgActCreate = 1 << iota
	plgActUpd
	plgActDel

	plgActMask = plgActCreate | plgActUpd | plgActDel
)

const (
	plgTopicMe = 1 << iota
	plgTopicFnd
	plgTopicP2P
	plgTopicGrp
	plgTopicSys
	plgTopicNew

	plgTopicCatMask = plgTopicMe | plgTopicFnd | plgTopicP2P | plgTopicGrp | plgTopicSys
)

const (
	plgFilterByTopicType = 1 << iota
	plgFilterByPacket
	plgFilterByAction
)

var (
	plgPacketNames = []string{
		"hi", "acc", "login", "sub", "leave", "pub", "get", "set", "del", "note",
		"data", "meta", "pres", "info",
	}

	plgTopicCatNames = []string{"me", "fnd", "p2p", "grp", "sys", "new"}
)

// Plugin defines client-side parameters of a gRPC plugin.
type Plugin struct {
	name    string
	timeout time.Duration
	// Filters for individual methods
	filterFireHose     *PluginFilter
	filterAccount      *PluginFilter
	filterTopic        *PluginFilter
	filterSubscription *PluginFilter
	filterMessage      *PluginFilter
	filterFind         bool
	failureCode        int
	failureText        string
	network            string
	addr               string

	conn   *grpc.ClientConn
	client chat.PluginClient
}

// PluginFilter is a enum which defines filtering types.
type PluginFilter struct {
	byPacket    int
	byTopicType int
	byAction    int
}
