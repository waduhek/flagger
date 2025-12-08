# --- Development container image ---
FROM golang:1.25.5 AS dev
WORKDIR /go/github.com/waduhek/flagger
COPY . .
RUN go build -o ./build/flagger ./cmd/flagger
EXPOSE 50051
CMD ["./build/flagger"]

# --- Debug container image ---
FROM golang:1.25.5 AS debug
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /go/github.com/waduhek/flagger
COPY . .
EXPOSE 50051 4040
CMD ["dlv", "debug", "./cmd/flagger", "--headless", "--listen", ":4040"]
