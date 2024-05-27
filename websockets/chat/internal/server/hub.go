package server

import (
	"strings"
	"sync"
	"sync/atomic"
	"time"

	auth "github.com/ahmad-khatib0/go/websockets/chat/internal/auth/types"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/constants"
	"github.com/ahmad-khatib0/go/websockets/chat/internal/store/types"
)

func (h *Hub) topicGet(name string) *Topic {
	if t, ok := h.topics.Load(name); ok {
		return t.(*Topic)
	}
	return nil
}

func (h *Hub) topicPut(name string, t *Topic) {
	h.numTopics++
	h.topics.Store(name, t)
}

func (h *Hub) topicDel(name string) {
	h.numTopics--
	h.topics.Delete(name)
}

func newHub() *Hub {
	h := &Hub{
		topics: &sync.Map{},
		// TODO: verify if these channels have to be buffered.
		routeCli:   make(chan *ClientComMessage, 4096),
		routeSrv:   make(chan *ServerComMessage, 4096),
		join:       make(chan *ClientComMessage, 256),
		unreg:      make(chan *topicUnreg, 256),
		rehash:     make(chan bool),
		meta:       make(chan *ClientComMessage, 128),
		userStatus: make(chan *userStatusReq, 128),
		shutdown:   make(chan chan<- bool),
	}

	globals.stats.RegisterInt("LiveTopics")
	globals.stats.RegisterInt("TotalTopics")

	globals.stats.RegisterInt("IncomingMessagesWebsockTotal")
	globals.stats.RegisterInt("OutgoingMessagesWebsockTotal")

	globals.stats.RegisterInt("IncomingMessagesLongpollTotal")
	globals.stats.RegisterInt("OutgoingMessagesLongpollTotal")

	globals.stats.RegisterInt("IncomingMessagesGrpcTotal")
	globals.stats.RegisterInt("OutgoingMessagesGrpcTotal")

	globals.stats.RegisterInt("FileDownloadsTotal")
	globals.stats.RegisterInt("FileUploadsTotal")

	globals.stats.RegisterInt("CtrlCodesTotal2xx")
	globals.stats.RegisterInt("CtrlCodesTotal3xx")
	globals.stats.RegisterInt("CtrlCodesTotal4xx")
	globals.stats.RegisterInt("CtrlCodesTotal5xx")

	globals.stats.RegisterHistogram("RequestLatency", requestLatencyDistribution)
	globals.stats.RegisterHistogram("OutgoingMessageSize", outgoingMessageSizeDistribution)

	go h.run()

	// Initialize 'sys' topic. It will be initialized either as master or proxy.
	h.join <- &ClientComMessage{RcptTo: "sys", Original: "sys"}

	return h
}

