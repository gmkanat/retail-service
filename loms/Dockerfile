FROM golang:1.22.5-alpine as builder
LABEL authors="blackm1nd"
RUN apk add --no-cache git make protoc

WORKDIR /loms

RUN go env -w GOCACHE=/go-cache
RUN go env -w GOMODCACHE=/gomod-cache


COPY loms/go.mod loms/go.sum ./
RUN --mount=type=cache,target=/gomod-cache go mod download

COPY api api
COPY loms/internal/config/stock-data.json ./
COPY api/openapiv2/loms.swagger.json ./
COPY . .


RUN  --mount=type=cache,target=/gomod-cache --mount=type=cache,target=/go-cache CGO_ENABLED=0 GOOS=linux go build -o /server loms/cmd/server/server.go

FROM alpine:latest
RUN apk add --no-cache bash
COPY --from=builder server /bin/server
COPY --from=builder /loms/stock-data.json /loms/stock-data.json
COPY --from=builder /loms/loms.swagger.json /loms/loms.swagger.json

EXPOSE 50051
EXPOSE 8088
ENTRYPOINT ["/bin/server"]
