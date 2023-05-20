package agent_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	api "github.com/Ahmadkhatib0/go/distributed-services/proglog/api/v1"
	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/agent"
	"github.com/Ahmadkhatib0/go/distributed-services/proglog/internal/config"
	"github.com/stretchr/testify/require"
	"github.com/travisjeffery/go-dynaport"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// TestAgent(t *testing.T) sets up a three-node cluster. The second and third nodes join the first node’s cluster
func TestAgent(t *testing.T) {

	// The serverTLSConfig defines the configuration of the certificate that’s served to clients
	serverTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.ServerCertFile,
		KeyFile:       config.ServerKeyFile,
		CAFile:        config.CAFile,
		Server:        true,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	// the peerTLSConfig defines the configuration of the certificate that’s
	// served between servers so they can connect with and replicate each other.
	peerTLSConfig, err := config.SetupTLSConfig(config.TLSConfig{
		CertFile:      config.RootClientCertFile,
		KeyFile:       config.RootClientKeyFile,
		CAFile:        config.CAFile,
		Server:        false,
		ServerAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	var agents []*agent.Agent
	for i := 0; i < 3; i++ {

		//  +-----------------------------------------------------------------------------------------------------------+
		//  |	Because we now have two addresses to configure in our service (the RPC address and the Serf address),  |
		//  |	and because we run our tests on a SINGLE HOST, we need two ports, We used the 0 port trick to test     |
		//  |	a gRPC Server and Client, to get a port automatically assigned to a listener by net.Listen,            |
		//  |	(but now we just want the port WITH NO LISTENER), so we use the dynaport library to allocate the two   |
		//  |	ports we need: one for our gRPC log connections and one for our Serf service discovery connections     |
		//  +-----------------------------------------------------------------------------------------------------------+
		ports := dynaport.Get(2)
		bindAddr := fmt.Sprintf("%s:%d", "127.0.0.1", ports[0])
		rpcPort := ports[1]
		dataDir, err := ioutil.TempDir("", "agent-test-log")
		require.NoError(t, err)

		var startJoinAddrs []string
		if i != 0 {
			startJoinAddrs = append(startJoinAddrs, agents[0].Config.BindAddr)
		}

		agent, err := agent.New(agent.Config{
			NodeName:        fmt.Sprintf("%d", i),
			StartJoinAddrs:  startJoinAddrs,
			BindAddr:        bindAddr,
			RPCPort:         rpcPort,
			DataDir:         dataDir,
			ACLModelFile:    config.ACLModelFile,
			ACLPolicyFile:   config.ACLPolicyFile,
			ServerTLSConfig: serverTLSConfig,
			PeerTLSConfig:   peerTLSConfig,
		})
		require.NoError(t, err)

		agents = append(agents, agent)
	}

	defer func() {
		for _, agent := range agents {
			err := agent.Shutdown()
			require.NoError(t, err)
			require.NoError(t, os.RemoveAll(agent.Config.DataDir)) // delete all test data
		}
	}()

	// make the test sleep for a few seconds to give the nodes time to discover each other
	time.Sleep(3 * time.Second)

	// it checks that we can produce to and consume from a single node.
	leaderClient := client(t, agents[0], peerTLSConfig)

	produceResponse, err := leaderClient.Produce(
		context.Background(),
		&api.ProduceRequest{Record: &api.Record{Value: []byte("foo")}},
	)
	require.NoError(t, err)

	consumeResponse, err := leaderClient.Consume(
		context.Background(),
		&api.ConsumeRequest{Offset: produceResponse.Offset},
	)
	require.NoError(t, err)

	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))

	// wait until replication has finished
	// 	Because our replication works asynchronously across servers, the logs pro-
	// duced to one server won’t be immediately available on the replica servers.
	time.Sleep(3 * time.Second)

	followerClient := client(t, agents[1], peerTLSConfig)
	consumeResponse, err = followerClient.Consume(
		context.Background(),
		&api.ConsumeRequest{Offset: produceResponse.Offset},
	)
	require.NoError(t, err)
	require.Equal(t, consumeResponse.Record.Value, []byte("foo"))

}

func client(t *testing.T, agent *agent.Agent, tlsConfig *tls.Config) api.LogClient {
	tlsCreds := credentials.NewTLS(tlsConfig)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(tlsCreds)}
	rpcAddr, err := agent.Config.RPCAddr()
	require.NoError(t, err)

	conn, err := grpc.Dial(fmt.Sprintf("%s", rpcAddr), opts...)
	require.NoError(t, err)

	client := api.NewLogClient(conn)
	return client
}
