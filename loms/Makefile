CYCLO_LIMIT = 10
COGNIT_LIMIT = 15

test-cover:
	go test $(shell go list ./internal/service/... ) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-e2e:
	go test ./e2e/... -v

install-gocyclo:
	go install github.com/fzipp/gocyclo/cmd/gocyclo@latest

install-gocognit:
	go install github.com/uudashr/gocognit/cmd/gocognit@latest

cyclo:
	gocyclo -over $(CYCLO_LIMIT) -ignore "_mock.go" ./

cognit:
	gocognit -over $(COGNIT_LIMIT) -ignore "_mock.go" ./

lint: install-gocyclo install-gocognit cyclo cognit
