.PHONY: lint
lint: lint-golangci

.PHONY: lint-golangci
lint-golangci:
	golangci-lint run

.PHONY: test
test:
	go test -v ./...

.PHONY: build
build:
	@go build -o bin/hexlet-path-size ./cmd/hexlet-path-size