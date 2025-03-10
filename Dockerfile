# --- Development container image ---
FROM golang:1.24.1 as dev
WORKDIR /go/github.com/waduhek/flagger
COPY . .
RUN go build -o ./build/flagger ./cmd/flagger
EXPOSE 50051
CMD ["./build/flagger"]

# --- Debug container image ---
FROM golang:1.24.1 as debug
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /go/github.com/waduhek/flagger
COPY . .
EXPOSE 50051 4040
CMD ["dlv", "debug", "./cmd/flagger", "--headless", "--listen", ":4040"]
