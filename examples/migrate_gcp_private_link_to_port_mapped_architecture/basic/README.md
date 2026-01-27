# Migration Example: GCP Private Link Legacy to Port-Mapped Architecture

This example demonstrates how to migrate from the legacy GCP Private Service Connect architecture to the port-mapped architecture.

## Migration Phases

### v1: Initial State (Legacy Architecture)
Shows the original configuration using legacy architecture:
- `mongodbatlas_privatelink_endpoint` without `port_mapping_enabled` (defaults to `false`).
- 50 Google Compute Addresses.
- 50 Google Compute Forwarding Rules.
- `mongodbatlas_privatelink_endpoint_service` with `endpoints` list (50 endpoints).
- `gcp_project_id` required.

### v2: Migration Phase (Both Architectures)
Demonstrates the migration approach:
- Creates a new `mongodbatlas_privatelink_endpoint` with `port_mapping_enabled = true`.
- Adds new GCP resources (1 address, 1 forwarding rule) alongside legacy resources.
- Creates a new `mongodbatlas_privatelink_endpoint_service` using `endpoint_service_id` and `private_endpoint_ip_address`.
- Allows testing the port-mapped architecture before removing legacy resources.

### v3: Final State (Port-Mapped Architecture Only)
Clean final configuration using only:
- `mongodbatlas_privatelink_endpoint` with `port_mapping_enabled = true`.
- 1 Google Compute Address.
- 1 Google Compute Forwarding Rule.
- `mongodbatlas_privatelink_endpoint_service` with `endpoint_service_id` and `private_endpoint_ip_address` (no `endpoints` list).

## Usage

1. Start with v1 to understand the original setup.
2. Apply v2 configuration to add the port-mapped architecture resources.
3. Update your application connection strings to use the port-mapped endpoint.
4. Verify that the port-mapped architecture works correctly.
5. Apply v3 configuration for the final clean state (removes legacy resources).

## Prerequisites

- MongoDB Atlas Terraform Provider with port-mapped architecture support.
- Valid MongoDB Atlas project ID.
- Google Cloud account with appropriate permissions.
- GCP project with network and subnet configured.

## Variables

```terraform
project_id          = "<ATLAS_PROJECT_ID>"
gcp_project_id      = "<GCP_PROJECT_ID>"
gcp_region          = "us-central1"
atlas_client_id     = "<ATLAS_CLIENT_ID>"     # Optional, can use env vars
atlas_client_secret = "<ATLAS_CLIENT_SECRET>" # Optional, can use env vars
cluster_name        = "<CLUSTER_NAME>"        # Optional: cluster whose connection string to output
endpoint_count      = 50                      # Optional: Number of endpoints for legacy architecture (defaults to 50, matches Atlas project's privateServiceConnectionsPerRegionGroup setting)
legacy_endpoint_service_id  = "legacy-endpoint-group" # Optional for v1 and v2: Endpoint service ID for legacy architecture (defaults to "legacy-endpoint-group")
new_endpoint_service_id     = "tf-test-port-mapped-endpoint" # Optional for v2 and v3: Endpoint service ID for port-mapped architecture (used as forwarding rule name and address name, defaults to "tf-test-port-mapped-endpoint")
```

**Note**: Network names, subnet names, IP addresses, and other resource names are hardcoded in the example configurations. You can modify these directly in the example files:
- Network name: `"my-network"`.
- Subnet name: `"my-subnet"`.
- IP CIDR range: `"10.0.0.0/16"`.
- Address IP (port-mapped): `"10.0.42.100"`.

Alternatively, set environment variables for authentication:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

## Important Notes

- **Downtime**: Plan for brief downtime during the migration process while changing the connection strings in your application.
- **State Management**: Backup your Terraform state before starting the migration.
- **Testing**: Test the port-mapped architecture in v2 before proceeding to v3.