func (h *Hub) run() {
	for {
		select {
		case join := <-h.join:
			// Handle a subscription request:
			// 1. Init topic
			// 1.1 If a new topic is requested, create it
			// 1.2 If a new subscription to an existing topic is requested:
			// 1.2.1 check if topic is already loaded
			// 1.2.2 if not, load it
			// 1.2.3 if it cannot be loaded (not found), fail
			// 2. Check access rights and reject, if appropriate
			// 3. Attach session to the topic
			// Is the topic already loaded?
			t := h.topicGet(join.RcptTo)
			if t == nil {
				// Topic does not exist or not loaded.
				t = &Topic{
					name:      join.RcptTo,
					xoriginal: join.Original,
					// Indicates a proxy topic.
					isProxy:   globals.cluster.isRemoteTopic(join.RcptTo),
					sessions:  make(map[*Session]perSessionData),
					clientMsg: make(chan *ClientComMessage, 192),
					serverMsg: make(chan *ServerComMessage, 64),
					reg:       make(chan *ClientComMessage, 256),
					unreg:     make(chan *ClientComMessage, 256),
					meta:      make(chan *ClientComMessage, 64),
					perUser:   make(map[types.Uid]perUserData),
					exit:      make(chan *shutDown, 1),
				}
				if globals.cluster != nil {
					if t.isProxy {
						t.proxy = make(chan *ClusterResp, 32)
						t.masterNode = globals.cluster.ring.Get(t.name)
					} else {
						// It's a master topic. Make a channel for handling
						// direct messages from the proxy.
						t.master = make(chan *ClusterSessUpdate, 8)
					}
				}

				// Topic is created in suspended state because it's not yet configured.
				t.markPaused(true)
				// Save topic now to prevent race condition.
				h.topicPut(join.RcptTo, t)

				// Configure the topic.
				go topicInit(t, join, h)
			} else {
				// Topic found.
				if t.isInactive() {
					// Topic is either not ready or being deleted.
					if join.sess.inflightReqs != nil {
						join.sess.inflightReqs.Done()
					}
					join.sess.queueOut(ErrLockedReply(join, join.Timestamp))
					continue
				}
				// Topic will check access rights and send appropriate {ctrl}
				select {
				case t.reg <- join:
				default:
					if join.sess.inflightReqs != nil {
						join.sess.inflightReqs.Done()
					}
					join.sess.queueOut(ErrServiceUnavailableReply(join, join.Timestamp))
					globals.l.Sugar().Errorf(
						"hub.join loop: topic's reg queue full",
						join.RcptTo,
						join.sess.sid,
						" - total queue len:",
						len(t.reg),
					)
				}
			}

		case msg := <-h.routeCli:
			// This is a message from a session not subscribed to topic
			// Route incoming message to topic if topic permits such routing.
			if dst := h.topicGet(msg.RcptTo); dst != nil {
				// Everything is OK, sending packet to known topic
				if dst.clientMsg != nil {
					select {
					case dst.clientMsg <- msg:
					default:
						globals.l.Sugar().Errorf("hub: topic's broadcast queue is full", dst.name)
					}
				} else {
					globals.l.Sugar().Warnf("hub: invalid topic category for broadcast", dst.name)
				}
			} else if msg.Note == nil {
				// Topic is unknown or offline.
				// Note is silently ignored, all other messages are reported as accepted to prevent
				// clients from guessing valid topic names.

				// TODO(gene): validate topic name, discarding invalid topics.

				globals.l.Sugar().Infof("Hub. Topic[%s] is unknown or offline", msg.RcptTo)

				msg.sess.queueOut(NoErrAcceptedExplicitTs(msg.Id, msg.RcptTo, types.TimeNow(), msg.Timestamp))
			}

		case msg := <-h.routeSrv:
			// This is a server message from a connection not subscribed to topic
			// Route incoming message to topic if topic permits such routing.
			if dst := h.topicGet(msg.RcptTo); dst != nil {
				// Everything is OK, sending packet to known topic
				select {
				case dst.serverMsg <- msg:
				default:
					globals.l.Sugar().Errorf("hub: topic's broadcast queue is full", dst.name)
				}
			} else if (strings.HasPrefix(msg.RcptTo, "usr") || strings.HasPrefix(msg.RcptTo, "grp")) &&
				globals.cluster.isRemoteTopic(msg.RcptTo) {
				// It is a remote topic.
				if err := globals.cluster.routeToTopicIntraCluster(msg.RcptTo, msg, msg.sess); err != nil {
					globals.l.Sugar().Infof("hub: routing to '%s' failed", msg.RcptTo)
				}
			}

		case msg := <-h.meta:
			// Metadata read or update from a user who is not attached to the topic.
			if msg.Get != nil {
				if msg.MetaWhat == constMsgMetaDesc {
					go replyOfflineTopicGetDesc(msg.sess, msg)
				} else {
					go replyOfflineTopicGetSub(msg.sess, msg)
				}
			} else if msg.Set != nil {
				go replyOfflineTopicSetSub(msg.sess, msg)
			}

		case status := <-h.userStatus:
			// Suspend/activate user's topics.
			go h.topicsStateForUser(status.forUser, status.state == types.StateSuspended)

		case unreg := <-h.unreg:
			reason := StopNone
			if unreg.del {
				reason = StopDeleted
			}
			if unreg.forUser.IsZero() {
				// The topic is being garbage collected or deleted.
				if err := h.topicUnreg(unreg.sess, unreg.rcptTo, unreg.pkt, reason); err != nil {
					globals.l.Sugar().Errorf("hub.topicUnreg failed:", err)
				}
			} else {
				// User is being deleted.
				go h.stopTopicsForUser(unreg.forUser, reason, unreg.done)
			}

		case <-h.rehash:
			// Cluster rehashing. Some previously local topics became remote,
			// and the other way round.
			// Such topics must be shut down at this node.
			h.topics.Range(func(_, t any) bool {
				topic := t.(*Topic)
				// Handle two cases:
				// 1. Master topic has moved out to another node.
				// 2. Proxy topic is running on a new master node
				//    (i.e. the master topic has moved to this node).
				if topic.isProxy != globals.cluster.isRemoteTopic(topic.name) {
					h.topicUnreg(nil, topic.name, nil, StopRehashing)
				}
				return true
			})

			// Check if 'sys' topic has migrated to this node.
			if h.topicGet("sys") == nil && !globals.cluster.isRemoteTopic("sys") {
				// Yes, 'sys' has migrated here. Initialize it.
				// The h.join is unbuffered. We must call from another goroutine. Otherwise deadlock.
				go func() {
					h.join <- &ClientComMessage{RcptTo: "sys", Original: "sys"}
				}()
			}

		case hubdone := <-h.shutdown:
			// start cleanup process
			topicsdone := make(chan bool)
			topicCount := 0
			h.topics.Range(func(_, topic any) bool {
				topic.(*Topic).exit <- &shutDown{done: topicsdone}
				topicCount++
				return true
			})

			for i := 0; i < topicCount; i++ {
				<-topicsdone
			}

			globals.l.Sugar().Infof("Hub shutdown completed with %d topics", topicCount)

			// let the main goroutine know we are done with the cleanup
			hubdone <- true

			return

		case <-time.After(constants.IdleSessionTimeout):
		}
	}
}

