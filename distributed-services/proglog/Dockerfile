# uses the golang:1.19-alpine image because we need the Go compiler,
# we split into tow steps because we don't need to clutter the container with dependencies, and perhaps 
# various system libraries These take up disk space, and we don’t need them after we have compiled our binary
FROM golang:1.19-alpine AS build

WORKDIR /go/src/proglog

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=0 go build -o /go/bin/proglog ./cmd/proglog

# the scratch empty image—the smallest Docker image.
FROM scratch

COPY --from=build /go/bin/proglog /bin/proglog

ENTRYPOINT ["/bin/proglog"]
