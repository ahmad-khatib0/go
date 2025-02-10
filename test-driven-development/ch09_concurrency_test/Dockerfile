# syntax=docker/dockerfile:1

FROM golang:1.23.6-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .

RUN go mod download
RUN go build ./chapter09/cmd
RUN apk --no-cache add curl
EXPOSE ${BOOKSWAP_PORT}

CMD [ "./cmd" ]