// Update state of all topics associated with the given user:
//
// * all p2p topics with the given user
//
// * group topics where the given user is the owner.
//
// 'me' and fnd' are ignored here because they are direcly tied to the user object.
func (h *Hub) topicsStateForUser(uid types.Uid, suspended bool) {
	h.topics.Range(func(name any, t any) bool {
		topic := t.(*Topic)
		if topic.cat == types.TopicCatMe || topic.cat == types.TopicCatFnd {
			return true
		}

		if _, isMember := topic.perUser[uid]; (topic.cat == types.TopicCatP2P && isMember) || topic.owner == uid {
			topic.markReadOnly(suspended)

			// Don't send "off" notification on suspension. They will be sent when the user is evicted.
		}
		return true
	})
}

// topicUnreg deletes or unregisters the topic:
//
// Cases:
//
// 1. Topic being deleted
//
// 1.1 Topic is online
//
// 1.1.1 If the requester is the owner or if it's the last sub in a p2p topic (p2p may be sent internally when the last user unsubscribes):
//
// 1.1.1.1 Tell topic to stop accepting requests.
//
// 1.1.1.2 Hub deletes the topic from database
//
// 1.1.1.3 Hub unregisters the topic
//
// 1.1.1.4 Hub informs the origin of success or failure
//
// 1.1.1.5 Hub forwards request to topic
//
// 1.1.1.6 Topic evicts all sessions
//
// 1.1.1.7 Topic exits the run() loop
//
// 1.1.2 If the requester is not the owner
//
// 1.1.2.1 Send it to topic to be treated like {leave unsub=true}
//
// 1.2 Topic is offline
//
// 1.2.1 If requester is the owner
//
// 1.2.1.1 Hub deletes topic from database
//
// 1.2.2 If not the owner
//
// 1.2.2.1 Delete subscription from DB
//
// 1.2.3 Hub informs the origin of success or failure
//
// 1.2.4 Send notification to subscribers that the topic was deleted
//
// 2. Topic is just being unregistered (topic is going offline)
//
// 2.1 Unregister it with no further action
func (h *Hub) topicUnreg(sess *Session, topic string, msg *ClientComMessage, reason int) error {
	now := types.TimeNow()

	// TODO: when channel is deleted unsubscribe all devices from channel's FCM topic.

	if reason == StopDeleted {
		var asUid types.Uid
		if msg != nil {
			asUid = types.ParseUserId(msg.AsUser)
		}

		// Case 1 (unregister and delete)
		if t := h.topicGet(topic); t != nil {
			// Case 1.1: topic is online
			if (!asUid.IsZero() && t.owner == asUid) || (t.cat == types.TopicCatP2P && t.subsCount() < 2) {
				// Case 1.1.1: requester is the owner or last sub in a p2p topic
				t.markPaused(true)
				hard := true
				if msg != nil && msg.Del != nil {
					// Soft-deleting does not make sense for p2p topics.
					hard = msg.Del.Hard || t.cat == types.TopicCatP2P
				}
				if err := globals.store.TopDelete(topic, t.isChan, hard); err != nil {
					t.markPaused(false)
					if sess != nil {
						sess.queueOut(ErrUnknownReply(msg, now))
					}
					return err
				}

				if sess != nil {
					sess.queueOut(NoErrReply(msg, now))
				}

				if t.isChan {
					// Notify channel subscribers that the channel is deleted.
					sendPush(pushForChanDelete(t.name, now))
				}

				h.topicDel(topic)
				t.markDeleted()
				t.exit <- &shutDown{reason: StopDeleted}
				globals.stats.IntStatsInc("LiveTopics", -1)

			} else {
				// Case 1.1.2: requester is NOT the owner or not empty P2P.
				msg.MetaWhat = constMsgDelTopic
				msg.sess = sess
				t.meta <- msg
			}

		} else {
			// Case 1.2: topic is offline.

			// Is user a channel subscriber? Use chnABC instead of grpABC and get only this user's subscription.
			var opts *types.QueryOpt
			if types.IsChannel(msg.Original) {
				topic = msg.Original
				opts = &types.QueryOpt{User: asUid}
			}

			// Get all subscribers of non-channel topics: we need to know how many are left and notify them.
			// Get only one subscription for channel users.
			subs, err := globals.store.TopGetTopicSubs(topic, opts)
			if err != nil {
				sess.queueOut(ErrUnknownReply(msg, now))
				return err
			}

			tcat := topicCat(topic)
			if len(subs) == 0 {
				if tcat == types.TopicCatP2P {
					// No subscribers: delete.
					globals.store.TopDelete(topic, false, true)
				}

				sess.queueOut(InfoNoActionReply(msg, now))
				return nil
			}

			// Find subscription of the current user.
			var sub *types.Subscription
			user := asUid.String()
			for i := range subs {
				if subs[i].User == user {
					sub = &subs[i]
					break
				}
			}

			if sub == nil {
				// If user has no subscription, tell him all is fine
				sess.queueOut(InfoNoActionReply(msg, now))
				return nil
			}

			if !(sub.ModeGiven & sub.ModeWant).IsOwner() {
				// Case 1.2.2.1 Not the owner, but possibly last subscription in a P2P topic.

				if tcat == types.TopicCatP2P && len(subs) < 2 {
					// This is a P2P topic and fewer than 2 subscriptions, delete the entire topic
					if err := globals.store.TopDelete(topic, false, msg.Del.Hard); err != nil {
						sess.queueOut(ErrUnknownReply(msg, now))
						return err
					}
					// Inform plugin that the topic was deleted.
					pluginTopic(&Topic{name: topic}, plgActDel)
				} else if err := globals.store.SubsDelete(topic, asUid); err != nil {
					// Not P2P or more than 1 subscription left.
					// Delete user's own subscription only
					if err == types.ErrNotFound {
						sess.queueOut(InfoNoActionReply(msg, now))
						err = nil
					} else {
						sess.queueOut(ErrUnknownReply(msg, now))
					}
					return err
				}

				// Notify user's other sessions that the subscription is gone
				presSingleUserOfflineOffline(asUid, msg.Original, "gone", nilPresParams, sess.sid)
				if tcat == types.TopicCatP2P && len(subs) == 2 {
					uname1 := asUid.UserId()
					uid2 := types.ParseUserId(msg.Original)
					// Tell user1 to stop sending updates to user2 without passing change to user1's sessions.
					presSingleUserOfflineOffline(asUid, uid2.UserId(), "?none+rem", nilPresParams, "")
					// Don't change the online status of user1, just ask user2 to stop notification exchange.
					// Tell user2 that user1 is offline but let him keep sending updates in case user1 resubscribes.
					presSingleUserOfflineOffline(uid2, uname1, "off", nilPresParams, "")
				}

				// Inform plugin that the subscription was deleted.
				pluginSubscription(sub, plgActDel)
			} else {
				// Case 1.2.1.1: owner, delete the group topic from db. Only group topics have owners.
				// We don't know if the group topic is a channel, but cleaning it as a channel does no harm
				// other than a small performance penalty.
				if err := globals.store.TopDelete(topic, true, msg.Del.Hard); err != nil {
					sess.queueOut(ErrUnknownReply(msg, now))
					return err
				}

				// Notify subscribers that the group topic is gone.
				presSubsOfflineOffline(topic, tcat, subs, "gone", &presParams{}, sess.sid)

				// Notify channel subscribers that the channel is deleted.
				// The push will not be delivered to anybody if the topic is not a channel.
				sendPush(pushForChanDelete(topic, now))

				// Inform plugin that the topic was deleted.
				pluginTopic(&Topic{name: topic}, plgActDel)
			}

			sess.queueOut(NoErrReply(msg, now))
		}
	} else {
		// Case 2: just unregister.
		// If t is nil, it's not registered, no action is needed
		if t := h.topicGet(topic); t != nil {
			t.markDeleted()
			h.topicDel(topic)

			t.exit <- &shutDown{reason: reason}

			globals.stats.IntStatsInc("LiveTopics", -1)
		}

		// sess && msg could be nil if the topic is being killed by timer or due to rehashing.
		if sess != nil && msg != nil {
			sess.queueOut(NoErrReply(msg, now))
		}
	}

	return nil
}

