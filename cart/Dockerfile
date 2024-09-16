FROM golang:1.22.5-alpine as builder
LABEL authors="blackm1nd"

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server/server.go

FROM scratch
COPY --from=builder server /bin/server
EXPOSE 8002
ENTRYPOINT ["/bin/server"]
