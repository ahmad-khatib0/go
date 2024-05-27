package server

import (
	"errors"
	"net"
	"net/rpc"
	"sync/atomic"

	pt "github.com/ahmad-khatib0/go/websockets/chat/internal/push/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/ringhash"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
	"go.uber.org/zap"
)

// TopicMaster is a gRPC endpoint which receives requests sent by proxy topic to master topic.
func (c *Cluster) TopicMaster(msg *ClusterReq, rejected *bool) error {
	*rejected = false

	node := c.nodes[msg.Node]
	if node == nil {
		globals.l.Sugar().Warnf("cluster TopicMaster: request from an unknown node", msg.Node)
		return nil
	}

	// Master maintains one multiplexing session per proxy topic per node.
	// Except channel topics:
	// * one multiplexing session for channel subscriptions.
	// * one multiplexing session for group subscriptions.
	var msid string
	if msg.CliMsg != nil && types.IsChannel(msg.CliMsg.Original) {
		// If it's a channel request, use channel name.
		msid = msg.CliMsg.Original
	} else {
		msid = msg.RcptTo
	}
	// Append node name.
	msid += "-" + msg.Node
	msess := globals.sessionStore.Get(msid)

	if msg.Gone {
		// Proxy topic is gone. Tear down the local auxiliary session.
		// If it was the last session, master topic will shut down as well.
		node.stopMultiplexingSession(msess)

		if t := globals.hub.topicGet(msg.RcptTo); t != nil && t.isChan {
			// If it's a channel topic, also stop the "chnX-" local auxiliary session.
			msidChn := types.GrpToChn(t.name) + "-" + msg.Node
			node.stopMultiplexingSession(globals.sessionStore.Get(msidChn))
		}

		return nil
	}

	if msg.Signature != c.ring.Signature() {
		globals.l.Sugar().Warnf("cluster TopicMaster: session signature mismatch", msg.RcptTo)
		*rejected = true
		return nil
	}

	// Create a new multiplexing session if needed.
	if msess == nil {
		// If the session is not found, create it.
		var count int
		msess, count = globals.sessionStore.NewSession(node, msid)
		node.lock.Lock()
		node.msess[msid] = struct{}{}
		node.lock.Unlock()

		globals.l.Sugar().Infof("cluster: multiplexing session started", msid, count)
		msess.proxiedTopic = msg.RcptTo
	}

	// This is a local copy of a remote session.
	var sess *Session
	// Sess is nil for user agent changes and deferred presence notification requests.
	if msg.Sess != nil {
		// We only need some session info. No need to copy everything.
		sess = &Session{
			proto: PROXY,
			// Multiplexing session which actually handles the communication.
			multi: msess,
			// Local parameters specific to this session.
			sid:         msg.Sess.Sid,
			userAgent:   msg.Sess.UserAgent,
			remoteAddr:  msg.Sess.RemoteAddr,
			lang:        msg.Sess.Lang,
			countryCode: msg.Sess.CountryCode,
			proxyReq:    msg.ReqType,
			background:  msg.Sess.Background,
			uid:         msg.Sess.Uid,
		}
	}

	if msg.CliMsg != nil {
		msg.CliMsg.sess = sess
		msg.CliMsg.init = true
	}

	switch msg.ReqType {
	case ProxyReqJoin:
		select {
		case globals.hub.join <- msg.CliMsg:
		default:
			// Reply with a 500 to the user.
			sess.queueOut(ErrUnknownReply(msg.CliMsg, msg.CliMsg.Timestamp))
			globals.l.Sugar().Warnf(
				"cluster: join req failed - hub.join queue full, topic ",
				msg.CliMsg.RcptTo,
				"; orig sid ",
				sess.sid,
			)
		}

	case ProxyReqLeave:
		if t := globals.hub.topicGet(msg.RcptTo); t != nil {
			t.unreg <- msg.CliMsg
		} else {
			globals.l.Sugar().Warnf("cluster: leave request for unknown topic", msg.RcptTo)
		}

	case ProxyReqMeta:
		if t := globals.hub.topicGet(msg.RcptTo); t != nil {
			select {
			case t.meta <- msg.CliMsg:
			default:
				sess.queueOut(ErrUnknownReply(msg.CliMsg, msg.CliMsg.Timestamp))
				globals.l.Sugar().Warnf(
					"cluster: meta req failed - topic.meta queue full, topic ",
					msg.CliMsg.RcptTo,
					"; orig sid ",
					sess.sid,
				)
			}

		} else {
			globals.l.Sugar().Warnf("cluster: meta request for unknown topic", msg.RcptTo)
		}

	case ProxyReqBroadcast:
		select {
		case globals.hub.routeCli <- msg.CliMsg:
		default:
			globals.l.Error("cluster: route req failed - hub.route queue full")
		}

	case ProxyReqBgSession, ProxyReqMeUserAgent:
		// sess could be nil
		if t := globals.hub.topicGet(msg.RcptTo); t != nil {
			if t.supd == nil {
				globals.l.Sugar().Panicf("cluster: invalid topic category in session update", t.name, msg.ReqType)
			}

			su := &sessionUpdate{}
			if msg.ReqType == ProxyReqBgSession {
				su.sess = sess
			} else {
				su.userAgent = sess.userAgent
			}

			t.supd <- su

		} else {
			globals.l.Sugar().Warnf("cluster: session update for unknown topic", msg.RcptTo, msg.ReqType)
		}

	default:
		globals.l.Sugar().Warnf("cluster: unknown request type", msg.ReqType, msg.RcptTo)
		*rejected = true
	}

	return nil
}