// Terminate all topics associated with the given user:
//
// * all p2p topics with the given user
//
// * group topics where the given user is the owner.
//
// * user's 'me' and 'fnd' topics.
func (h *Hub) stopTopicsForUser(uid types.Uid, reason int, alldone chan<- bool) {
	var done chan bool
	if alldone != nil {
		done = make(chan bool, 128)
	}

	count := 0
	h.topics.Range(func(name any, t any) bool {
		topic := t.(*Topic)
		if _, isMember := topic.perUser[uid]; (topic.cat != types.TopicCatGrp && isMember) ||
			topic.owner == uid {
			topic.markDeleted()
			h.topics.Delete(name)

			// This call is non-blocking unless some other routine tries to stop it at the same time.
			topic.exit <- &shutDown{reason: reason, done: done}

			// Just send to p2p topics here.
			if topic.cat == types.TopicCatP2P && len(topic.perUser) == 2 {
				presSingleUserOfflineOffline(topic.p2pOtherUser(uid), uid.UserId(), "gone", nilPresParams, "")
			}
			count++
		}
		return true
	})

	globals.stats.IntStatsInc("LiveTopics", -count)

	if alldone != nil {
		for i := 0; i < count; i++ {
			<-done
		}
		alldone <- true
	}
}

// replyOfflineTopicGetDesc reads a minimal topic Desc from the database.
// The requester may or maynot be subscribed to the topic.
func replyOfflineTopicGetDesc(sess *Session, msg *ClientComMessage) {
	now := types.TimeNow()
	desc := &MsgTopicDesc{}
	asUid := types.ParseUserId(msg.AsUser)
	topic := msg.RcptTo

	if strings.HasPrefix(topic, "grp") || topic == "sys" {
		stopic, err := globals.store.TopGet(topic)
		if err != nil {
			globals.l.Sugar().Infof("replyOfflineTopicGetDesc", err)
			sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
			return
		}
		if stopic == nil {
			sess.queueOut(ErrTopicNotFoundReply(msg, now))
			return
		}

		desc.CreatedAt = &stopic.CreatedAt
		desc.UpdatedAt = &stopic.UpdatedAt
		desc.Public = stopic.Public
		desc.Trusted = stopic.Trusted
		desc.IsChan = stopic.UseBt
		if stopic.Owner == msg.AsUser {
			desc.DefaultAcs = &MsgDefaultAcsMode{
				Auth: stopic.Access.Auth.String(),
				Anon: stopic.Access.Anon.String(),
			}
		}
		// Report appropriate access level. Could be overridden below if subscription exists.
		desc.Acs = &MsgAccessMode{}
		if sess.authLvl == auth.LevelAuth || sess.authLvl == auth.LevelRoot {
			desc.Acs.Mode = stopic.Access.Auth.String()
		} else if sess.authLvl == auth.LevelAnon {
			desc.Acs.Mode = stopic.Access.Anon.String()
		}
	} else {
		// 'me' and p2p topics
		uid := types.ZeroUid
		if strings.HasPrefix(topic, "usr") {
			// User specified as usrXXX
			uid = types.ParseUserId(topic)
			topic = asUid.P2PName(uid)

		} else if strings.HasPrefix(topic, "p2p") {
			// User specified as p2pXXXYYY
			uid1, uid2, _ := types.ParseP2P(topic)
			if uid1 == asUid {
				uid = uid2
			} else if uid2 == asUid {
				uid = uid1
			}
		}

		if uid.IsZero() {
			globals.l.Sugar().Warnf("replyOfflineTopicGetDesc: malformed p2p topic name")
			sess.queueOut(ErrMalformedReply(msg, now))
			return
		}

		suser, err := globals.store.UsersGet(uid)
		if err != nil {
			sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
			return
		}

		if suser == nil {
			sess.queueOut(ErrUserNotFoundReply(msg, now))
			return
		}

		desc.CreatedAt = &suser.CreatedAt
		desc.UpdatedAt = &suser.UpdatedAt
		desc.Public = suser.Public
		desc.Trusted = suser.Trusted
		if sess.authLvl == auth.LevelRoot {
			desc.State = suser.State.String()
		}

		// Report appropriate access level. Could be overridden below if subscription exists.
		desc.Acs = &MsgAccessMode{}
		if sess.authLvl == auth.LevelAuth || sess.authLvl == auth.LevelRoot {
			desc.Acs.Mode = suser.Access.Auth.String()
		} else if sess.authLvl == auth.LevelAnon {
			desc.Acs.Mode = suser.Access.Anon.String()
		}
	}

	sub, err := globals.store.SubsGetSubs(topic, asUid, false)
	if err != nil {
		globals.l.Sugar().Warnf("replyOfflineTopicGetDesc:", err)
		sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
		return
	}

	if sub != nil {
		desc.Private = sub.Private
		// FIXME: suspended topics should get no AW access.
		desc.Acs = &MsgAccessMode{
			Want:  sub.ModeWant.String(),
			Given: sub.ModeGiven.String(),
			Mode:  (sub.ModeGiven & sub.ModeWant).String(),
		}
	}

	sess.queueOut(&ServerComMessage{
		Meta: &MsgServerMeta{
			Id: msg.Id, Topic: msg.Original, Timestamp: &now, Desc: desc,
		},
	})
}

