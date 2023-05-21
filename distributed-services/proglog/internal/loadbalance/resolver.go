package loadbalance

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"

	api "github.com/Ahmadkhatib0/go/distributed-services/proglog/api/v1"
)

//  +--------------------------------------------------------------------------------+
//  | gRPC uses the builder pattern for resolvers and pickers, so each has a builder |
//  | interface and an implementation interface. Because the builder interfaces have |
//  | one simple method—Build()—we’ll implement both interfaces with one type.       |
//  +--------------------------------------------------------------------------------+

// Resolver is the type we’ll implement into gRPC’s resolver.Builder and resolver.Resolver interfaces
type Resolver struct {
	mu            sync.Mutex
	clientConn    resolver.ClientConn // resolverConn is the resolver’s own client connection to the server
	resolverConn  *grpc.ClientConn    // clientConn connection is the user’s client connection and gRPC passes it to the resolver for the resolver to update with the servers it discovers.
	serviceConfig *serviceconfig.ParseResult
	logger        *zap.Logger
}

const Name = "proglog"

var _ resolver.Builder = (*Resolver)(nil)

func init() {
	resolver.Register(&Resolver{})
}

func (r *Resolver) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {

	r.logger = zap.L().Named("resolver")
	r.clientConn = cc

	var dialOpts []grpc.DialOption
	if opts.DialCreds != nil {
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(opts.DialCreds))
	}

	r.serviceConfig = r.clientConn.ParseServiceConfig(fmt.Sprintf(`{"loadBalancingConfig":[{"%s":{}}]}`, Name))

	var err error

	r.resolverConn, err = grpc.Dial(target.Endpoint, dialOpts...)
	if err != nil {
		return nil, err
	}

	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (r *Resolver) Scheme() string {
	return Name
}

var _ resolver.Resolver = (*Resolver)(nil)

// ResolveNow() gRPC calls it to resolve the target, discover the servers,
//
// and update the client  connection with the servers.
func (r *Resolver) ResolveNow(resolver.ResolveNowOptions) {
	// gRPC may call ResolveNow() concurrently, so we use a mutex to protect access across goroutines.
	r.mu.Lock()
	defer r.mu.Unlock()

	client := api.NewLogClient(r.resolverConn) // get cluster and then set on cc attributes
	ctx := context.Background()
	res, err := client.GetServers(ctx, &api.GetServersRequest{})

	if err != nil {
		r.logger.Error("failed to resolve server", zap.Error(err))
		return
	}

	var addrs []resolver.Address
	for _, server := range res.Servers {
		addrs = append(addrs, resolver.Address{
			Addr:       server.RpcAddr,
			Attributes: attributes.New("is_leader", server.IsLeader),
		})
	}

	r.clientConn.UpdateState(resolver.State{Addresses: addrs, ServiceConfig: r.serviceConfig})
}

// Close() closes the resolver. In our resolver, we close the connection to our server created in Build()
func (r *Resolver) Close() {

	if err := r.resolverConn.Close(); err != nil {
		r.logger.Error("failed to close conn", zap.Error(err))
	}
}