// TopicProxy is a gRPC endpoint at topic proxy which receives topic master responses.
func (Cluster) TopicProxy(msg *ClusterResp, unused *bool) error {
	// This cluster member received a response from the topic master to be forwarded to the topic.
	// Find appropriate topic, send the message to it.
	if t := globals.hub.topicGet(msg.RcptTo); t != nil {
		msg.SrvMsg.uid = types.ParseUserId(msg.SrvMsg.AsUser)
		t.proxy <- msg
	} else {
		globals.l.Sugar().Warnf("cluster: master response for unknown topic", msg.RcptTo)
	}

	return nil
}

// Route endpoint receives intra-cluster messages destined for the nodes hosting the topic.
// Called by Hub.route channel consumer for messages send without attaching to topic first.
func (c *Cluster) Route(msg *ClusterRoute, rejected *bool) error {
	logError := func(err string) {
		sid := ""
		if msg.Sess != nil {
			sid = msg.Sess.Sid
		}

		globals.l.Error(err, zap.String("", sid))
		*rejected = true
	}

	*rejected = false
	if msg.Signature != c.ring.Signature() {
		logError("cluster Route: session signature mismatch")
		return nil
	}

	if msg.SrvMsg == nil {
		// TODO: maybe panic here.
		logError("cluster Route: nil server message")
		return nil
	}

	select {
	case globals.hub.routeSrv <- msg.SrvMsg:
	default:
		logError("cluster Route: server busy")
	}
	return nil
}

