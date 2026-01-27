---
name: autogen-resource-from-mms
description: Autogenerate a Terraform resource from an OpenAPI spec in an MMS repository branch or PR. Use when the user wants to generate Terraform provider code from a pre-release API spec, mentions MMS branches, or references CLOUDP tickets.
---

# Autogenerate Terraform Resource from MMS OpenAPI Spec

This skill automates the process of generating a Terraform resource from an OpenAPI spec in the private `10gen/mms` repository.

## When to Use

- User wants to generate a Terraform resource from a feature branch OpenAPI spec
- User mentions an MMS PR number or branch name (e.g., "PR 153849", "CLOUDP-375419")
- User wants to test code generation with a pre-release API
- User asks to autogenerate a resource from a non-production OpenAPI spec

## Prerequisites

- Access to the private `10gen/mms` GitHub repository
- GitHub CLI (`gh`) authenticated with access to `10gen/mms`
- Node.js installed (for the OpenAPI transformer)
- Go installed (for the code generator)
- Resource configuration added to `tools/codegen/config.yml`

## Input Requirements

Ask the user for:

1. **MMS Branch or PR**: Either:
   - Branch name in `10gen/mms` (e.g., `CLOUDP-375419`, `master`)
   - PR number from `10gen/mms` (e.g., `153849`)

2. **Resource Name**: The resource name as defined in `config.yml`
   - Examples: `log_integration`, `service_account`, `alert_configuration_api`

## Instructions

### Step 1: Verify Resource Configuration

Check if the resource is configured in `tools/codegen/config.yml`. If not, ask the user to add the configuration first.

### Step 2: Resolve PR to Branch (if needed)

If the user provides a PR number, resolve it to a branch name:

```bash
gh pr view <PR_NUMBER> --repo 10gen/mms --json headRefName -q '.headRefName'
```

### Step 3: Fetch the OpenAPI Spec

Run the fetch script to download the spec from the MMS branch:

```bash
scripts/fetch-mms-spec.sh <BRANCH_NAME>
```

This downloads the spec to `tools/codegen/atlasapispec/raw-multi-version-api-spec.json`.

### Step 4: Flatten the OpenAPI Spec

Run the flatten script:

```bash
scripts/flatten-spec.sh
```

### Step 5: Generate the Resource

Run the code generator:

```bash
go run ./tools/codegen/main.go <resource_name>
```

### Step 6: Verify Generated Code

After generation:
1. Check for build errors: `go build ./...`
2. Review the generated schema in `internal/serviceapi/<packagename>/resource_schema.go`
3. Verify field types, descriptions, and sensitive markings

## Generated Files

The generator creates:
- `tools/codegen/models/<resource_name>.yaml` - Resource model
- `internal/serviceapi/<packagename>/resource.go` - Resource implementation
- `internal/serviceapi/<packagename>/resource_schema.go` - Resource schema
- `internal/serviceapi/<packagename>/data_source.go` - Data source (if configured)
- `internal/serviceapi/<packagename>/data_source_schema.go` - Data source schema
- `internal/serviceapi/<packagename>/plural_data_source.go` - Plural data source (if list configured)

## Post-Generation Steps

Remind the user to:
1. Register the resource with the Terraform provider (if new)
2. Add acceptance tests
3. Add documentation in `docs/resources/`
4. Update CHANGELOG.md

## Cleanup

After development, restore the production spec:

```bash
make generate-autogen-api-spec
```

## Troubleshooting

### Authentication Failed (401/403)
- Check GitHub CLI auth: `gh auth status`
- Verify access to `10gen/mms` repo

### Flattening Fails
- Validate JSON: `jq . tools/codegen/atlasapispec/raw-multi-version-api-spec.json > /dev/null`
- Check network connectivity for npx

### Code Generation Fails
- Verify resource is in `config.yml`
- Check OpenAPI spec has expected paths/schemas
- Review error for missing references
