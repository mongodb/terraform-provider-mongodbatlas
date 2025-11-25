# Migration Example: Project Invitation to Cloud User Project Assignment

This example demonstrates how to migrate from the deprecated `mongodbatlas_project_invitation` resource to the new `mongodbatlas_cloud_user_project_assignment` resource.

## Migration Phases

### v1: Initial State (Deprecated Resource)
Shows the original configuration using deprecated `mongodbatlas_project_invitation` for pending invitations.

### v2: Migration Phase (Re-creation with Removed Block)
Demonstrates the migration approach:
- Adds new `mongodbatlas_cloud_user_project_assignment` resource
- Uses `removed` block to cleanly remove old resource from state
- Shows both removed block and manual state removal options

### v3: Final State (New Resource Only)
Clean final configuration using only `mongodbatlas_cloud_user_project_assignment`.

## Important Notes

- **Pending invites only**: This migration applies only to PENDING project invitations that still exist in your Terraform configuration
- **Re-creation approach**: The new resources and data sources cannot discover pending invites created by the deprecated resource, so we re-create them with the new resource
- **Accepted invites**: If users already accepted invitations, the provider removed them from state and you should remove them from configuration (no migration needed)

## Usage

1. Start with v1 to understand the original setup with pending invitations
2. Apply v2 configuration to re-create invites with new resource and remove old resource
3. Apply v3 configuration for the final clean state

## Prerequisites

- MongoDB Atlas Terraform Provider 2.0.0 or later
- Valid MongoDB Atlas project ID
- Pending project invitations in your configuration (not yet accepted)

## Variables

Set these variables for all versions:

```terraform
client_id     = "your-service-account-client-id"   # Optional, can use env vars
client_secret = "your-service-account-client-secret"  # Optional, can use env vars
project_id  = "your-project-id"
username    = "user@example.com"                # User with pending invitation
roles       = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_ONLY"]
```

Alternatively, set environment variables:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```