// replyOfflineTopicGetSub reads user's subscription from the database.
// Only own subscription is available.
// The requester must be subscribed but need not be attached.
func replyOfflineTopicGetSub(sess *Session, msg *ClientComMessage) {
	now := types.TimeNow()

	if msg.Get.Sub != nil && msg.Get.Sub.User != "" && msg.Get.Sub.User != msg.AsUser {
		sess.queueOut(ErrPermissionDeniedReply(msg, now))
		return
	}

	topicName := msg.RcptTo
	if types.IsChannel(msg.Original) {
		topicName = msg.Original
	}

	ssub, err := globals.store.SubsGetSubs(topicName, types.ParseUserId(msg.AsUser), true)
	if err != nil {
		globals.l.Sugar().Warnf("replyOfflineTopicGetSub:", err)
		sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
		return
	}

	if ssub == nil {
		sess.queueOut(ErrNotFoundExplicitTs(msg.Id, msg.Original, now, msg.Timestamp))
		return
	}

	sub := MsgTopicSub{}
	if ssub.DeletedAt == nil {
		sub.UpdatedAt = &ssub.UpdatedAt
		sub.Acs = MsgAccessMode{
			Want:  ssub.ModeWant.String(),
			Given: ssub.ModeGiven.String(),
			Mode:  (ssub.ModeGiven & ssub.ModeWant).String(),
		}

		// Fnd is asymmetric: desc.private is a string, but sub.private is a []string.
		if types.GetTopicCat(msg.RcptTo) != types.TopicCatFnd {
			sub.Private = ssub.Private
		}

		sub.User = types.ParseUid(ssub.User).UserId()

		if (ssub.ModeGiven & ssub.ModeWant).IsReader() && (ssub.ModeWant & ssub.ModeGiven).IsJoiner() {
			sub.DelId = ssub.DelId
			sub.ReadSeqId = ssub.ReadSeqId
			sub.RecvSeqId = ssub.RecvSeqId
		}

	} else {
		sub.DeletedAt = ssub.DeletedAt
	}

	sess.queueOut(&ServerComMessage{
		Meta: &MsgServerMeta{
			Id: msg.Id, Topic: msg.Original, Timestamp: &now, Sub: []MsgTopicSub{sub},
		},
	})
}

