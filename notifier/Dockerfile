FROM golang:1.22.5-alpine as builder
LABEL authors="blackm1nd"

WORKDIR /notifier

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY notifier/go.mod notifier/go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download

COPY api api
COPY . .


RUN  --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache CGO_ENABLED=0 GOOS=linux go build -o /main notifier/cmd/notifier/main.go

FROM alpine:latest
RUN apk add --no-cache bash

COPY --from=builder /main /bin/main

ENTRYPOINT ["/bin/main"]
