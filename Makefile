BINARY      := cloudlens
BUILD_DIR   := execs
INSTALL_DIR := $(HOME)/.local/bin
PACKAGE     := github.com/one2nc/$(BINARY)
VERSION     := v0.1.4
GIT_REV     ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
ifeq ($(shell uname), Darwin)
BUILD_TIME  ?= $(shell TZ=UTC date -j -f "%s" "$(shell date +%s)" +"%Y-%m-%dT%H:%M:%SZ")
else
BUILD_TIME  ?= $(shell date -u -d @$(shell date +%s) +"%Y-%m-%dT%H:%M:%SZ")
endif
LDFLAGS     := -ldflags "-w -s \
	-X '$(PACKAGE)/cmd.version=$(VERSION)' \
	-X '$(PACKAGE)/cmd.commit=$(GIT_REV)' \
	-X '$(PACKAGE)/cmd.date=$(BUILD_TIME)'"

.PHONY: all build install clean run test test-v cover setup setup-down

all: build ## Build the default artifact

build: ## Build cloudlens
	go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY) main.go

install: build ## Install cloudlens to ~/.local/bin (re-codesigns on macOS)
	cp $(BUILD_DIR)/$(BINARY) $(INSTALL_DIR)/$(BINARY)
ifeq ($(shell uname),Darwin)
	codesign --force --sign - $(INSTALL_DIR)/$(BINARY)
endif

clean: ## Remove build artifacts
	rm -rf $(BUILD_DIR)/$(BINARY)

run: build ## Build and run cloudlens
	./$(BUILD_DIR)/$(BINARY)

test: ## Run all tests
	go test ./...

test-v: ## Run tests with verbose output
	go test -v ./...

cover: ## Run tests with coverage report
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

setup: ## Start localstack via docker-compose
	docker-compose up -d

setup-down: ## Stop and remove localstack containers
	docker ps -a --format "{{.ID}} {{.Names}}" | grep cloudlens| awk '{print $$1}'| xargs docker stop | xargs docker rm -v

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help
