---
page_title: "Migration Guide: GCP Private Link Legacy to Port-Mapped Architecture"
---

# Migration Guide: GCP Private Link Legacy to Port-Mapped Architecture

## Overview

This guide explains how to migrate from the legacy GCP Private Service Connect architecture to the port-mapped architecture for MongoDB Atlas private link endpoints.

The steps in this guide are for migrating Terraform-managed GCP private link endpoint resources, `mongodbatlas_privatelink_endpoint`, and `mongodbatlas_privatelink_endpoint_service`. The legacy architecture requires dedicated resources for each Atlas node. The port-mapped architecture uses a single set of resources to support up to 150 nodes through port mapping, enabling direct targeting of specific nodes using only one customer IP address.

**Note:** Migration to the port-mapped architecture is recommended but **not required**. If you are currently using the legacy architecture, you may continue to do so. This guide is for users who wish to adopt the port-mapped architecture for simplified management and reduced resource overhead.

## Why Migrate to Port-Mapped Architecture?

The legacy architecture has two main limitations:

1. **IP Exhaustion**: Atlas defaults to 50 private service connections per region group (50 forwarding rules and 50 IP addresses), which can lead to IP address exhaustion in your GCP project.

2. **Static Configuration**: Changing the number of private service connections per region group requires a full private service connect redeployment, causing friction when changing cluster configurations.

The port-mapped architecture addresses these limitations by using a single set of resources to support up to 150 nodes, requiring only 1 Google Compute Address and 1 Google Compute Forwarding Rule.

## Architecture Comparison

The following table shows the key differences between the legacy and port-mapped architectures:

| Aspect | Legacy Architecture | Port-Mapped Architecture |
|--------|---------------------|------------------------|
| `mongodbatlas_privatelink_endpoint.port_mapping_enabled` | Not set (defaults to `false`) | Must be set to `true` |
| `google_compute_address` count | One per Atlas node | 1 address (total, supports up to 150 nodes) |
| `google_compute_forwarding_rule` count | One per Atlas node | 1 forwarding rule (total, supports up to 150 nodes) |
| `mongodbatlas_privatelink_endpoint_service.endpoint_service_id` | Required (can be any identifier string) | Required (is the forwarding rule name) |
| `mongodbatlas_privatelink_endpoint_service.private_endpoint_ip_address` | Not used | Required (the IP address of the forwarding rule) |
| `mongodbatlas_privatelink_endpoint_service.endpoints` | Required (one endpoint per Atlas node) | Not used |
| `mongodbatlas_privatelink_endpoint_service.gcp_project_id` | Required | Required |
| `mongodbatlas_privatelink_endpoint_service.endpoint_group_names` | A list of endpoint group names associated with the private endpoint service | A list of private endpoint names associated with the private endpoint service |
| `mongodbatlas_privatelink_endpoint_service.service_attachment_names` | A list of service attachments connected to the private endpoint service (one per Atlas node) | A list of one service attachment connected to the private endpoint service |
| Connection String Format | Uses `pl-0` identifier (e.g., `cluster0-pl-0.a0b1c2.domain.com`) | Uses `psc-0` identifier (e.g., `cluster0-psc-0.a0b1c2.domain.com`). **Exception:** Cross-cloud clusters that span a region with a port-mapped endpoint service continue using `pl-0`. |

## Before You Begin

- **Backup your Terraform state file** before making any changes.
- **Test the process in a non-production environment** if possible.
- Ensure you have the necessary GCP permissions to create and delete Compute Addresses and Forwarding Rules.

### Important Considerations

#### Cannot Modify Existing mongodbatlas_privatelink_endpoint

**You cannot modify an existing `mongodbatlas_privatelink_endpoint` to enable port mapping.** The `port_mapping_enabled` attribute must be set when the `mongodbatlas_privatelink_endpoint` is first created. If you need to migrate, you must:

