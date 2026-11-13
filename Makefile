APP_NAME   := server
CMD_DIR    := ./cmd/server
BIN_DIR    := ./bin
BIN        := $(BIN_DIR)/$(APP_NAME)

.PHONY: all build run fmt lint clean

all: build

## fmt: format all Go code
fmt:
	go fmt ./...

## lint: run go vet and golangci-lint
lint:
	go vet ./...
	golangci-lint run ./...

## build: format, lint, then compile cmd/server into ./bin/server
build: fmt lint
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD_DIR)

## run: format, lint, then run cmd/server (pass extra args with ARGS="...")
run: fmt lint
	go run $(CMD_DIR) $(ARGS)

## clean: remove build artifacts
clean:
	rm -rf $(BIN_DIR)
