
GIT_COMMIT ?= $(shell git rev-parse --short HEAD)
VERSION ?= $(shell git describe --tags --abbrev=0)
BUILD_TIME ?= $(shell date --iso-8601=seconds)
CGO_ENABLED ?= 0
LDFLAGS := -extldflags "-static"
LDFLAGS += -X github.com/mozillazg/docker-credential-acr-helper/pkg/version.Version=$(VERSION)
LDFLAGS += -X github.com/mozillazg/docker-credential-acr-helper/pkg/version.GitCommit=$(GIT_COMMIT)
LDFLAGS += -X github.com/mozillazg/docker-credential-acr-helper/pkg/version.Timestamp=$(BUILD_TIME)

.PHONY: build
build:
	CGO_ENABLED=$(CGO_ENABLED) go build -ldflags "$(LDFLAGS)" -a -o docker-credential-acr-helper \
	cmd/docker-credential-acr-helper/main.go

.PHONY: test
test:
	go test -v ./...
