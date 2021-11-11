TEST?=$$(go list ./... | grep -v /integrationtesting)
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=mongodbatlas

BINARY_NAME=terraform-provider-mongodbatlas
DESTINATION=./bin/$(BINARY_NAME)

WEBSITE_REPO=github.com/hashicorp/terraform-website

GOFLAGS=-mod=vendor
GOOPTS="-p 2"

GITTAG=$(shell git describe --always --tags)
VERSION=$(GITTAG:v%=%)
LINKER_FLAGS=-s -w -X 'github.com/mongodb/terraform-provider-mongodbatlas/version.ProviderVersion=${VERSION}'

GOLANGCI_VERSION=v1.41.1

export PATH := $(shell go env GOPATH)/bin:$(PATH)
export SHELL := env PATH=$(PATH) /bin/bash

default: build

.PHONY: build
build: fmtcheck
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

.PHONY: install
install: fmtcheck
	go install -ldflags="$(LINKER_FLAGS)"

.PHONY: test
test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

.PHONY: testacc
testacc: fmtcheck
	@$(eval VERSION=acc)
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -cover -ldflags="$(LINKER_FLAGS)"

.PHONY: fmt
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

.PHONY: fmtcheck
fmtcheck: # Currently required by tf-deploy compile
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: websitefmtcheck
websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

.PHONY: lint-fix
lint-fix:
	@echo "==> Checking source code against linters..."
	golangci-lint run --fix

.PHONY: lint
lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run

.PHONY: tools
tools:  ## Install dev tools
	@echo "==> Installing dependencies..."
	go install github.com/client9/misspell/cmd/misspell@latest
	go install github.com/terraform-linters/tflint@v0.31.0
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: check
check: test lint

.PHONY: test-compile
test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

.PHONY: website
website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: website-lint
website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

.PHONY: website-test
website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: terratest
terratest: fmtcheck
	@$(eval VERSION=acc)
	 go test $$(go list ./... | grep  /integrationtesting) -v -parallel 20 $(TESTARGS) -timeout 120m -cover -ldflags="$(LINKER_FLAGS)"

.PHONY: tflint
tflint: fmtcheck
	@scripts/tflint.sh
