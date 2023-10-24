# --- Development container image ---
FROM golang:1.21.3 as dev
WORKDIR /go/github.com/waduhek/flagger
COPY . .
EXPOSE 50051
CMD ["go", "run", "./cmd/flagger"]

# --- Debug container image ---
FROM golang:1.21.3 as debug
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /go/github.com/waduhek/flagger
COPY . .
EXPOSE 50051 4040
CMD ["dlv", "debug", "./cmd/flagger", "--headless", "--listen", ":4040"]
