package server

import (
	"crypto/tls"
	"io"
	"time"

	"github.com/ahmad-khatib0/go/websockets/chat-protobuf/chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/peer"
)

type grpcNodeServer struct {
	chat.UnimplementedNodeServer
}

func (sess *Session) closeGrpc() {
	if sess.proto == GRPC {
		sess.lock.Lock()
		sess.grpcnode = nil
		sess.lock.Unlock()
	}
}

// Equivalent of starting a new session and a read loop in one.
func (*grpcNodeServer) MessageLoop(stream chat.Node_MessageLoopServer) error {
	sess, count := globals.sessionStore.NewSession(stream, "")
	if p, ok := peer.FromContext(stream.Context()); ok {
		sess.remoteAddr = p.Addr.String()
	}

	globals.l.Sugar().Infof("grpc: session started", sess.sid, sess.remoteAddr, count)
	defer func() {
		sess.closeGrpc()
		sess.cleanUp(false)
	}()

	go sess.writeGrpcLoop()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			globals.l.Sugar().Errorf("grpc: recv", sess.sid, err)
			return err
		}

		globals.l.Sugar().Infof("grpc in:", truncateStringIfTooLong(in.String()), sess.sid)
		globals.stats.IntStatsInc("IncomingMessagesGrpcTotal", 1)
		sess.dispatch(pbCliDeserialize(in))

		sess.lock.Lock()
		if sess.grpcnode == nil {
			sess.lock.Unlock()
			break
		}
		sess.lock.Unlock()
	}

	return nil
}

func (sess *Session) sendMessageGrpc(msg any) bool {
	if len(sess.send) > sendQueueLimit {
		globals.l.Sugar().Errorf("grpc: outbound queue limit exceeded", sess.sid)
		return false
	}

	globals.stats.IntStatsInc("OutgoingMessagesGrpcTotal", 1)
	if err := grpcWrite(sess, msg); err != nil {
		globals.l.Sugar().Errorf("grpc: write", sess.sid, err)
		return false
	}
	return true
}

func (sess *Session) writeGrpcLoop() {
	defer func() {
		sess.closeGrpc() // exit MessageLoop
	}()

	for {
		select {
		case msg, ok := <-sess.send:
			if !ok {
				// channel closed
				return
			}
			switch v := msg.(type) {
			case []*ServerComMessage: // batch of unserialized messages
				for _, msg := range v {
					w := sess.serializeAndUpdateStats(msg)
					if !sess.sendMessageGrpc(w) {
						return
					}
				}
			case *ServerComMessage: // single unserialized message
				w := sess.serializeAndUpdateStats(v)
				if !sess.sendMessageGrpc(w) {
					return
				}
			default: // serialized message
				if !sess.sendMessageGrpc(v) {
					return
				}
			}

		case <-sess.bkgTimer.C:
			if sess.background {
				sess.background = false
				sess.onBackgroundTimer()
			}

		case msg := <-sess.stop:
			// Shutdown requested, don't care if the message is delivered
			if msg != nil {
				grpcWrite(sess, msg)
			}
			return

		case topic := <-sess.detach:
			sess.delSub(topic)
		}
	}
}

func grpcWrite(sess *Session, msg any) error {
	if out := sess.grpcnode; out != nil {
		// Will panic if msg is not of *pbx.ServerMsg type. This is an intentional panic.
		return out.Send(msg.(*chat.ServerMsg))
	}
	return nil
}

func serveGrpc(addr string, kaEnabled bool, tlsConf *tls.Config) (*grpc.Server, error) {
	if addr == "" {
		return nil, nil
	}

	lis, err := netListener(addr)
	if err != nil {
		return nil, err
	}

	secure := ""
	var opts []grpc.ServerOption
	opts = append(opts, grpc.MaxRecvMsgSize(int(globals.maxMessageSize)))
	if tlsConf != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(tlsConf)))
		secure = " secure"
	}

	if kaEnabled {
		kepConfig := keepalive.EnforcementPolicy{
			MinTime:             1 * time.Second, // If a client pings more than once every second, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}
		opts = append(opts, grpc.KeepaliveEnforcementPolicy(kepConfig))

		kpConfig := keepalive.ServerParameters{
			Time:    60 * time.Second, // Ping the client if it is idle for 60 seconds to ensure the connection is still active
			Timeout: 20 * time.Second, // Wait 20 second for the ping ack before assuming the connection is dead
		}
		opts = append(opts, grpc.KeepaliveParams(kpConfig))
	}

	srv := grpc.NewServer(opts...)
	chat.RegisterNodeServer(srv, &grpcNodeServer{})
	globals.l.Sugar().Infof("gRPC/%s%s server is registered at [%s]", grpc.Version, secure, addr)

	go func() {
		if err := srv.Serve(lis); err != nil {
			globals.l.Sugar().Errorf("gRPC server failed:", err)
		}
	}()

	return srv, nil
}