// Sends user cache update to user's Master node where the cache actually resides.
//
// The request is extected to contain users who reside at remote nodes only.
func (c *Cluster) routeUserReq(req *UserCacheReq) error {
	// Index requests by cluster node.
	reqByNode := make(map[string]*UserCacheReq)

	if req.PushRcpt != nil {
		// Request to send push notifications. Create separate packets for each affected cluster node.
		for uid, recipient := range req.PushRcpt.To {
			n := c.nodeForTopic(uid.UserId())
			if n == nil {
				return errors.New("attempt to update user at a non-existent node (1)")
			}
			r := reqByNode[n.name]
			if r == nil {
				r = &UserCacheReq{
					Node: c.thisNodeName,
					PushRcpt: &pt.Receipt{
						Payload: req.PushRcpt.Payload,
						To:      make(map[types.Uid]pt.Recipient),
					},
				}
			}
			r.PushRcpt.To[uid] = recipient
			reqByNode[n.name] = r
		}

	} else if len(req.UserIdList) > 0 {
		// Request to add/remove some users from cache.
		for _, uid := range req.UserIdList {
			n := c.nodeForTopic(uid.UserId())
			if n == nil {
				return errors.New("attempt to update user at a non-existent node (2)")
			}
			r := reqByNode[n.name]
			if r == nil {
				r = &UserCacheReq{Node: c.thisNodeName, Inc: req.Inc}
			}

			r.UserIdList = append(r.UserIdList, uid)
			reqByNode[n.name] = r
		}

	} else if req.Gone {
		// Message that the user is deleted is sent to all nodes.
		r := &UserCacheReq{Node: c.thisNodeName, UserIdList: req.UserIdList, Gone: true}
		for _, n := range c.nodes {
			reqByNode[n.name] = r
		}
	}

	if len(reqByNode) > 0 {
		for nodeName, r := range reqByNode {
			n := c.nodes[nodeName]
			var rejected bool
			err := n.call("Cluster.UserCacheUpdate", r, &rejected)
			if rejected {
				err = errors.New("master node out of sync")
			}
			if err != nil {
				return err
			}
		}
		return nil
	}

	// Update to cached values.
	n := c.nodeForTopic(req.UserId.UserId())
	if n == nil {
		return errors.New("attempt to update user at a non-existent node (3)")
	}
	req.Node = c.thisNodeName
	var rejected bool
	err := n.call("Cluster.UserCacheUpdate", req, &rejected)
	if rejected {
		err = errors.New("master node out of sync")
	}
	return err
}

// Given topic name, find appropriate cluster node to route message to.
func (c *Cluster) nodeForTopic(topic string) *ClusterNode {
	key := c.ring.Get(topic)
	if key == c.thisNodeName {
		globals.l.Error("cluster: request to route to self")
		// Do not route to self
		return nil
	}

	node := c.nodes[key]
	if node == nil {
		globals.l.Sugar().Warnf("cluster: no node for topic", topic, key)
	}
	return node
}

// isRemoteTopic checks if the given topic is handled by this node or a remote node.
func (c *Cluster) isRemoteTopic(topic string) bool {
	if c == nil {
		// Cluster not initialized, all topics are local
		return false
	}
	return c.ring.Get(topic) != c.thisNodeName
}

// genLocalTopicName is just like genTopicName(), but the generated name belongs to the current cluster node.
func (c *Cluster) genLocalTopicName() string {
	topic := genTopicName()
	if c == nil {
		// Cluster not initialized, all topics are local
		return topic
	}

	// TODO: if cluster is large it may become too inefficient.
	for c.ring.Get(topic) != c.thisNodeName {
		topic = genTopicName()
	}
	return topic
}

// isPartitioned checks if the cluster is partitioned due to network or other
//
// failure and if the current node is a part of the smaller partition.
func (c *Cluster) isPartitioned() bool {
	if c == nil || c.fo == nil {
		// Cluster not initialized or failover disabled therefore not partitioned.
		return false
	}

	c.fo.activeNodesLock.RLock()
	result := (len(c.nodes)+1)/2 >= len(c.fo.activeNodes)
	c.fo.activeNodesLock.RUnlock()

	return result
}

