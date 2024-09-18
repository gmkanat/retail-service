CYCLO_LIMIT = 10
COGNIT_LIMIT = 15

test-cover:
	# ignore mock files
	go test $(shell go list ./internal/pkg/cart/... | grep -v ./internal/pkg/cart/mocks) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

test-e2e:
	go test ./e2e/... -v

cyclo:
	gocyclo -over $(CYCLO_LIMIT) -ignore "_mock.go" ./

cognit:
	gocognit -over $(COGNIT_LIMIT) -ignore "_mock.go" ./

lint: cyclo cognit

bench:
	go test ./internal/pkg/cart/repository -bench=. -benchmem
