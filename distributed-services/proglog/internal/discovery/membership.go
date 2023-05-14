package discovery

import (
	"net"

	"github.com/hashicorp/serf/serf"
	"go.uber.org/zap"
)

// Membership is our type wrapping Serf to provide discovery and cluster membership to our service.
type Membership struct {
	Config
	handler Handler
	serf    *serf.Serf
	events  chan serf.Event
	logger  *zap.Logger
}

func New(handler Handler, config Config) (*Membership, error) {
	c := &Membership{
		Config:  config,
		handler: handler,
		logger:  zap.L().Named("membership"),
	}
	if err := c.setupSerf(); err != nil {
		return nil, err
	}
	return c, nil
}

type Config struct {
	NodeName       string
	BindAddr       string
	Tags           map[string]string
	StartJoinAddrs []string
}

func (m *Membership) setupSerf() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", m.BindAddr)
	if err != nil {
		return err
	}

	config := serf.DefaultConfig()
	config.Init()
	config.MemberlistConfig.BindAddr = addr.IP.String()

	config.MemberlistConfig.BindPort = addr.Port

	m.events = make(chan serf.Event)

	config.EventCh = m.events

	config.Tags = m.Tags
	config.NodeName = m.Config.NodeName

	m.serf, err = serf.Create(config)

	if err != nil {

		return err

	}

	go m.eventHandler()
	if m.StartJoinAddrs != nil {
		_, err = m.serf.Join(m.StartJoinAddrs, true)
		if err != nil {
			return err
		}
	}
	return nil
}

type Handler interface {
	Join(name, addr string) error
	Leave(name string) error
}

func (m *Membership) eventHandler() {
	// Notice that Serf may coalesce multiple members updates into one event. For
	// example, say ten nodes join around the same time; in that case, Serf will send
	// you one join event with ten members, so that’s why we iterate over the event’s members.

	for e := range m.events {
		switch e.EventType() {
		case serf.EventMemberJoin:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					continue
				}
				m.handleJoin(member)
			}
		case serf.EventMemberLeave, serf.EventMemberFailed:
			for _, member := range e.(serf.MemberEvent).Members {
				if m.isLocal(member) {
					return
				}
				m.handleLeave(member)
			}
		}
	}
}

func (m *Membership) handleJoin(member serf.Member) {
	if err := m.handler.Join(member.Name, member.Tags["rpc_addr"]); err != nil {
		m.logError(err, "failed to join", member)
	}
}

func (m *Membership) handleLeave(member serf.Member) {
	if err := m.handler.Leave(member.Name); err != nil {
		m.logError(err, "failed to leave", member)
	}
}

// isLocal() returns whether the given Serf member is the local member by checking the members’ names.
func (m *Membership) isLocal(member serf.Member) bool {
	return m.serf.LocalMember().Name == member.Name
}

// Members() returns a point-in-time snapshot of the cluster’s Serf members.
func (m *Membership) Members() []serf.Member {
	return m.serf.Members()
}

// Leave() tells this member to leave the Serf cluster.
func (m *Membership) Leave() error {
	return m.serf.Leave()
}

// logError() logs the given error and message.
func (m *Membership) logError(err error, msg string, member serf.Member) {
	m.logger.Error(
		msg,
		zap.Error(err),
		zap.String("name", member.Name),
		zap.String("rpc_addr", member.Tags["rpc_addr"]),
	)
}