func (c *Cluster) makeClusterReq(reqType ProxyReqType, msg *ClientComMessage, topic string, sess *Session) *ClusterReq {
	req := &ClusterReq{
		Node:        c.thisNodeName,
		Signature:   c.ring.Signature(),
		Fingerprint: c.fingerprint,
		ReqType:     reqType,
		RcptTo:      topic,
	}

	var uid types.Uid

	if msg != nil {
		req.CliMsg = msg
		uid = types.ParseUserId(req.CliMsg.AsUser)
	}

	if sess != nil {
		if uid.IsZero() {
			uid = sess.uid
		}

		req.Sess = &ClusterSess{
			Uid:         uid,
			AuthLvl:     sess.authLvl,
			RemoteAddr:  sess.remoteAddr,
			UserAgent:   sess.userAgent,
			Ver:         sess.ver,
			Lang:        sess.lang,
			CountryCode: sess.countryCode,
			DeviceID:    sess.deviceID,
			Platform:    sess.platf,
			Sid:         sess.sid,
			Background:  sess.background,
		}
	}
	return req
}

// Forward client request message from the Topic Proxy to the Topic Master (cluster node which owns the topic).
func (c *Cluster) routeToTopicMaster(reqType ProxyReqType, msg *ClientComMessage, topic string, sess *Session) error {
	if c == nil {
		// Cluster may be nil due to shutdown.
		return nil
	}

	if sess != nil && reqType != ProxyReqLeave {
		if atomic.LoadInt32(&sess.terminating) > 0 {
			// The session is terminating.
			// Do not forward any requests except "leave" to the topic master.
			return nil
		}
	}

	req := c.makeClusterReq(reqType, msg, topic, sess)

	// Find the cluster node which owns the topic, then forward to it.
	n := c.nodeForTopic(topic)
	if n == nil {
		return errors.New("node for topic not found")
	}
	return n.proxyToMasterAsync(req)
}

// Forward server response message to the node that owns topic.
func (c *Cluster) routeToTopicIntraCluster(topic string, msg *ServerComMessage, sess *Session) error {
	if c == nil {
		// Cluster may be nil due to shutdown.
		return nil
	}

	n := c.nodeForTopic(topic)
	if n == nil {
		return errors.New("node for topic not found (intra)")
	}

	route := &ClusterRoute{
		Node:        c.thisNodeName,
		Signature:   c.ring.Signature(),
		Fingerprint: c.fingerprint,
		SrvMsg:      msg,
	}

	if sess != nil {
		route.Sess = &ClusterSess{Sid: sess.sid}
	}
	return n.route(route)
}

// Topic proxy terminated. Inform remote Master node that the proxy is gone.
func (c *Cluster) topicProxyGone(topicName string) error {
	if c == nil {
		// Cluster may be nil due to shutdown.
		return nil
	}

	// Find the cluster node which owns the topic, then forward to it.
	n := c.nodeForTopic(topicName)
	if n == nil {
		return errors.New("node for topic not found")
	}

	req := c.makeClusterReq(ProxyReqLeave, nil, topicName, nil)
	req.Gone = true
	return n.proxyToMasterAsync(req)
}

// Proxied session is being closed at the Master node.
func (sess *Session) closeRPC() {
	if sess.isMultiplex() {
		globals.l.Sugar().Infof("cluster: session proxy closed", sess.sid)
	}
}

// Start accepting connections.
func (c *Cluster) start() {
	addr, err := net.ResolveTCPAddr("tcp", c.listenOn)
	if err != nil {
		globals.l.Fatal("", zap.Error(err))
	}

	c.inbound, err = net.ListenTCP("tcp", addr)

	if err != nil {
		globals.l.Fatal("", zap.Error(err))
	}

	for _, n := range c.nodes {
		go n.reconnect()
		n.rpcDone = make(chan *rpc.Call, len(c.nodes)*clusterRpcCompletionBuffer)
		n.p2mSender = make(chan *ClusterReq, clusterProxyToMasterBuffer)
		go n.asyncRpcLoop()
		go n.p2mSenderLoop()
	}

	if c.fo != nil {
		go c.run()
	}

	err = rpc.Register(c)
	if err != nil {
		globals.l.Fatal("", zap.Error(err))
	}

	go rpc.Accept(c.inbound)

	globals.l.Sugar().Infof(
		"Cluster of %d nodes initialized, node '%s' is listening on [%s]",
		len(globals.cluster.nodes)+1,
		globals.cluster.thisNodeName,
		c.listenOn,
	)
}

