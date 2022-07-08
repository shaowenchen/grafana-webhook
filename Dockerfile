ARG GO_VERSION=1.17

FROM golang:${GO_VERSION}-alpine AS builder
RUN apk update && apk add alpine-sdk git
RUN mkdir -p /builder
WORKDIR /builder
COPY . .
RUN make binary

FROM alpine:latest
WORKDIR /
COPY --from=builder /builder/bin/app .
COPY --from=builder /builder/default.toml .
CMD [ "/app", "-c", "./default.toml"]
EXPOSE 8000
