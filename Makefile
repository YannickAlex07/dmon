.PHONY: test
test:
	go test ./...

.PHONY: coverage
coverage:
	go test ./pkg/... -coverpkg=./... -coverprofile coverage.out
	go tool cover -html=coverage.out