// replyOfflineTopicSetSub updates Desc.Private and Sub.Mode when the topic is not
//
// loaded in memory. Only Private and Mode are updated and only for the requester. The
//
// requester must be subscribed to the topic but does not need to be attached.
func replyOfflineTopicSetSub(sess *Session, msg *ClientComMessage) {
	now := types.TimeNow()

	if (msg.Set.Desc == nil || msg.Set.Desc.Private == nil) && (msg.Set.Sub == nil || msg.Set.Sub.Mode == "") {
		sess.queueOut(InfoNotModifiedReply(msg, now))
		return
	}

	if msg.Set.Sub != nil && msg.Set.Sub.User != "" && msg.Set.Sub.User != msg.AsUser {
		sess.queueOut(ErrPermissionDeniedReply(msg, now))
		return
	}

	asUid := types.ParseUserId(msg.AsUser)

	topicName := msg.RcptTo
	if types.IsChannel(msg.Original) {
		topicName = msg.Original
	}

	sub, err := globals.store.SubsGetSubs(topicName, asUid, false)
	if err != nil {
		globals.l.Sugar().Warnf("replyOfflineTopicSetSub get sub:", err)
		sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
		return
	}

	if sub == nil {
		sess.queueOut(ErrNotFoundExplicitTs(msg.Id, msg.Original, now, msg.Timestamp))
		return
	}

	update := make(map[string]any)
	if msg.Set.Desc != nil && msg.Set.Desc.Private != nil {
		private, ok := msg.Set.Desc.Private.(map[string]any)
		if !ok {
			update = map[string]any{"Private": msg.Set.Desc.Private}
		} else if private, changed := mergeInterfaces(sub.Private, private); changed {
			update = map[string]any{"Private": private}
		}
	}

	if msg.Set.Sub != nil && msg.Set.Sub.Mode != "" {
		var modeWant types.AccessMode
		if err = modeWant.UnmarshalText([]byte(msg.Set.Sub.Mode)); err != nil {
			globals.l.Sugar().Warnf("replyOfflineTopicSetSub mode:", err)
			sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
			return
		}

		if modeWant.IsOwner() != sub.ModeWant.IsOwner() {
			// No ownership changes here.
			sess.queueOut(ErrPermissionDeniedReply(msg, now))
			return
		}

		if types.GetTopicCat(msg.RcptTo) == types.TopicCatP2P {
			// For P2P topics ignore requests exceeding types.ModeCP2P and do not allow
			// removal of 'A' permission.
			modeWant = modeWant&types.ModeCP2P | types.ModeApprove
		}

		if modeWant != sub.ModeWant {
			update["ModeWant"] = modeWant
			// Cache it for later use
			sub.ModeWant = modeWant
		}
	}

	if len(update) > 0 {
		err = globals.store.SubsUpdate(topicName, asUid, update)
		if err != nil {
			globals.l.Sugar().Warnf("replyOfflineTopicSetSub update:", err)
			sess.queueOut(decodeStoreErrorExplicitTs(err, msg.Id, msg.Original, now, msg.Timestamp, nil))
		} else {

			var params any
			if update["ModeWant"] != nil {
				params = map[string]any{
					"acs": MsgAccessMode{
						Given: sub.ModeGiven.String(),
						Want:  sub.ModeWant.String(),
						Mode:  (sub.ModeGiven & sub.ModeWant).String(),
					},
				}
			}
			sess.queueOut(NoErrParamsReply(msg, now, params))
		}
	} else {
		sess.queueOut(InfoNotModifiedReply(msg, now))
	}
}

