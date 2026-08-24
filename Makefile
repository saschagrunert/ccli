GO ?= go

GOLANGCI_LINT_VERSION = v2.13.1

.PHONY: lint
lint:
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run
