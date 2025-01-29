//go:build tools
// +build tools

package tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/rakyll/statik"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

// What this code does is pretty simple, It’s just a list of blank imports of the protoc plugins.
// The reason we’re doing this is because we’re not using them directly in the code, but we just
// want to install them to our local machine, so that protoc can use them to generate codes for us.
// and they are available in the mod file