// statusChangeBits sets or removes given bits from t.status
func (t *Topic) statusChangeBits(bits int32, set bool) {
	for {
		oldStatus := atomic.LoadInt32(&t.status)
		newStatus := oldStatus
		if set {
			newStatus |= bits
		} else {
			newStatus &= ^bits
		}
		if newStatus == oldStatus {
			break
		}
		if atomic.CompareAndSwapInt32(&t.status, oldStatus, newStatus) {
			break
		}
	}
}

// markLoaded indicates that topic subscribers have been loaded into memory.
func (t *Topic) markLoaded() {
	t.statusChangeBits(topicStatusLoaded, true)
}

// markPaused pauses or unpauses the topic. When the topic is paused all
// messages are rejected.
func (t *Topic) markPaused(pause bool) {
	t.statusChangeBits(topicStatusPaused, pause)
}

// markDeleted marks topic as being deleted.
func (t *Topic) markDeleted() {
	t.statusChangeBits(topicStatusMarkedDeleted, true)
}

// markReadOnly suspends/un-suspends the topic: adds or removes the 'read-only' flag.
func (t *Topic) markReadOnly(readOnly bool) {
	t.statusChangeBits(topicStatusReadOnly, readOnly)
}

// isInactive checks if topic is paused or being deleted.
func (t *Topic) isInactive() bool {
	return (atomic.LoadInt32(&t.status) & (topicStatusPaused | topicStatusMarkedDeleted)) != 0
}

func (t *Topic) isReadOnly() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusReadOnly) != 0
}

func (t *Topic) isLoaded() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusLoaded) != 0
}

func (t *Topic) isDeleted() bool {
	return (atomic.LoadInt32(&t.status) & topicStatusMarkedDeleted) != 0
}

// Get topic name suitable for the given client
func (t *Topic) original(uid types.Uid) string {
	if t.cat == types.TopicCatP2P {
		if pud, ok := t.perUser[uid]; ok {
			return pud.topicName
		}
		panic("Invalid P2P topic")
	}

	if t.cat == types.TopicCatGrp && t.isChan {
		if t.perUser[uid].isChan {
			// This is a channel reader.
			return types.GrpToChn(t.xoriginal)
		}
	}
	return t.xoriginal
}

// Get ID of the other user in a P2P topic
func (t *Topic) p2pOtherUser(uid types.Uid) types.Uid {
	if t.cat == types.TopicCatP2P {
		// Try to find user in subscribers.
		for u2 := range t.perUser {
			if u2.Compare(uid) != 0 {
				return u2
			}
		}
	}

	// Even when one user is deleted, the subscription must be restored
	// before p2pOtherUser is called.
	panic("Not a valid P2P topic")
}

// Get per-session value of fnd.Public
func (t *Topic) fndGetPublic(sess *Session) any {
	if t.cat == types.TopicCatFnd {
		if t.public == nil {
			return nil
		}
		if pubmap, ok := t.public.(map[string]any); ok {
			return pubmap[sess.sid]
		}
		panic("Invalid Fnd.Public type")
	}
	panic("Not Fnd topic")
}

