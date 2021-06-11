FROM golang:1.16.5-alpine  as build

LABEL maintainer "github.com/jxsl13"

RUN apk --update add git openssh && \
    rm -rf /var/lib/apt/lists/* && \
    rm /var/cache/apk/*

WORKDIR /build
COPY . ./
COPY go.* ./

ENV CGO_ENABLED=0
ENV GOOS=linux 

RUN go get -d && go build -a -ldflags '-w -extldflags "-static"' -o tw-server-status-discord-bot .


FROM alpine:latest as minimal

ENV DISCORD_TOKEN=""
ENV DISCORD_CHANNEL_ID=""
ENV DICORD_OWNER=""
ENV TEEWORLDS_SERVERS=""
ENV REFRESH_INTERVAL="60s"
ENV CUSTOM_FLAGS=""

WORKDIR /app
COPY --from=build /build/tw-server-status-discord-bot .
VOLUME ["/data", "/app/.env"]
ENTRYPOINT ["/app/tw-server-status-discord-bot"]