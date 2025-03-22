FROM golang:1.24-bookworm

WORKDIR /app

ENV CGO_ENABLED=1

RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      build-essential \
      libsqlite3-dev

COPY . .

RUN go mod tidy


RUN go build -o vpn_server main.go

FROM debian:bookworm
RUN apt-get update \
 && DEBIAN_FRONTEND=noninteractive \
    apt-get install --no-install-recommends --assume-yes \
      ca-certificates \
      libsqlite3-0 \
      docker.io

WORKDIR /root/

COPY --from=0 /app/vpn_server .

CMD ["/root/vpn_server"]