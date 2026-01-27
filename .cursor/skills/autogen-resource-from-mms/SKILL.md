---
name: autogen-resource-from-mms
description: Autogenerate a Terraform resource from an OpenAPI spec. Supports production, development, and MMS branch/PR environments. Use when generating Terraform provider code from any API spec source.
---

# Autogenerate Terraform Resource from OpenAPI Spec

This skill automates the process of generating a Terraform resource from an OpenAPI spec from different environments.

## When to Use

- User wants to generate a Terraform resource from the production or dev OpenAPI spec
- User mentions an MMS PR number or branch name (e.g., "PR 153849", "CLOUDP-375419")
- User wants to test code generation with a pre-release API
- User asks to autogenerate a resource

## Supported Environments

| Environment | Source | Description |
|-------------|--------|-------------|
| `prod` | mongodb/openapi (main) | Production OpenAPI spec |
| `dev` | mongodb/openapi (dev) | Development OpenAPI spec |
| `mms:<branch>` | 10gen/mms branch | Pre-release from MMS branch |
| `mms-pr:<number>` | 10gen/mms PR | Pre-release from MMS PR |

## Prerequisites

- GitHub CLI (`gh`) installed and authenticated
- For MMS environments: access to the private `10gen/mms` repository
- Node.js installed (for the OpenAPI transformer)
- Go installed (for the code generator)
- Resource configuration added to `tools/codegen/config.yml`

## Input Requirements

Ask the user for:

1. **Environment**: One of:
   - `prod` or `production` - Production spec
   - `dev` or `development` - Development spec
   - `mms:<branch>` - MMS branch (e.g., `mms:CLOUDP-375419`)
   - `mms-pr:<number>` - MMS PR (e.g., `mms-pr:153849`)

2. **Resource Name**: The resource name as defined in `config.yml`
   - Examples: `log_integration`, `service_account`, `alert_configuration_api`

## Instructions

**Important**: All commands must be run from the terraform-provider-mongodbatlas repository root directory.

### Quick Start: Run All Steps at Once

Use the all-in-one script:

```bash
.cursor/skills/autogen-resource-from-mms/scripts/autogen-resource.sh <environment> <resource_name>
```

Examples:
```bash
# From production spec
.cursor/skills/autogen-resource-from-mms/scripts/autogen-resource.sh prod log_integration

# From development spec
.cursor/skills/autogen-resource-from-mms/scripts/autogen-resource.sh dev log_integration

# From MMS branch
.cursor/skills/autogen-resource-from-mms/scripts/autogen-resource.sh mms:CLOUDP-375419 log_integration

# From MMS PR
.cursor/skills/autogen-resource-from-mms/scripts/autogen-resource.sh mms-pr:153849 log_integration
```

### Step-by-Step Instructions

#### Step 1: Verify Resource Configuration

Check if the resource is configured in `tools/codegen/config.yml`. If not, ask the user to add the configuration first.

#### Step 2: Fetch the OpenAPI Spec

```bash
.cursor/skills/autogen-resource-from-mms/scripts/fetch-spec.sh <environment>
```

Examples:
```bash
.cursor/skills/autogen-resource-from-mms/scripts/fetch-spec.sh prod
.cursor/skills/autogen-resource-from-mms/scripts/fetch-spec.sh dev
.cursor/skills/autogen-resource-from-mms/scripts/fetch-spec.sh mms:CLOUDP-375419
.cursor/skills/autogen-resource-from-mms/scripts/fetch-spec.sh mms-pr:153849
```

#### Step 3: Flatten the OpenAPI Spec

```bash
.cursor/skills/autogen-resource-from-mms/scripts/flatten-spec.sh
```

#### Step 4: Generate the Resource

```bash
go run ./tools/codegen/main.go <resource_name>
```

#### Step 5: Verify Generated Code

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
- For MMS: verify access to `10gen/mms` repo

### Flattening Fails
- Validate JSON/YAML spec syntax
- Check network connectivity for npx

### Code Generation Fails
- Verify resource is in `config.yml`
- Check OpenAPI spec has expected paths/schemas
- Review error for missing references
