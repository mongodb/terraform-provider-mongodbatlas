ifdef ACCTEST_PACKAGES
		# remove newlines and blanks coming from GH Actions
    ACCTEST_PACKAGES := $(strip $(subst $(newline),, $(ACCTEST_PACKAGES)))
else
    ACCTEST_PACKAGES := "./..."
endif

ACCTEST_TIMEOUT?=300m
PARALLEL_GO_TEST?=50

BINARY_NAME=terraform-provider-mongodbatlas
DESTINATION=./bin/$(BINARY_NAME)

GITTAG=$(shell git describe --always --tags)
VERSION=$(GITTAG:v%=%)
LINKER_FLAGS=-s -w -X 'github.com/mongodb/terraform-provider-mongodbatlas/version.ProviderVersion=${VERSION}'

GOLANGCI_VERSION=v2.10.1

export PATH := $(shell go env GOPATH)/bin:$(PATH)
export SHELL := env PATH=$(PATH) /bin/bash

default: fix

.PHONY: help
help:
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' | sort

.PHONY: fix
fix: ## Fix, format, and build Go code (default target)
	gofmt -s -w .
	golangci-lint run --fix
	go mod tidy
	go fix ./...
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

.PHONY: verify
verify: ## Verify Go code without modifying files. Usage: make verify [files="file1.go file2.go"]
	@bad_fmt=$$(gofmt -l -s $(or $(files),.)); \
	if [ -n "$$bad_fmt" ]; then echo "ERROR: gofmt issues:"; echo "$$bad_fmt"; exit 1; fi
ifdef files
	golangci-lint run $(addsuffix ...,$(sort $(dir $(files))))
	go fix -diff $(addsuffix ...,$(sort $(dir $(files))))
else
	golangci-lint run
	go mod tidy -diff
	go fix -diff ./...
endif

.PHONY: build
build: ## Compile the provider binary
	go build -ldflags "$(LINKER_FLAGS)" -o $(DESTINATION)

.PHONY: clean-atlas-org
clean-atlas-org: ## Run a test to clean all projects and pending resources in an Atlas org, supports export DRY_RUN=false (default=true)
	@$(eval export MONGODB_ATLAS_CLEAN_ORG?=true)
	@$(eval export DRY_RUN?=true)
	go test -count=1 'github.com/mongodb/terraform-provider-mongodbatlas/internal/testutil/clean' -timeout 3600s -parallel=250 -run 'TestCleanProjectAndClusters' -v -ldflags="$(LINKER_FLAGS)"

.PHONY: test
test: ## Run unit tests
	@$(eval export HTTP_MOCKER_REPLAY?=true)
	@$(eval export MONGODB_ATLAS_ORG_ID?=111111111111111111111111)
	@$(eval export MONGODB_ATLAS_PROJECT_ID?=111111111111111111111111)
	@$(eval export MONGODB_ATLAS_CLUSTER_NAME?=mocked-cluster)
	@$(eval export MONGODB_ATLAS_PUBLIC_KEY=)
	@$(eval export MONGODB_ATLAS_PRIVATE_KEY=)
	@$(eval export MONGODB_ATLAS_CLIENT_ID=)
	@$(eval export MONGODB_ATLAS_CLIENT_SECRET=)
	@$(eval export MONGODB_ATLAS_ACCESS_TOKEN=)
	go test ./... -timeout=120s -parallel=$(PARALLEL_GO_TEST) -race

.PHONY: testmact
testmact: ## Run MacT tests (mocked acc tests)
	@$(eval export ACCTEST_REGEX_RUN?=^TestAccMockable)
	@$(eval export HTTP_MOCKER_REPLAY?=true)
	@$(eval export HTTP_MOCKER_CAPTURE?=false)
	@$(eval export MONGODB_ATLAS_ORG_ID?=111111111111111111111111)
	@$(eval export MONGODB_ATLAS_PROJECT_ID?=111111111111111111111111)
	@$(eval export MONGODB_ATLAS_CLUSTER_NAME?=mocked-cluster)
	@if [ "$(ACCTEST_PACKAGES)" = "./..." ]; then \
		echo "Error: ACCTEST_PACKAGES must be explicitly set for testmact target, './...' is not allowed"; \
		exit 1; \
	fi
	TF_ACC=1 go test $(ACCTEST_PACKAGES) -run '$(ACCTEST_REGEX_RUN)' -v -parallel $(PARALLEL_GO_TEST) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags="$(LINKER_FLAGS)"

