GO ?= go

BUILD_PATH := $(shell pwd)/build
BUILD_BIN_PATH := ${BUILD_PATH}/bin
GOLANGCI_LINT := ${BUILD_BIN_PATH}/golangci-lint

${GOLANGCI_LINT}:
	export \
		VERSION=v1.54.2 \
		URL=https://raw.githubusercontent.com/golangci/golangci-lint \
		BINDIR=${BUILD_BIN_PATH} && \
	curl -sfL $$URL/$$VERSION/install.sh | sh -s $$VERSION

.PHONY: lint
lint: ${GOLANGCI_LINT}
	GL_DEBUG=gocritic ${GOLANGCI_LINT} run
