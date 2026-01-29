# Migration Example: GCP Private Service Connect to Port-Mapped Architecture

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
# Provider variables
atlas_client_id     = "<ATLAS_CLIENT_ID>"     # Optional, can use env vars
atlas_client_secret = "<ATLAS_CLIENT_SECRET>" # Optional, can use env vars

# Common variables between architectures
project_id          = "<ATLAS_PROJECT_ID>"
gcp_project_id      = "<GCP_PROJECT_ID>"
gcp_region          = "<GCP_REGION>"
cluster_name        = "<CLUSTER_NAME>"
network_name        = "<NETWORK_NAME>"
subnet_name         = "<SUBNET_NAME>"
subnet_ip_cidr_range = "<SUBNET_IP_CIDR_RANGE>"

# Legacy architecture variables
legacy_endpoint_count      = <LEGACY_ENDPOINT_COUNT>
legacy_endpoint_service_id  = "<LEGACY_ENDPOINT_SERVICE_ID>"
legacy_address_name_prefix = "<LEGACY_ADDRESS_NAME_PREFIX>"
legacy_address_base_ip     = "<LEGACY_ADDRESS_BASE_IP>"

# Port-mapped architecture variables
port_mapped_endpoint_service_id     = "<PORT_MAPPED_ENDPOINT_SERVICE_ID>"
port_mapped_address_ip = "<PORT_MAPPED_ADDRESS_IP>"
```

**Note**: For this migration guide, some sample values have been set as defaults for network names, subnet names, IP addresses, and other resource names. You can override any of these in your `terraform.tfvars` file or via command line flags.

Alternatively, set environment variables for authentication:
```bash
export MONGODB_ATLAS_CLIENT_ID="<ATLAS_CLIENT_ID>"
export MONGODB_ATLAS_CLIENT_SECRET="<ATLAS_CLIENT_SECRET>"
```

## Important Notes

- **Downtime**: Plan for brief downtime during the migration process while changing the connection strings in your application.
- **State Management**: Backup your Terraform state before starting the migration.
- **Testing**: Test the port-mapped architecture in v2 before proceeding to v3.
