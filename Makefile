GOCMD?=go

APP_VERSION?=$(shell git describe --dirty --tags --match "v[0-9]*" )

PKG_SRC=$(shell find . -type f -name '*.go')
GOOS=$(shell $(GOCMD) env GOOS)
GOARCH=$(shell $(GOCMD) env GOARCH)

BIN_NAME=issuectl-$(GOOS)-$(GOARCH) # -$(APP_VERSION)

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

BUILD_LDFLAGS=-ldflags "\
-X main.BuildVersion=$(APP_VERSION)"

##@ General

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development
.PHONY: test
test: ## Run tests
	go test ./...


.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: lint
lint:
	golangci-lint run -c .golangci.yml

.PHONY: check
check: ## Run all static checks
check: lint vet fmt

##@ Build

issuectl: ## Build issuectl binary
issuectl: $(PKG_SRC) cmd
	go build $(BUILD_LDFLAGS) -o $(BIN_NAME) cmd/main.go
