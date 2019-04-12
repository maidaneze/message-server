FROM golang:1.12 AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    sqlite3

WORKDIR /server

COPY . ./

# Building using -mod=vendor, which will utilize the v
RUN CGO_ENABLED=1 GOOS=linux go build -mod=vendor -o app

EXPOSE 8080

CMD ["./app"]
