GO ?= go
GOMODCACHE ?= $(CURDIR)/.cache/gomod
GOCACHE ?= $(CURDIR)/.cache/go-build

export GOMODCACHE
export GOCACHE

.PHONY: fmt lint test run-controller run-api run-webhook

fmt:
	$(GO) fmt ./...

lint:
	@files="$$(git ls-files '*.go')"; \
	unformatted="$$(gofmt -l $$files)"; \
	test -z "$$unformatted" || (echo "The following files need gofmt:" && echo "$$unformatted" && exit 1)
	$(GO) vet ./...

test:
	$(GO) test ./...

run-controller:
	$(GO) run ./cmd/controller

run-api:
	$(GO) run ./cmd/api

run-webhook:
	$(GO) run ./cmd/webhook
