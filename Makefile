MODULE      := gofscraper
BINARY      := gofscraper
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT      := $(shell git rev-parse --short HEAD 2>/dev/null || echo unknown)
BUILD_DATE  := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS     := -s -w \
	-X $(MODULE)/pkg/version.Version=$(VERSION) \
	-X $(MODULE)/pkg/version.Commit=$(COMMIT) \
	-X $(MODULE)/pkg/version.BuildDate=$(BUILD_DATE)

.PHONY: build run test lint clean generate

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BINARY) ./cmd/gofscraper

run: build
	./bin/$(BINARY)

test:
	go test ./... -race -count=1

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

generate:
	sqlc generate
