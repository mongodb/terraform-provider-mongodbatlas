# Module Maintainer - GCP Private Link Legacy to Port-Mapped Architecture

If you own and maintain modules to manage your Terraform resources, the purpose of this example is to demonstrate how a Terraform module definition can migrate from the legacy GCP Private Service Connect architecture to the port-mapped architecture while minimizing impact to its clients. The [module user example](../module_user/README.md) explains the same process from the module user point of view.

The example contains three module versions which represent the three steps of the migration:

Step | Purpose | Architecture
--- | --- | ---
[Step 1](./v1) | Baseline | Legacy architecture (uses `endpoints` list)
[Step 2](./v2) | Add port-mapped support | Creates both legacy and port-mapped architectures simultaneously
[Step 3](./v3) | Port-mapped only | Port-mapped architecture only (removes legacy support)

The rest of this document summarizes the different implementations:

- [Step 1: Module `v1` Implementation Summary](#step-1-module-v1-implementation-summary)
  - [`variables.tf`](#variablestf)
  - [`main.tf`](#maintf)
  - [`outputs.tf`](#outputstf)
- [Step 2: Module `v2` Implementation Changes and Highlights](#step-2-module-v2-implementation-changes-and-highlights)
  - [`variables.tf`](#variablestf-1)
  - [`main.tf`](#maintf-1)
  - [`outputs.tf`](#outputstf-1)
- [Step 3: Module `v3` Implementation Changes and Highlights](#step-3-module-v3-implementation-changes-and-highlights)
  - [`variables.tf`](#variablestf-2)
  - [`main.tf`](#maintf-2)
  - [`outputs.tf`](#outputstf-2)

## Step 1: Module `v1` Implementation Summary

This module creates GCP private link resources using the legacy architecture.

### [`variables.tf`](v1/variables.tf)

An abstraction for the `mongodbatlas_privatelink_endpoint` and `mongodbatlas_privatelink_endpoint_service` resources:
- Exposes variables for project ID, GCP project ID, region, and network configuration
- Legacy architecture requires an `endpoints` list with multiple endpoints (configurable via `legacy_endpoint_count` variable)

### [`main.tf`](v1/main.tf)

It uses the legacy architecture:
- Creates `mongodbatlas_privatelink_endpoint` without `port_mapping_enabled` (defaults to `false`)
- Creates multiple Google Compute Addresses (number matches `legacy_endpoint_count` variable)
- Creates multiple Google Compute Forwarding Rules (number matches `legacy_endpoint_count` variable)
- Creates `mongodbatlas_privatelink_endpoint_service` with a `dynamic "endpoints"` block

### [`outputs.tf`](v1/outputs.tf)

- Exposes the endpoint service ID and connection string
- Outputs `connection_string_legacy` - the connection string of the private endpoint with legacy architecture
- Outputs the full `mongodbatlas_privatelink_endpoint_service` resource for reference

## Step 2: Module `v2` Implementation Changes and Highlights

This is the new version of the module where support for port-mapped architecture is added. The implementation creates both legacy and port-mapped architectures simultaneously, allowing users to test the port-mapped architecture while keeping legacy resources intact during migration.

### [`variables.tf`](v2/variables.tf)

- Adds port-mapped architecture variables (`new_endpoint_service_id`, `port_mapped_endpoint_ip`)
- Keeps all existing variables for backward compatibility
- No `port_mapping_enabled` variable - both architectures are always created

### [`main.tf`](v2/main.tf)

- Creates both legacy and port-mapped architectures simultaneously:
  - Legacy architecture: Creates `mongodbatlas_privatelink_endpoint.legacy` (without `port_mapping_enabled`), multiple Google Compute Addresses (number matches `legacy_endpoint_count` variable), multiple Forwarding Rules (number matches `legacy_endpoint_count` variable), and `mongodbatlas_privatelink_endpoint_service.legacy` with `endpoints` list
  - Port-mapped architecture: Creates `mongodbatlas_privatelink_endpoint.new` (with `port_mapping_enabled = true`), 1 Google Compute Address, 1 Forwarding Rule, and `mongodbatlas_privatelink_endpoint_service.new` with `endpoint_service_id` and `private_endpoint_ip_address`
- Both architectures coexist in the same configuration, allowing parallel testing during migration

### [`outputs.tf`](v2/outputs.tf)

- Outputs connection strings for both legacy and port-mapped architectures
- Outputs separate resources for legacy (`mongodbatlas_privatelink_endpoint_legacy`, `mongodbatlas_privatelink_endpoint_service_legacy`) and new (`mongodbatlas_privatelink_endpoint_new`, `mongodbatlas_privatelink_endpoint_service_new`)
- Maintains backward compatibility with `v1` outputs

## Step 3: Module `v3` Implementation Changes and Highlights

This module removes support for the legacy architecture and only supports the port-mapped architecture. A major version bump would typically accompany this module version since we remove input variables and change the module behavior.

### [`variables.tf`](v3/variables.tf)

- Removes variables related to legacy architecture (e.g., `legacy_endpoint_count`, `endpoint_base_name`, `legacy_endpoint_service_id`, `endpoint_base_ip`)
- Keeps only port-mapped architecture variables (`new_endpoint_service_id`, `port_mapped_endpoint_ip`)
- Keeps common variables (project_id, gcp_project_id, gcp_region, network_name, subnet_name, subnet_ip_cidr_range, cluster_name)

### [`main.tf`](v3/main.tf)

- Creates only port-mapped architecture resources:
  - Creates `mongodbatlas_privatelink_endpoint.new` with `port_mapping_enabled = true`
  - Creates only 1 Google Compute Address
  - Creates only 1 Google Compute Forwarding Rule
  - Creates `mongodbatlas_privatelink_endpoint_service.new` with `endpoint_service_id` and `private_endpoint_ip_address` (no `endpoints` list)
- All legacy resources are removed

### [`outputs.tf`](v3/outputs.tf)

- Simplifies outputs to only expose port-mapped architecture connection strings
- Removes legacy-specific outputs
- Outputs `connection_string_new` - the connection string of the private endpoint with port-mapped architecture