func (c *Cluster) shutdown() {
	if globals.cluster == nil {
		return
	}
	for _, n := range c.nodes {
		close(n.rpcDone)
		close(n.p2mSender)
	}

	globals.cluster.proxyEventQueue.Stop()
	globals.cluster = nil

	c.inbound.Close()

	if c.fo != nil {
		c.fo.done <- true
	}

	for _, n := range c.nodes {
		n.done <- true
	}

	globals.l.Info("Cluster shut down")
}

// Recalculate the ring hash using provided list of nodes or only nodes in a non-failed state.
//
// Returns the list of nodes used for ring hash.
func (c *Cluster) rehash(nodes []string) []string {
	ring := ringhash.New(clusterHashReplicas, nil)

	var ringKeys []string

	if nodes == nil {
		for _, node := range c.nodes {
			ringKeys = append(ringKeys, node.name)
		}
		ringKeys = append(ringKeys, c.thisNodeName)
	} else {
		ringKeys = append(ringKeys, nodes...)
	}
	ring.Add(ringKeys...)

	c.ring = ring

	return ringKeys
}

// invalidateProxySubs iterates over sessions proxied on this node and for each session
//
// sends "{pres term}" informing that the topic subscription (attachment) was lost:
//
// - Called immediately after Cluster.rehash() for all relocated topics (forNode == "").
//
// - Called for topics hosted at a specific node when a node restart is detected.
// TODO: consider resubscribing to topics instead of forcing sessions to resubscribe.
func (c *Cluster) invalidateProxySubs(forNode string) {
	sessions := make(map[*Session][]string)
	globals.hub.topics.Range(func(_, v any) bool {
		topic := v.(*Topic)
		if !topic.isProxy {
			// Topic isn't a proxy.
			return true
		}
		if forNode == "" {
			if topic.masterNode == c.ring.Get(topic.name) {
				// The topic hasn't moved. Continue.
				return true
			}
		} else if topic.masterNode != forNode {
			// The topic is hosted at a different node than the restarted node.
			return true
		}

		for s, psd := range topic.sessions {
			// FIXME: 'me' topic must be the last one in the list for each topic.
			sessions[s] = append(sessions[s], topicNameForUser(topic.name, psd.uid, psd.isChanSub))
		}
		return true
	})

	for s, topicsToTerminate := range sessions {
		s.presTermDirect(topicsToTerminate)
	}
}

// gcProxySessions terminates orphaned proxy sessions at a master node for all lost nodes (allNodes minus activeNodes).
// The session is orphaned when the origin node is gone.
func (c *Cluster) gcProxySessions(activeNodes []string) {
	allNodes := []string{c.thisNodeName}
	for name := range c.nodes {
		allNodes = append(allNodes, name)
	}

	_, failedNodes, _ := stringSliceDelta(allNodes, activeNodes)
	for _, node := range failedNodes {
		// Iterate sessions of a failed node
		c.gcProxySessionsForNode(node)
	}
}

// gcProxySessionsForNode terminates orphaned proxy sessions at a master node for the given node.
//
// For example, a remote node is restarted or the cluster is rehashed without the node.
func (c *Cluster) gcProxySessionsForNode(node string) {
	n := c.nodes[node]
	n.lock.Lock()
	msess := n.msess
	n.msess = make(map[string]struct{})
	n.lock.Unlock()
	for sid := range msess {
		if sess := globals.sessionStore.Get(sid); sess != nil {
			sess.stop <- nil
		}
	}
}

