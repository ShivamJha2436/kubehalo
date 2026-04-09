GO ?= go
GOMODCACHE ?= $(CURDIR)/.cache/gomod
GOCACHE ?= $(CURDIR)/.cache/go-build

export GOMODCACHE
export GOCACHE

.PHONY: fmt test run-controller run-api run-webhook

fmt:
	$(GO) fmt ./...

test:
	$(GO) test ./...

run-controller:
	$(GO) run ./cmd/controller

run-api:
	$(GO) run ./cmd/api

run-webhook:
	$(GO) run ./cmd/webhook
