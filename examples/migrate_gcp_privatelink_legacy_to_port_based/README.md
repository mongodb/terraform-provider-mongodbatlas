# Migration Example: GCP Private Link Legacy Architecture to Port-Based Architecture

This example demonstrates how to migrate from the legacy GCP Private Service Connect architecture to the new port-based architecture.

## Migration Phases

### v1: Initial State (Legacy Architecture)
Shows the original configuration using legacy architecture:
- `mongodbatlas_privatelink_endpoint` without `port_mapping_enabled` (defaults to `false`)
- 50 Google Compute Addresses
- 50 Google Compute Forwarding Rules
- `mongodbatlas_privatelink_endpoint_service` with `endpoints` list (50 endpoints)
- `gcp_project_id` required

### v2: Migration Phase (Both Architectures)
Demonstrates the migration approach:
- Creates a new `mongodbatlas_privatelink_endpoint` with `port_mapping_enabled = true`
- Adds new GCP resources (1 address, 1 forwarding rule) alongside legacy resources
- Creates a new `mongodbatlas_privatelink_endpoint_service` using `endpoint_service_id` and `private_endpoint_ip_address`
- Allows testing the new architecture before removing legacy resources

### v3: Final State (Port-Based Architecture Only)
Clean final configuration using only:
- `mongodbatlas_privatelink_endpoint` with `port_mapping_enabled = true`
- 1 Google Compute Address
- 1 Google Compute Forwarding Rule
- `mongodbatlas_privatelink_endpoint_service` with `endpoint_service_id` and `private_endpoint_ip_address` (no `endpoints` list)

## Usage

1. Start with v1 to understand the original setup
2. Apply v2 configuration to add the new port-based architecture resources
3. Update your application connection strings to use the new endpoint
4. Verify that the new architecture works correctly
5. Apply v3 configuration for the final clean state (removes legacy resources)

## Prerequisites

- MongoDB Atlas Terraform Provider with port-based architecture support
- Valid MongoDB Atlas project ID
- Google Cloud account with appropriate permissions
- GCP project with network and subnet configured

## Variables

```terraform
project_id          = "<ATLAS_PROJECT_ID>"
gcp_project_id      = "<GCP_PROJECT_ID>"
gcp_region          = "us-central1"
atlas_client_id     = "<ATLAS_CLIENT_ID>"     # Optional, can use env vars
atlas_client_secret = "<ATLAS_CLIENT_SECRET>" # Optional, can use env vars
cluster_name        = "<CLUSTER_NAME>"        # Optional: cluster whose connection string to output
```

**Note**: Network names, subnet names, IP addresses, and other resource names are hardcoded in the example configurations. You can modify these directly in the example files:
- Network name: `"my-network"`
- Subnet name: `"my-subnet"`
- IP CIDR range: `"10.0.0.0/16"`
- Address name (port-based): `"tf-test-port-based-endpoint"`
- Address IP (port-based): `"10.0.42.100"`

Alternatively, set environment variables for authentication:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

## Important Notes

- **Downtime**: Plan for brief downtime while changing the connection strings in your application after migrating to v2
- **State Management**: Backup your Terraform state before starting the migration
- **Testing**: Test the new architecture in v2 before proceeding to v3
