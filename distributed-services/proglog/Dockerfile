# uses the golang:1.19-alpine image because we need the Go compiler,
# we split into tow steps because we don't need to clutter the container with dependencies, and perhaps 
# various system libraries These take up disk space, and we don’t need them after we have compiled our binary
FROM golang:1.22 AS build

WORKDIR /go/src/proglog

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog

# install the grpc_health_probe executable
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.2 && \ 
  wget -qO/go/bin/grpc_health_probe \ 
  https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \ 
  chmod +x /go/bin/grpc_health_probe

# the scratch empty image—the smallest Docker image.
FROM scratch

COPY --from=build /go/bin/proglog /bin/proglog

COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe

ENTRYPOINT ["/bin/proglog"]
