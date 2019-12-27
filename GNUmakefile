TEST?=./...
PKG_NAME=mongodbatlas
export GO111MODULE := on
export PATH := ./bin:$(PATH)

default: build

build:
	go install ./$(PKG_NAME)

test:
	go test $(TEST) -timeout=30s -parallel=4 -cover

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run $(TEST) -E gofmt -E golint -E misspell

check: test lint

tools:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s v1.21.0

.PHONY: build test fmt lint check tools
