TEST?=./...
PKG_NAME=mongodbatlas

default: build

build: fmtcheck
	go install


test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@GOGC=30 golangci-lint run ./$(PKG_NAME)

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint


.PHONY: build test fmt fmtcheck lint tools
