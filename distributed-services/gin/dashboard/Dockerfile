FROM golang:1.22

WORKDIR /go/src/github.com/api

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=0 /go/src/github.com/api/app .
COPY templates templates/
COPY assets assets/

CMD ["./app"]  
