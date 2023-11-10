# Compiles the protobufs and generates the stubs
.PHONY: proto
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/**/*.proto

# Lint the code
.PHONY: vet
vet:
	go vet ./...

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

# --- Targets for running ---

# Run debug image
.PHONY: run-debug
run-debug:
	docker compose up flagger-debug

# Run development image
.PHONY: run-dev
run-dev:
	docker compose up flagger-dev

# Run other containers
.PHONY: run-others
run-others:
	docker compose up mongo redis

# --- Testing and benchmarking targets ---
.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. -run=^# -count=5 -benchmem ./...

# --- Targets for tearing down containers ---

# Teardown everything
.PHONY: teardown-all
teardown-all:
	docker compose down

# Teardown other containers
.PHONY: teardown-others
teardown-others:
	docker compose down mongo redis

# Teardown dev container
.PHONY: teardown-dev
teardown-dev:
	docker compose down flagger-dev

# Teardown debug container
.PHONY: teardown-debug
teardown-debug:
	docker compose down flagger-debug
