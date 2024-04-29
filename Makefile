# Compiles the protobufs and generates the stubs
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto

# Lint the code
.PHONY: lint
lint:
	docker run --rm -v $$(pwd):/app -w /app golangci/golangci-lint golangci-lint run

# Run gofmt
.PHONY: format
format:
	gofmt -w -l \
		$$(find . -path ./_docker -prune -o -type f -name '*.go' -printf "%p ")

# --- Targets for building ---

# Build debug Docker image
.PHONY: build-debug
build-debug:
	docker build --target debug -t waduhek/flagger:debug .

# Build development Docker image
.PHONY: build-dev
build-dev:
	docker build --target dev -t waduhek/flagger:dev .

# --- Testing and benchmarking targets ---
.PHONY: test
test:
	go test ./...

.PHONY: integration-test
integration-test:
	./scripts/run-integration-tests.sh

.PHONY: test-all
test-all: test integration-test

.PHONY: bench
bench:
	go test -bench=. -run=^# -count=5 -benchmem ./...
