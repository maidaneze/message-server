FROM golang:1.12

ADD ./ /go/src/github.com/maidaneze/message-server
WORKDIR /go/src/github.com/maidaneze/message-server

ENV GO111MODULE=on

RUN apt-get update && apt-get install -y --no-install-recommends \
    sqlite3

EXPOSE 8080

CMD ["go", "run", "server.go"]