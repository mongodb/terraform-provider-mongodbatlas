# Migration Example: Organization-Level API Keys to Service Accounts

This example demonstrates how to migrate from organization-level Programmatic API Key (PAK) resources to Service Account (SA) resources.

## Migration Phases

### v1: Initial State (PAK Resources)
Shows the original configuration using PAK resources:
- `mongodbatlas_api_key` for organization-level API key
- `mongodbatlas_api_key_project_assignment` for project assignment
- `mongodbatlas_access_list_api_key` for IP access list entry

### v2: Migration Phase (Both PAK and SA Resources)
Demonstrates the migration approach:
- Adds Service Account resources alongside existing PAK resources
- Includes output to capture the Service Account secret
- Allows testing Service Accounts before removing PAKs

### v3: Final State (SA Resources Only)
Clean final configuration using only:
- `mongodbatlas_service_account` for organization-level service account
- `mongodbatlas_service_account_project_assignment` for project assignment
- `mongodbatlas_service_account_access_list_entry` for IP access list entry

## Usage

1. Start with v1 to understand the original setup
2. Apply v2 configuration to add Service Account resources
3. Retrieve and securely store the Service Account secret from the output
4. Verify that both PAK and SA authentication methods work correctly
5. Apply v3 configuration for the final clean state

## Prerequisites

- MongoDB Atlas Terraform Provider with Service Account support
- Valid MongoDB Atlas organization and project IDs
- Appropriate permissions to manage API keys and Service Accounts

## Variables

Set these variables for all versions:

```terraform
atlas_client_id     = "<ATLAS_CLIENT_ID>"     # Optional, can use env vars
atlas_client_secret = "<ATLAS_CLIENT_SECRET>" # Optional, can use env vars
org_id              = "your-organization-id"
project_id          = "your-project-id"
org_roles           = ["ORG_MEMBER"]
project_roles       = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
cidr_block          = "192.168.1.100/32"
```

For v2 and v3, also set:
```terraform
service_account_name        = "example-service-account" # Optional
secret_expires_after_hours = 2160                      # Optional, 90 days default
```

Alternatively, set environment variables:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```