// Assign per-session fnd.Public. Returns true if value has been changed.
func (t *Topic) fndSetPublic(sess *Session, public any) bool {
	if t.cat != types.TopicCatFnd {
		panic("Not Fnd topic")
	}

	var pubmap map[string]any
	var ok bool
	if t.public != nil {
		if pubmap, ok = t.public.(map[string]any); !ok {
			// This could only happen if fnd.public is assigned outside of this function.
			panic("Invalid Fnd.Public type")
		}
	}
	if pubmap == nil {
		pubmap = make(map[string]any)
	}

	if public != nil {
		pubmap[sess.sid] = public
	} else {
		ok = (pubmap[sess.sid] != nil)
		delete(pubmap, sess.sid)
		if len(pubmap) == 0 {
			pubmap = nil
		}
	}
	t.public = pubmap
	return ok
}

// Remove per-session value of fnd.Public.
func (t *Topic) fndRemovePublic(sess *Session) {
	if t.public == nil {
		return
	}
	// FIXME: case of a multiplexing session won't work correctly.
	// Maybe handle it at the proxy topic.
	if pubmap, ok := t.public.(map[string]any); ok {
		delete(pubmap, sess.sid)
		return
	}
	panic("Invalid Fnd.Public type")
}

func (t *Topic) accessFor(authLvl auth.Level) types.AccessMode {
	return selectAccessMode(authLvl, t.accessAnon, t.accessAuth, getDefaultAccess(t.cat, true, false))
}

// subsCount returns the number of topic subscribers
func (t *Topic) subsCount() int {
	if t.cat == types.TopicCatP2P {
		count := 0
		for uid := range t.perUser {
			if !t.perUser[uid].deleted {
				count++
			}
		}
		return count
	}
	return len(t.perUser)
}

// Add session record. 'user' may be different from sess.uid.
func (t *Topic) addSession(sess *Session, asUid types.Uid, isChanSub bool) {
	s := sess
	if sess.multi != nil {
		s = s.multi
	}

	if pssd, ok := t.sessions[s]; ok {
		// Subscription already exists.
		if s.isMultiplex() && !sess.background {
			// This slice is expected to be relatively short.
			// Not doing anything fancy here like maps or sorting.
			pssd.muids = append(pssd.muids, asUid)
			t.sessions[s] = pssd
		}
		// Maybe panic here.
		return
	}

	if s.isMultiplex() {
		if sess.background {
			t.sessions[s] = perSessionData{}
		} else {
			t.sessions[s] = perSessionData{muids: []types.Uid{asUid}, isChanSub: isChanSub}
		}
	} else {
		t.sessions[s] = perSessionData{uid: asUid, isChanSub: isChanSub}
	}
}

// Disconnects session from topic if either one of the following is true:
// * 's' is an ordinary session AND ('asUid' is zero OR 'asUid' matches subscribed user).
// * 's' is a multiplexing session and it's being dropped all together ('asUid' is zero ).
// If 's' is a multiplexing session and asUid is not zero, it's removed from the list of session
// users 'muids'.
// Returns perSessionData if it was found and true if session was actually detached from topic.
func (t *Topic) remSession(sess *Session, asUid types.Uid) (*perSessionData, bool) {
	s := sess
	if sess.multi != nil {
		s = s.multi
	}
	pssd, ok := t.sessions[s]
	if !ok {
		// Session not found at all.
		return nil, false
	}

	if pssd.uid == asUid || asUid.IsZero() {
		delete(t.sessions, s)
		return &pssd, true
	}

	for i := range pssd.muids {
		if pssd.muids[i] == asUid {
			pssd.muids[i] = pssd.muids[len(pssd.muids)-1]
			pssd.muids = pssd.muids[:len(pssd.muids)-1]
			t.sessions[s] = pssd
			if len(pssd.muids) == 0 {
				delete(t.sessions, s)
				return &pssd, true
			}

			return &pssd, false
		}
	}

	return nil, false
}

// Check if topic has any online (non-background) users.
func (t *Topic) isOnline() bool {
	// Find at least one non-background session.
	for s, pssd := range t.sessions {
		if s.isMultiplex() && len(pssd.muids) > 0 {
			return true
		}
		if !s.background {
			return true
		}
	}
	return false
}

// Verifies if topic can be access by the provided name: access any topic as non-channel, access channel as channel.
// Returns true if access is for channel, false if not and error if access is invalid.
func (t *Topic) verifyChannelAccess(asTopic string) (bool, error) {
	if !types.IsChannel(asTopic) {
		return false, nil
	}
	if t.isChan {
		return true, nil
	}
	return false, types.ErrNotFound
}
