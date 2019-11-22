GO ?= go

BUILD_PATH := $(shell pwd)/build
BUILD_BIN_PATH := ${BUILD_PATH}/bin
GOLANGCI_LINT := ${BUILD_BIN_PATH}/golangci-lint

${GOLANGCI_LINT}:
	export \
		VERSION=v1.21.0 \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=${BUILD_BIN_PATH} && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION

.PHONY: lint
lint: ${GOLANGCI_LINT}
	${GOLANGCI_LINT} run
