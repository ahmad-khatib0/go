package agent

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/ahmad-khatib0/go/distributed-services/proglog/internal/auth"
	"github.com/ahmad-khatib0/go/distributed-services/proglog/internal/discovery"
	"github.com/ahmad-khatib0/go/distributed-services/proglog/internal/log"
	"github.com/ahmad-khatib0/go/distributed-services/proglog/internal/server"
	"github.com/hashicorp/raft"
	"github.com/soheilhy/cmux"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Agent runs on every service instance, setting up and connecting all the different components.
type Agent struct {
	Config
	mux          cmux.CMux
	log          *log.DistributedLog
	server       *grpc.Server
	membership   *discovery.Membership
	shutdown     bool
	shutdowns    chan struct{}
	shutdownLock sync.Mutex
}

type Config struct {
	Bootstrap       bool
	ServerTLSConfig *tls.Config
	PeerTLSConfig   *tls.Config
	DataDir         string
	BindAddr        string
	RPCPort         int
	NodeName        string
	StartJoinAddrs  []string
	ACLModelFile    string
	ACLPolicyFile   string
}

func (c Config) RPCAddr() (string, error) {
	host, _, err := net.SplitHostPort(c.BindAddr)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s:%d", host, c.RPCPort), nil
}

func New(config Config) (*Agent, error) {
	a := &Agent{
		Config:    config,
		shutdowns: make(chan struct{}),
	}

	setup := []func() error{
		a.setupLogger,
		a.setupMux,
		a.setupLog,
		a.setupServer,
		a.setupMembership,
	}

	for _, fn := range setup {
		if err := fn(); err != nil {
			return nil, err
		}
	}

	go a.serve()
	return a, nil
}

func (a *Agent) setupLogger() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(logger)
	return nil
}

func (a *Agent) setupLog() error {

	raftLn := a.mux.Match(func(reader io.Reader) bool {
		b := make([]byte, 1)
		if _, err := reader.Read(b); err != nil {
			return false
		}
		return bytes.Compare(b, []byte{byte(log.RaftRPC)}) == 0
	})
	// We identify Raft connections by reading one byte and checking that the byte
	// matches the byte we set up our outgoing Raft connections to write

	//  configure the distributed log’s Raft to use our multiplexed listener
	logConfig := log.Config{}
	logConfig.Raft.StreamLayer = log.NewStreamLayer(raftLn, a.Config.ServerTLSConfig, a.Config.PeerTLSConfig)

	rpcAddr, err := a.Config.RPCAddr()
	if err != nil {
		return err
	}
	logConfig.Raft.BindAddr = rpcAddr
	logConfig.Raft.LocalID = raft.ServerID(a.Config.NodeName)
	logConfig.Raft.Bootstrap = a.Config.Bootstrap

	// configure and create the distributed log
	a.log, err = log.NewDistributedLog(a.Config.DataDir, logConfig)
	if err != nil {
		return err
	}

	if a.Config.Bootstrap {
		err = a.log.WaitForLeader(3 * time.Second)
	}

	return err
}

func (a *Agent) setupServer() error {
	authorizer := auth.New(a.Config.ACLModelFile, a.Config.ACLPolicyFile)
	serverConfig := &server.Config{
		CommitLog:   a.log,
		Authorizer:  authorizer,
		GetServerer: a.log,
	}

	var opts []grpc.ServerOption
	if a.Config.ServerTLSConfig != nil {
		creds := credentials.NewTLS(a.Config.ServerTLSConfig)
		opts = append(opts, grpc.Creds(creds))
	}

	var err error
	a.server, err = server.NewGRPCServer(serverConfig, opts...)
	if err != nil {
		return err
	}

	// Because we’ve multiplexed two connection types (Raft and gRPC) and we added a matcher for the
	// Raft connections, we know all other connections must be gRPC connections. We use cmux.Any()
	// because it matches any connections
	grpcLn := a.mux.Match(cmux.Any())
	go func() {
		// we tell our gRPC server to serve on the multiplexed listener.
		if err := a.server.Serve(grpcLn); err != nil {
			_ = a.Shutdown()
		}
	}()

	return err
}

func (a *Agent) setupMembership() error {
	rpcAddr, err := a.Config.RPCAddr()
	if err != nil {
		return err
	}

	a.membership, err = discovery.New(a.log, discovery.Config{
		NodeName:       a.Config.NodeName,
		BindAddr:       a.Config.BindAddr,
		Tags:           map[string]string{"rpc_addr": rpcAddr},
		StartJoinAddrs: a.Config.StartJoinAddrs,
	})

	return err
}

func (a *Agent) Shutdown() error {
	a.shutdownLock.Lock()
	defer a.shutdownLock.Unlock()

	// ensures that the agent will shut down once even if people call Shutdown() multiple times
	if a.shutdown {
		return nil
	}
	a.shutdown = true

	close(a.shutdowns)

	shutdown := []func() error{
		// Leaving the membership so that other servers will see that this server has left the cluster
		// and so that this server doesn’t receive discovery events anymore;
		a.membership.Leave,

		// Closing the replicator so it doesn’t continue to replicate;
		// a.replicator.Close,

		func() error {
			// Gracefully stopping the server, which stops the server from accepting new connections and
			// blocks until all the pending RPCs have finished;
			a.server.GracefulStop()
			return nil
		},

		a.log.Close,
	}

	for _, fn := range shutdown {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

// setupMux() creates a listener on our RPC address that’ll accept both Raft and
//
// gRPC connections and then creates the mux with the listener. The mux will accept
//
// connections on that listener and match connections based on your configured rules.
func (a *Agent) setupMux() error {
	rpcAddr := fmt.Sprintf(":%d", a.Config.RPCPort)
	ln, err := net.Listen("tcp", rpcAddr)

	if err != nil {
		return err
	}

	a.mux = cmux.New(ln)
	return nil
}

func (a *Agent) serve() error {
	if err := a.mux.Serve(); err != nil {
		_ = a.Shutdown()
		return err
	}

	return nil
}
