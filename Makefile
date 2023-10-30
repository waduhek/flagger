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

# Build debug Docker image
.PHONY: build-debug
build-debug:
	docker build --target debug -t waduhek/flagger:debug .

# Run debug image
.PHONY: run-debug
run-debug:
	docker run --env-file ./.env \
		-p 50051:50051 -p 4040:4040 \
		--name flagger-debug \
		waduhek/flagger:debug

# Build development Docker image
.PHONY: build-dev
build-dev:
	docker build --target dev -t waduhek/flagger:dev .

# Run development image
run-dev:
	docker run --env-file ./.env \
		-p 50051:50051 \
		--name flagger-dev \
		waduhek/flagger:dev

# --- Testing and benchmarking targets ---
.PHONY: test
test:
	go test ./...

.PHONY: bench
bench:
	go test -bench=. -run=^# -count=5 -benchmem ./...