1. Create a new `mongodbatlas_privatelink_endpoint` with `port_mapping_enabled = true`.
2. Create new GCP resources (1 address, 1 forwarding rule).
3. Create a new `mongodbatlas_privatelink_endpoint_service` linking to the new `mongodbatlas_privatelink_endpoint`.
4. Update your application connection strings.
5. Delete unused resources.

#### Downtime

**Downtime occurs during the migration process when updating application connection strings**, not during Terraform operations. You can maintain both your legacy and port-mapped architectures in the same region during the transition. This ensures a stable migration path before you tear down the original resource.

After creating the port-mapped resources in Step 2, you will need to test and update your application connection strings to use the port-mapped private endpoint. You can retrieve the updated connection string from your cluster's private endpoint configuration.

---

## Migration Steps

For complete migration examples showing the step-by-step transition from legacy to port-mapped architecture, see the [GCP Private Link migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture).

- **Direct Resource Management**: If you are managing the resources directly (not using modules), see the [basic migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/basic).
- **Module Maintainers**: If you own and maintain modules to manage your private link resources, see the [module maintainer example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/module_maintainer) to learn how to update your module to support port-mapped architecture while maintaining backward compatibility.
- **Module Users**: If you are using a Terraform module to manage your private link resources, see the [module user example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/module_user) to learn how to upgrade to a module version that supports port-mapped architecture.

For working examples of each architecture, see the [legacy architecture example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp) and the [port-mapped architecture example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped).

### Step 1: Initial Configuration - Legacy Architecture Only

Original configuration with legacy architecture. The number of endpoints is configurable via the `endpoint_count` variable (defaults to 50) and should match your Atlas project's `privateServiceConnectionsPerRegionGroup` setting. See [Set One Project Limit](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-setgrouplimit) for more information on the default value.

```terraform
# Create mongodbatlas_privatelink_endpoint with legacy architecture
resource "mongodbatlas_privatelink_endpoint" "legacy" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  # port_mapping_enabled is not set (defaults to false for legacy architecture)
}

# Create a Google Network
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = "my-network"
}

# Create a Google Sub Network
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Create Google Addresses (required for legacy architecture)
resource "google_compute_address" "legacy" {
  count        = var.endpoint_count
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-legacy${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.legacy]
}

# Create Forwarding rules (required for legacy architecture)
resource "google_compute_forwarding_rule" "legacy" {
  count                 = var.endpoint_count
  target                = mongodbatlas_privatelink_endpoint.legacy.service_attachment_names[count.index]
  project               = google_compute_address.legacy[count.index].project
  region                = google_compute_address.legacy[count.index].region
  name                  = google_compute_address.legacy[count.index].name
  ip_address            = google_compute_address.legacy[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Create mongodbatlas_privatelink_endpoint_service with legacy architecture
resource "mongodbatlas_privatelink_endpoint_service" "legacy" {
  project_id               = mongodbatlas_privatelink_endpoint.legacy.project_id
  private_link_id          = mongodbatlas_privatelink_endpoint.legacy.private_link_id
  provider_name            = "GCP"
  # Note: endpoint_service_id can be any identifier string for legacy architecture.
  # It's used only as an identifier and doesn't need to match any GCP resource name.
  endpoint_service_id      = "legacy-endpoint-group"
  gcp_project_id           = var.gcp_project_id
  # Legacy architecture requires the endpoints list with all endpoints
  dynamic "endpoints" {
    for_each = google_compute_address.legacy

    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.legacy[endpoints.key].name
    }
  }

  depends_on = [google_compute_forwarding_rule.legacy]
}
```

### Step 2: Create Port-Mapped Endpoint (Parallel Setup)

**Resource Naming:** When creating the port-mapped resources, consider using different names to avoid conflicts during the parallel setup phase. For example:
- Legacy: `google_compute_address.default` (with count)
- New: `google_compute_address.new` (single resource)

1. **Add the port-mapped mongodbatlas_privatelink_endpoint alongside your existing legacy resources:**

