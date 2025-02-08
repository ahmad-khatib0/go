FROM golang:1.23.6-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY . .

RUN go mod download 
RUN go build ./cmd

EXPOSE ${BOOKSWAP_PORT}

CMD [ "./cmd" ]

