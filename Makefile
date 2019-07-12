export GO111MODULE=off

GO ?= go

BUILD_PATH := $(shell pwd)/build
BUILD_BIN_PATH := ${BUILD_PATH}/bin
GOLANGCI_LINT := ${BUILD_BIN_PATH}/golangci-lint

define go-build
	$(shell cd `pwd` && $(GO) build -o $(BUILD_BIN_PATH)/$(shell basename $(1)) $(1))
	@echo > /dev/null
endef

${GOLANGCI_LINT}:
	$(call go-build,./vendor/github.com/golangci/golangci-lint/cmd/golangci-lint)

.PHONY: lint
lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} run

.PHONY: vendor
vendor:
	export GO111MODULE=on \
		$(GO) mod tidy && \
		$(GO) mod vendor && \
		$(GO) mod verify
