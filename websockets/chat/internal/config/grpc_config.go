package config

type GrpcConfig struct {
	// Address:port to listen for gRPC clients
	Listen string `json:"listen"`
	// Enable handling of gRPC keepalives https://github.com/grpc/grpc/blob/master/doc/keepalive.md
	// This sets server's GRPC_ARG_KEEPALIVE_TIME_MS to 60 seconds instead of the default 2 hours.
	Keepalive bool `json:"keepalive"`
}