// clusterWriteLoop implements write loop for multiplexing (proxy) session at a node which hosts master topic.
// The session is a multiplexing session, i.e. it handles requests for multiple sessions at origin.
func (sess *Session) clusterWriteLoop(forTopic string) {
	terminate := true
	defer func() {
		if terminate {
			sess.closeRPC()
			globals.sessionStore.Delete(sess)
			sess.inflightReqs = nil
			sess.unsubAll()
		}
	}()

	for {
		select {
		case msg, ok := <-sess.send:
			if !ok || sess.clnode.endpoint == nil {
				// channel closed
				return
			}
			srvMsg := msg.(*ServerComMessage)
			response := &ClusterResp{SrvMsg: srvMsg}
			if srvMsg.sess == nil {
				response.OrigSid = "*"
			} else {
				response.OrigReqType = srvMsg.sess.proxyReq
				response.OrigSid = srvMsg.sess.sid
				srvMsg.AsUser = srvMsg.sess.uid.UserId()

				switch srvMsg.sess.proxyReq {
				case ProxyReqJoin, ProxyReqLeave, ProxyReqMeta, ProxyReqBgSession, ProxyReqMeUserAgent, ProxyReqCall:
				// Do nothing
				case ProxyReqBroadcast, ProxyReqNone:
					if srvMsg.Data != nil || srvMsg.Pres != nil || srvMsg.Info != nil {
						response.OrigSid = "*"
					} else if srvMsg.Ctrl == nil {
						globals.l.Sugar().Warnf(
							"cluster: request type not set in clusterWriteLoop",
							sess.sid,
							srvMsg.describe(),
							"src_sid:",
							srvMsg.sess.sid,
						)
					}
				default:
					globals.l.Sugar().Panicf("cluster: unknown request type in clusterWriteLoop", srvMsg.sess.proxyReq)
				}
			}

			srvMsg.RcptTo = forTopic
			response.RcptTo = forTopic

			if err := sess.clnode.masterToProxyAsync(response); err != nil {
				globals.l.Sugar().Warnf("cluster: response to proxy failed \"%s\": %s", sess.sid, err.Error())
				return
			}

		case msg := <-sess.stop:
			if msg == nil {
				// Terminating multiplexing session.
				return
			}
			// There are two cases of msg != nil:
			//  * user is being deleted
			//  * node shutdown
			// In both cases the msg does not need to be forwarded to the proxy.

		case <-sess.detach:
			return
		default:
			terminate = false
			return
		}
	}
}

// ******************************************************************
// ******************************************************************
// ******************************************************************
// User cache & push notifications management. These are calls received by the
//
// Master from Proxy. The Proxy expects no payload to be returned by the master.

// UserCacheUpdate endpoint receives updates to user's cached values as well as sends push notifications.
func (c *Cluster) UserCacheUpdate(msg *UserCacheReq, rejected *bool) error {
	if msg.Gone {
		// User is deleted. Evict all user's sessions.
		globals.sessionStore.EvictUser(msg.UserId, "")

		if globals.cluster.isRemoteTopic(msg.UserId.UserId()) {
			// No need to delete user's cache if user is remote.
			return nil
		}
	}

	usersRequestFromCluster(msg)
	return nil
}

// Ping is a gRPC endpoint which receives ping requests from peer nodes.Used to detect node restarts.
func (c *Cluster) Ping(ping *ClusterPing, unused *bool) error {
	node := c.nodes[ping.Node]
	if node == nil {
		globals.l.Sugar().Warnf("cluster Ping from unknown node", ping.Node)
		return nil
	}

	if node.fingerprint == 0 {
		// This is the first connection to remote node.
		node.fingerprint = ping.Fingerprint
	} else if node.fingerprint != ping.Fingerprint {
		// Remote node restarted.
		node.fingerprint = ping.Fingerprint
		c.invalidateProxySubs(ping.Node)
		c.gcProxySessionsForNode(ping.Node)
	}

	return nil
}
