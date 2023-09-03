FROM alpine:latest

RUN apk update && apk upgrade

RUN apk add bash

COPY base.yaml .

COPY main .

EXPOSE 8081

CMD ["./main"]