.PHONY: testmact-capture
testmact-capture: ## Capture HTTP traffic for MacT tests
	@$(eval export ACCTEST_REGEX_RUN?=^TestAccMockable)
	@$(eval export HTTP_MOCKER_REPLAY?=false)
	@$(eval export HTTP_MOCKER_CAPTURE?=true)
	@if [ "$(ACCTEST_PACKAGES)" = "./..." ]; then \
		echo "Error: ACCTEST_PACKAGES must be explicitly set for testmact-capture target, './...' is not allowed"; \
		exit 1; \
	fi
	TF_ACC=1 go test $(ACCTEST_PACKAGES) -run '$(ACCTEST_REGEX_RUN)' -v -parallel $(PARALLEL_GO_TEST) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags="$(LINKER_FLAGS)"

.PHONY: testacc
testacc: ## Run acc & mig tests (acceptance & migration tests)
	@$(eval export ACCTEST_REGEX_RUN?=^TestAcc)
	TF_ACC=1 go test $(ACCTEST_PACKAGES) -run '$(ACCTEST_REGEX_RUN)' -v -parallel $(PARALLEL_GO_TEST) $(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags="$(LINKER_FLAGS)"

.PHONY: testaccgov
testaccgov: ## Run Government cloud-provider acc & mig tests
	TF_ACC=1 go test ./... -run 'TestAccProjectRSGovProject_CreateWithProjectOwner' -v -parallel 1 "$(TESTARGS) -timeout $(ACCTEST_TIMEOUT) -ldflags=$(LINKER_FLAGS) "


.PHONY: tools
tools:  ## Install the dev tools (dependencies)
	@echo "==> Installing dependencies..."
	go telemetry off # disable sending telemetry data, more info: https://go.dev/doc/telemetry
	go install github.com/icholy/gomajor@latest
	go install github.com/terraform-linters/tflint@v0.61.0
	go install github.com/rhysd/actionlint/cmd/actionlint@latest
	go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
	go install github.com/hashicorp/go-changelog/cmd/changelog-build@latest
	go install github.com/hashicorp/go-changelog/cmd/changelog-entry@latest
	go install golang.org/x/tools/cmd/goimports@latest
	go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin $(GOLANGCI_VERSION)

.PHONY: docs
docs: ## Give URL to test Terraform documentation
	@echo "Use this site to preview markdown rendering: https://registry.terraform.io/tools/doc-preview"


.PHONY: tflint
tflint: ## Linter for Terraform files in examples/ dir (avoid `internal/**/testdata/main*.tf`), disable terraform_required_providers rule as we intentionally omit the provider version
	tflint --chdir=examples/ -f compact --recursive --minimum-failure-severity=warning --disable-rule=terraform_required_providers

.PHONY: tf-validate
tf-validate: ## Validate Terraform files
	scripts/tf-validate.sh

.PHONY: resign-commits
resign-commits: ## Rebase commits ahead of master and re-sign them with GPG. Usage: make resign-commits [base=master]
	./scripts/resign-commits.sh $(or $(base),master)

.PHONY: link-git-hooks
link-git-hooks: ## Install Git hooks
	@echo "==> Installing all git hooks..."
	find .git/hooks -type l -exec rm {} \;
	find .githooks -type f -exec ln -sf ../../{} .git/hooks/ \;

.PHONY: update-atlas-sdk
update-atlas-sdk: ## Update the Atlas SDK dependency
	./scripts/update-sdk.sh

# e.g. run: make scaffold resource_name=streamInstance type=resource
# type - valid values: `resource`, `data-source`, `plural-data-source`.
# details on usage can be found in contributing/development-best-practices.md under "Scaffolding initial Code and File Structure"
.PHONY: scaffold
scaffold: ## Create scaffolding for a new resource
	@go run ./tools/scaffold/*.go $(resource_name) $(type)
	@echo "Reminder: configure the new $(type) in provider.go"

# Generate flattened API spec used by codegen
# api_spec_url (optional) - URL to the OpenAPI spec (default: https://raw.githubusercontent.com/mongodb/openapi/main/openapi/v2.yaml).
.PHONY: autogen-update-api-spec
autogen-update-api-spec:
	@scripts/generate-autogen-api-spec.sh $(api_spec_url)

# Generate resources using API spec present in tools/codegen/atlasapispec/multi-version-api-spec.flattened.yml
# resource_name (optional) - If not provided all configured resource models will be generated.
# resource_tier (optional) - Valid values: `prod`, `internal` (default: all).
# step (optional) - Valid values: `model-gen`, `code-gen` (default: both).
# e.g. make autogen-generate-resources resource_name=search_deployment_api
.PHONY: autogen-generate-resources
autogen-generate-resources:
	@go run ./tools/codegen/main.go $(if $(resource_name),--resource-name $(resource_name),) $(if $(resource_tier),--resource-tier $(resource_tier),) $(if $(step),--step $(step),)

# Complete generation pipeline: Fetch latest API Spec -> update resource models -> generate resource code
# api_spec_url (optional) - URL to the OpenAPI spec (default: https://raw.githubusercontent.com/mongodb/openapi/main/openapi/v2.yaml).
# resource_name (optional) - If not provided all configured resources code will be generated
# resource_tier (optional) - Valid values: `prod`, `internal` (default: all)
# step (optional) - Valid values: `model-gen`, `code-gen` (default: both).
# e.g. make autogen-pipeline resource_tier=prod
.PHONY: autogen-pipeline
autogen-pipeline: autogen-update-api-spec autogen-generate-resources

.PHONY: generate-doc
# e.g. run: make generate-doc resource_name=search_deployment
# generate the resource documentation via tfplugindocs
generate-doc: ## Auto-generate the documentation for a resource
	@scripts/generate-doc.sh ${resource_name}

# generate the resource documentation via tfplugindocs for all resources that have templates
# note: if you also want to generate documentation for autogen resources, first you need to register resources in the provider using `enable-autogen` and
# generate the code using `autogen-pipeline`.
# autogenerated resource docs make use of generic resource template which is limited in its content, no examples (either tf config or import commands) are present
.PHONY: generate-docs-all
generate-docs-all: ## Auto-generate the documentation for all resources
	@scripts/generate-docs-all.sh

.PHONY: update-tf-compatibility-matrix
update-tf-compatibility-matrix: ## Update Terraform Compatibility Matrix documentation
	./scripts/update-tf-compatibility-matrix.sh

.PHONY: update-tf-version-in-repository
update-tf-version-in-repository: ## Update Terraform versions
	./scripts/update-tf-version-in-repository.sh

.PHONY: update-changelog-unreleased-section
update-changelog-unreleased-section: ## Update changelog unreleased section
	./scripts/update-changelog-unreleased-section.sh
  
.PHONY: generate-changelog-entry
generate-changelog-entry: ## Generate a changelog entry in a PR
	./scripts/generate-changelog-entry.sh

.PHONY: check-changelog-entry-file
check-changelog-entry-file: ## Check a changelog entry file in a PR
	go run ./tools/check-changelog-entry-file/*.go

.PHONY: jira-release-version
jira-release-version: ## Update Jira version in a release
	go run ./tools/jira-release-version/*.go

.PHONY: access-token-create
access-token-create: ## Create a new OAuth2 access token from Service Account credentials
	@go run ./tools/access-token/*.go create

.PHONY: access-token-revoke
access-token-revoke: ## Revoke an OAuth2 access token. Usage: make access-token-revoke token=<token>
	@go run ./tools/access-token/*.go revoke $(token)

.PHONY: enable-autogen
enable-autogen: ## Enable use of autogen resources and datasources in the provider
	@go run tools/enable-autogen/main.go

.PHONY: delete-lines ${filename} ${delete}
delete-lines:
	rm -f file.tmp
	grep -v "${delete}" "${filename}" > file.tmp
	mv file.tmp ${filename}
	goimports -w ${filename}

.PHONY: add-lines ${filename} ${find} ${add}
add-lines:
	rm -f file.tmp
	sed 's/${find}/${add}${find}/' "${filename}" > "file.tmp"
	mv file.tmp ${filename}

.PHONY: add-lines-if-missing ${filename} ${resource}
add-lines-if-missing:
	@if ! grep -q "${resource}.Resource," "${filename}" 2>/dev/null; then \
		make add-lines filename=${filename} find="project.Resource," add="${resource}.Resource,\n"; \
	fi

.PHONY: add-datasource-if-exists ${filename} ${resource}
add-datasource-if-exists:
	@if [ -f "internal/serviceapi/${resource}/data_source.go" ] && ! grep -q "${resource}.DataSource," "${filename}" 2>/dev/null; then \
		make add-lines filename=${filename} find="project.DataSource," add="${resource}.DataSource,\n"; \
	fi

.PHONY: change-lines ${filename} ${find} ${new}
change-lines:
	rm -f file.tmp
	sed 's/${find}/${new}/' "${filename}" > "file.tmp"
	mv file.tmp ${filename}

.PHONY: gen-purls
gen-purls: # Generate purls on linux os
	./scripts/compliance/generate-purls.sh

.PHONY: generate-sbom
generate-sbom: ## Generate SBOM
	./scripts/compliance/generate-sbom.sh

.PHONY: upload-sbom
upload-sbom: ## Upload SBOM
	./scripts/compliance/upload-sbom.sh

.PHONY: augment-sbom
augment-sbom: ## Augment SBOM
	./scripts/compliance/augment-sbom.sh