```terraform
# New: Create mongodbatlas_privatelink_endpoint with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint" "new" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true
}

# New: Create Google Address (1 address for port-mapped architecture)
# Note: Uses existing network and subnet from Step 1
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-mapped-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.new]
}

# New: Create Forwarding Rule (1 rule for port-mapped architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# New: Create mongodbatlas_privatelink_endpoint_service with port-mapped architecture
resource "mongodbatlas_privatelink_endpoint_service" "new" {
  project_id                = mongodbatlas_privatelink_endpoint.new.project_id
  private_link_id           = mongodbatlas_privatelink_endpoint.new.private_link_id
  provider_name             = "GCP"
  endpoint_service_id       = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id            = var.gcp_project_id
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes. You should see:
   - A new `mongodbatlas_privatelink_endpoint.new` resource being created.
   - New GCP resources (1 address, 1 forwarding rule) being created.
   - A new `mongodbatlas_privatelink_endpoint_service.new` resource being created.
   - Your existing legacy resources remain unchanged.

2. Run `terraform apply` to create the port-mapped endpoint resources.

3. **Update your application connection strings** to use the port-mapped endpoint. You can retrieve the connection string from your cluster's private endpoint configuration. **This is when downtime occurs** - update connection strings and restart your applications.

   **Note:** The port-mapped connection strings will have a different format than legacy connection strings. Legacy connection strings use the `pl-0` identifier (e.g., `cluster0-pl-0.a0b1c2.domain.com`), while port-mapped connection strings use the `psc-0` identifier (e.g., `cluster0-psc-0.a0b1c2.domain.com`). **Exception:** For cross-cloud clusters that span a region with a port-mapped endpoint service, the connection string will continue using the `pl-0` identifier. For all other cases (multi-region or single-region clusters), the connection string will use the `psc-0` identifier. Make sure to update all application connection strings accordingly.

4. Test your application connectivity with the port-mapped endpoint to ensure everything works correctly.

5. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove Legacy Resources

Once you have verified that the port-mapped endpoint works correctly and your applications are using it, remove the legacy resources from your configuration:

```terraform
# from Step 2, port-mapped architecture
resource "mongodbatlas_privatelink_endpoint" "new" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true
}

# from Step 1, also used for the port-mapped architecture
resource "google_compute_network" "default" {
  project = var.gcp_project_id
  name    = "my-network"
}

# from Step 1, also used for the port-mapped architecture
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# from Step 2, port-mapped architecture
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-mapped-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.new]
}

# from Step 2, port-mapped architecture
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# from Step 2, port-mapped architecture
resource "mongodbatlas_privatelink_endpoint_service" "new" {
  project_id                  = mongodbatlas_privatelink_endpoint.new.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.new.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id             = var.gcp_project_id
}
```

1. Run `terraform plan` to verify:
   - Legacy endpoint resources are planned for destruction.
   - Legacy GCP resources (addresses and forwarding rules matching your `endpoint_count` variable) are planned for destruction.
   - Only the port-mapped architecture resources remain.
   - No unexpected changes.

2. Run `terraform apply` to finalize the migration. This will:
   - Delete the legacy `mongodbatlas_privatelink_endpoint_service` resource.
   - Delete the legacy `mongodbatlas_privatelink_endpoint` resource.
   - Delete the legacy Google Compute Addresses (number matches your `endpoint_count` variable).
   - Delete the legacy Google Compute Forwarding Rules (number matches your `endpoint_count` variable).

3. Verify that your applications and infrastructure continue to work with the port-mapped endpoint.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

---

## Additional Resources

- [GCP Private Service Connect Documentation](https://www.mongodb.com/docs/atlas/security-private-endpoint/)
- [Private Endpoint Resource Documentation](../resources/privatelink_endpoint.md)
- [Private Endpoint Service Resource Documentation](../resources/privatelink_endpoint_service.md)
- [Legacy Architecture Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp)
- [Port-Mapped Architecture Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)
