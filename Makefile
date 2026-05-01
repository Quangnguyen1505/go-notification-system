-include .env
export

GO ?= go

all: build test

deps:
	$(GO) mod tidy
	$(GO) mod download

build: deps
	$(GO) build -o build/proxy.exe ./cmd/proxy
	$(GO) build -o build/notification.exe ./cmd/notification

test:
	$(GO) test ./... -count=1

run:
	$(MAKE) -j2 run-proxy run-notification

run-proxy: deps
	$(GO) run ./cmd/proxy

run-notification: deps
	$(GO) run ./cmd/notification

.PHONY: all deps build test run run-proxy run-notification