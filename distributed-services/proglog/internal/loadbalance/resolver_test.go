package loadbalance_test

import (
	"net"
	"testing"

	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/config"
	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/loadbalance"
	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/server"

	api "github.com/Ahmadkhatib0/go/distributed-services/proglog/api/v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

func TestResolver(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	// when listening on port 0, port number is automatically chosen. The Addr method of
	// Listener can be used to discover the chosen port.
	require.NoError(t, err)

	tlsConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})

	require.NoError(t, err)
	serverCreds := credentials.NewTLS(tlsConfig)

	// mock GetServerers so we can set what servers the resolver should find.
	srv, err := server.NewGRPCServer(&server.Config{GetServerer: &getServers{}}, grpc.Creds(serverCreds))
	require.NoError(t, err)

	go srv.Serve(l)

	conn := &clientConn{}

	tlsConfig, err = config.SetupTLSConfig(
		config.TLSConfig{
			CertFile:      config.RootClientCertFile,
			KeyFile:       config.RootClientKeyFile,
			CAFile:        config.CAFile,
			Server:        false,
			ServerAddress: "127.0.0.1",
		})
	require.NoError(t, err)

	clientCreds := credentials.NewTLS(tlsConfig)
	opts := resolver.BuildOptions{DialCreds: clientCreds}
	r := &loadbalance.Resolver{}

	// configures the target end-point to point to the server we set up in first step
	_, err = r.Build(resolver.Target{Endpoint: l.Addr().String()}, conn, opts)
	require.NoError(t, err)

	wantState := resolver.State{
		Addresses: []resolver.Address{{ // client will connect to those servers
			Addr:       "localhost:9001",
			Attributes: attributes.New("is_leader", true),
		}, {
			Addr:       "localhost:9002",
			Attributes: attributes.New("is_leader", false),
		}},
	}
	require.Equal(t, wantState, conn.state)

	conn.state.Addresses = nil

	r.ResolveNow(resolver.ResolveNowOptions{})
	require.Equal(t, wantState, conn.state)
}

type getServers struct{}

func (s *getServers) GetServers() ([]*api.Server, error) {
	return []*api.Server{{
		Id:       "leader",
		RpcAddr:  "localhost:9001",
		IsLeader: true,
	}, {
		Id:      "follower",
		RpcAddr: "localhost:9002",
	}}, nil
}

// clientConn implements resolver.ClientConn, and its job is to keep a reference to the state the resolver
// updated it with so that we can verify that the resolver updates the client connection with the correct data
type clientConn struct {
	resolver.ClientConn
	state resolver.State
}

func (c *clientConn) UpdateState(state resolver.State) error {
	c.state = state
	return nil
}

func (c *clientConn) ReportError(err error) {}

func (c *clientConn) NewAddress(addrs []resolver.Address) {}

func (c *clientConn) NewServiceConfig(config string) {}

func (c *clientConn) ParseServiceConfig(config string) *serviceconfig.ParseResult {
	return nil
}
