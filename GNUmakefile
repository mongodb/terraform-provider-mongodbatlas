TEST?=$$(go list ./... | grep -v /integrationtesting)
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=mongodbatlas
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

build: fmtcheck
	go install -ldflags="$(LINKER_FLAGS)"

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	@$(eval VERSION=acc)
	TF_ACC=1 go test $(TEST) -v -parallel 20 $(TESTARGS) -timeout 300m -cover -ldflags="$(LINKER_FLAGS)"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./main.go
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

websitefmtcheck:
	@sh -c "'$(CURDIR)/scripts/websitefmtcheck.sh'"

lint-fix:
	@echo "==> Checking source code against linters..."
	golangci-lint run --fix

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run

tools:  ## Install dev tools
	@echo "==> Installing dependencies..."
	GO111MODULE=on go install github.com/client9/misspell/cmd/misspell
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

check: test lint

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

.PHONY: build test testacc fmt fmtcheck lint check tools test-compile website website-lint website-test

terratest: fmtcheck
	@$(eval VERSION=acc)
	 go test $$(go list ./... | grep  /integrationtesting) -v -parallel 20 $(TESTARGS) -timeout 120m -cover -ldflags="$(LINKER_FLAGS)"
