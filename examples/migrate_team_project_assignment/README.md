# Migration Example: Team Project Attribute to Team Project Assignment

This example demonstrates how to migrate from the deprecated `mongodbatlas_project.teams` attribute to the new `mongodbatlas_team_project_assignment` resource.

## Migration Phases

### v1: Initial State (Deprecated Resource)
Shows the original configuration using deprecated `mongodbatlas_project.teams` attribute for team assignments.

### v2: Final State (New Resource Only)
Update the configuration to use `mongodbatlas_team_project_assignment` & migrate away from deprecated `mongodbatlas_project.teams`.

## Usage

1. Start with v1 to understand the original setup with team assignments
2. Apply v2 configuration to import existing assignments with new resource and no longer use deprecated attribute teams

## Prerequisites

- MongoDB Atlas Terraform Provider 2.0.0 or later
- Valid MongoDB Atlas and Team IDs

## Variables

Set these variables for all versions:

```terraform
public_key   = "your-mongodb-atlas-public-key"   # Optional, can use env vars
private_key  = "your-mongodb-atlas-private-key"  # Optional, can use env vars
team_id_1    = "your-team-id-1"                  # Team to assign
team_id_2    = "your-team-id-2"                  # Another team to assign
team_1_roles = ["GROUP_OWNER"]
team_2_roles = ["GROUP_READ_ONLY", "GROUP_DATA_ACCESS_READ_WRITE"]
```

Alternatively, set environment variables:
```bash
export MONGODB_ATLAS_PUBLIC_KEY="your-public-key"
export MONGODB_ATLAS_PRIVATE_KEY="your-private-key"
```