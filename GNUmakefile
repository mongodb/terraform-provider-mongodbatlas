
ifdef ACCTEST_PACKAGES
		# remove newlines and blanks coming from GH Actions
    ACCTEST_PACKAGES := $(strip $(subst $(newline),, $(ACCTEST_PACKAGES)))
else
    ACCTEST_PACKAGES := "./..."
endif

ACCTEST_REGEX_RUN?=^TestAcc
ACCTEST_TIMEOUT?=300m
PARALLEL_GO_TEST?=50

BINARY_NAME=terraform-provider-mongodbatlas
DESTINATION=./bin/$(BINARY_NAME)

GOFLAGS=-mod=vendor
GITTAG=$(shell git describe --always --tags)
VERSION=$(GITTAG:v%=%)
LINKER_FLAGS=-s -w -X 'github.com/mongodb/terraform-provider-mongodbatlas/version.ProviderVersion=${VERSION}'

GOLANGCI_VERSION=v1.61.0 # Also update golangci-lint GH action in code-health.yml when updating this version

export PATH := $(shell go env GOPATH)/bin:$(PATH)
export SHELL := env PATH=$(PATH) /bin/bash

default: build

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: fmt fmtcheck
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

.PHONY: install
install: fmtcheck
	go install -ldflags="$(LINKER_FLAGS)"

.PHONY: test
test: fmtcheck
	go test ./... -timeout=30s -parallel=4 -race

.PHONY: testacc
testacc: fmtcheck
	@$(eval VERSION=acc)
	TF_ACC=1 go test $(ACCTEST_PACKAGES) -run '$(ACCTEST_REGEX_RUN)' -v -parallel $(PARALLEL_GO_TEST) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags="$(LINKER_FLAGS)"

.PHONY: testaccgov
testaccgov: fmtcheck
	@$(eval VERSION=acc)
	TF_ACC=1 go test ./... -run 'TestAccProjectRSGovProject_CreateWithProjectOwner' -v -parallel 1 "$(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags=$(LINKER_FLAGS) "

.PHONY: fmt
fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w .

.PHONY: fmtcheck
fmtcheck: ## Currently required by tf-deploy compile
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: lint-fix
lint-fix:
	@echo "==> Fixing linters errors..."
	fieldalignment -json -fix ./...
	golangci-lint run --fix

.PHONY: lint
lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run

.PHONY: tools
tools:  ## Install dev tools
	@echo "==> Installing dependencies..."
	go telemetry off # disable sending telemetry data, more info: https://go.dev/doc/telemetry
	go install github.com/icholy/gomajor@latest
	go install github.com/terraform-linters/tflint@v0.52.0
	go install github.com/rhysd/actionlint/cmd/actionlint@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	go install github.com/hashicorp/terraform-plugin-codegen-openapi/cmd/tfplugingen-openapi@latest
	go install github.com/hashicorp/terraform-plugin-codegen-framework/cmd/tfplugingen-framework@latest
	go install github.com/hashicorp/go-changelog/cmd/changelog-build@latest
	go install github.com/hashicorp/go-changelog/cmd/changelog-entry@latest
	go install golang.org/x/tools/cmd/goimports@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: docs
docs:
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"

.PHONY: tflint
tflint: fmtcheck
	tflint -f compact --recursive --minimum-failure-severity=warning

.PHONY: tf-validate
tf-validate: fmtcheck
	scripts/tf-validate.sh

.PHONY: link-git-hooks
link-git-hooks: ## Install git hooks
	@echo "==> Installing all git hooks..."
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

.PHONY: update-atlas-sdk
update-atlas-sdk: ## Update the atlas-sdk dependency
	./scripts/update-sdk.sh

# e.g. run: make scaffold resource_name=streamInstance type=resource
# - type argument can have the values: `resource`, `data-source`, `plural-data-source`.
# details on usage can be found in contributing/development-best-practices.md under "Scaffolding initial Code and File Structure"
.PHONY: scaffold
scaffold:
	@go run ./tools/scaffold/*.go $(resource_name) $(type)
	@echo "Reminder: configure the new $(type) in provider.go"

# e.g. run: make scaffold-schemas resource_name=streamInstance
# details on usage can be found in contributing/development-best-practices.md under "Scaffolding Schema and Model Definitions"
.PHONY: scaffold-schemas
scaffold-schemas:
	@scripts/schema-scaffold.sh $(resource_name)


.PHONY: generate-doc
# e.g. run: make generate-doc resource_name=search_deployment
# generate the resource documentation via tfplugindocs
generate-doc: 
	@scripts/generate-doc.sh ${resource_name}

# generate the resource documentation via tfplugindocs for all resources that have templates
.PHONY: generate-docs-all
generate-docs-all: 
	@scripts/generate-docs-all.sh

.PHONY: update-tf-compatibility-matrix
update-tf-compatibility-matrix: ## Update Terraform Compatibility Matrix documentation
	./scripts/update-tf-compatibility-matrix.sh

.PHONY: update-changelog-unreleased-section
update-changelog-unreleased-section:
	./scripts/update-changelog-unreleased-section.sh
  
.PHONY: generate-changelog-entry
generate-changelog-entry:
	./scripts/generate-changelog-entry.sh

.PHONY: check-changelog-entry-file
check-changelog-entry-file:
	go run ./tools/check-changelog-entry-file/*.go

.PHONY: jira-release-version
jira-release-version:
	go run ./tools/jira-release-version/*.go

.PHONY: enable-advancedclustertpf
enable-advancedclustertpf:
	make delete-lines filename="./internal/provider/provider_sdk2.go" delete="mongodbatlas_advanced_cluster"
	make add-lines filename=./internal/provider/provider.go find="project.Resource," add="advancedclustertpf.Resource,"
	make add-lines filename=./internal/provider/provider.go find="project.DataSource," add="advancedclustertpf.DataSource,"
	make add-lines filename=./internal/provider/provider.go find="project.PluralDataSource," add="advancedclustertpf.PluralDataSource,"

.PHONY: delete-lines ${filename} ${delete}
delete-lines:
	rm -f file.tmp
	grep -v "${delete}" "${filename}" > file.tmp
	mv file.tmp ${filename}
	goimports -w ${filename}

.PHONY: add-lines ${filename} ${find} ${add}
add-lines:
	rm -f file.tmp
	sed 's/${find}/${find}${add}/' "${filename}" > "file.tmp"
	mv file.tmp ${filename}
	goimports -w ${filename}
