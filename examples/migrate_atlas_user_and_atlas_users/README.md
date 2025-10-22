# Migration Example: Atlas User to Cloud User Org Assignment

This example demonstrates how to migrate from the deprecated `mongodbatlas_atlas_user` and `mongodbatlas_atlas_users` data sources to their replacements.

## Migration Phases

### v1: Initial State (Deprecated Data Sources)
Shows the original configuration using deprecated data sources:
- `mongodbatlas_atlas_user` for single user reads
- `mongodbatlas_atlas_users` for user lists

### v2: Migration Phase (Both Old and New)
Demonstrates the migration approach:
- Adds new data sources alongside old ones
- Shows attribute mapping examples
- Validates new data sources work before removing old ones

### v3: Final State (New Data Sources Only)
Clean final configuration using only:
- `mongodbatlas_cloud_user_org_assignment` for single user reads
- `mongodbatlas_organization.users`, `mongodbatlas_project.users`, `mongodbatlas_team.users` for user lists

## Usage

1. Start with v1 to understand the original setup
2. Apply v2 configuration to add new data sources
3. Verify the new data sources return expected data
4. Update your references using the attribute mappings shown
5. Apply v3 configuration for the final clean state

## Prerequisites

- MongoDB Atlas Terraform Provider 2.0.0 or later
- Valid MongoDB Atlas organization, project, and team IDs
- Existing users in your organization

## Variables

Set these variables for all versions:

```terraform
client_id     = "<ATLAS_CLIENT_ID>"   # Optional, can use env vars
client_secret = "<ATLAS_CLIENT_SECRET>" # Optional, can use env vars
org_id     = "your-organization-id"
project_id = "your-project-id"  
team_id    = "your-team-id"
user_id    = "existing-user-id"
username   = "existing-user@example.com"
```

Alternatively, set environment variables:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```
