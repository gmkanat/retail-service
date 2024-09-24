FROM golang:1.22.5-alpine as builder
LABEL authors="blackm1nd"

WORKDIR /cart

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache

COPY cart/go.mod cart/go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download

COPY api api
COPY . .


RUN  --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache CGO_ENABLED=0 GOOS=linux go build -o /server cart/cmd/server/server.go

FROM alpine:latest
RUN apk add --no-cache bash

COPY --from=builder /server /bin/server
EXPOSE 8002
ENTRYPOINT ["/bin/server"]
