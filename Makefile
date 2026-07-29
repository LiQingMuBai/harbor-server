.DEFAULT_GOAL := help

GO ?= go
BIN_DIR ?= bin
BIN_PREFIX ?= harbor

API_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-api
ADMIN_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-admin
WSS_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-wss
CDN_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-cdn
TASK_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-task
TASK_DATA_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-task-data
TASK_JOBS_BIN ?= $(BIN_DIR)/$(BIN_PREFIX)-task-jobs

GOFLAGS ?= -trimpath
LDFLAGS ?= -s -w

.PHONY: help
help:
	@printf "%s\n" \
	"Targets:" \
	"  make build            Build all cmd binaries into ./$(BIN_DIR)" \
	"  make build BIN_PREFIX=cointrade   Build with custom binary prefix" \
	"  make build-api        Build $(API_BIN)" \
	"  make build-admin      Build $(ADMIN_BIN)" \
	"  make build-wss        Build $(WSS_BIN)" \
	"  make build-cdn        Build $(CDN_BIN)" \
	"  make build-task       Build $(TASK_BIN)" \
	"  make build-task-data  Build $(TASK_DATA_BIN)" \
	"  make build-task-jobs  Build $(TASK_JOBS_BIN)" \
	"  make clean            Remove ./$(BIN_DIR)"

.PHONY: build
build: build-api build-admin build-wss build-cdn build-task build-task-data build-task-jobs

$(BIN_DIR):
	@mkdir -p $(BIN_DIR)

.PHONY: build-api
build-api: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(API_BIN) ./cmd/api

.PHONY: build-admin
build-admin: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(ADMIN_BIN) ./cmd/admin

.PHONY: build-wss
build-wss: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(WSS_BIN) ./cmd/wss

.PHONY: build-cdn
build-cdn: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(CDN_BIN) ./cmd/cdn

.PHONY: build-task
build-task: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(TASK_BIN) ./cmd/task

.PHONY: build-task-data
build-task-data: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(TASK_DATA_BIN) ./cmd/task-data

.PHONY: build-task-jobs
build-task-jobs: $(BIN_DIR)
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(TASK_JOBS_BIN) ./cmd/task-jobs

.PHONY: clean
clean:
	@rm -rf $(BIN_DIR)
