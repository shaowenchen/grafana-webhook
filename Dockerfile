ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git && rm -rf /var/cache/apk/*

RUN mkdir -p /builder
WORKDIR /builder

COPY . .
RUN make binary

FROM alpine:latest
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*

WORKDIR /
COPY --from=builder /builder/bin/app .
COPY --from=builder /builder/conf/run.toml .
CMD [ "/app", "-c", "./run.toml"]
EXPOSE 8000
