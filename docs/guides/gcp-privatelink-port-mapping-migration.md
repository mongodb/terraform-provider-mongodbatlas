---
page_title: "Migration Guide: GCP Private Service Connect to Port-Mapped Architecture"
---

# Migration Guide: GCP Private Service Connect to Port-Mapped Architecture

**Objective**: Migrate from the legacy GCP Private Service Connect architecture (one service attachment per Atlas node) to the port-mapped architecture (one service attachment for up to 150 nodes) for `mongodbatlas_privatelink_endpoint` and `mongodbatlas_privatelink_endpoint_service` resources.

-> **Note:** Migration to the port-mapped architecture is recommended but **not required**. You may continue using the legacy architecture.

## Why Migrate to Port-Mapped Architecture?

The legacy architecture has two main limitations:

1. **IP Exhaustion**: Atlas defaults to 50 private service connections per region group (50 forwarding rules and 50 IP addresses), which can lead to IP address exhaustion in your GCP project.

2. **Static Configuration**: Changing the number of private service connections per region group requires a full private service connect redeployment, causing friction when changing cluster configurations.

The port-mapped architecture addresses these limitations by using one service attachment that can be mapped to up to 150 nodes via ports designated per node, requiring only 1 Google Compute Address and 1 Google Compute Forwarding Rule.

## Architecture Comparison

| Resource | Legacy | Port-Mapped |
|----------|--------|-------------|
| `google_compute_address` | One per node | 1 total |
| `google_compute_forwarding_rule` | One per node | 1 total |

**Key attribute changes:**

- `mongodbatlas_privatelink_endpoint.port_mapping_enabled`: Set to `true` (legacy defaults to `false`).
- `mongodbatlas_privatelink_endpoint_service.endpoint_service_id`: Forwarding rule name (legacy: endpoint group name).
- `mongodbatlas_privatelink_endpoint_service.private_endpoint_ip_address`: Required (not used in legacy).
- `mongodbatlas_privatelink_endpoint_service.endpoints`: Not used (required in legacy).

## Best Practices Before Migrating

- Backup your Terraform state file before making any changes.
- Test the process in a non-production environment if possible.
- Ensure you have the necessary GCP permissions to create and delete Compute Addresses and Forwarding Rules.

-> **Note:** You cannot modify an existing `mongodbatlas_privatelink_endpoint` to enable port mapping. You must create new resources alongside the legacy ones, then remove the legacy resources after migration.

-> **Note:** Downtime occurs when updating application connection strings, not during Terraform operations. You can run both architectures in parallel during the transition.

---

## Migration Steps

For complete migration examples showing the step-by-step transition from legacy to port-mapped architecture, see the [GCP Private Link migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture).

- **Direct Resource Management**: If you are managing the resources directly (not using modules), see the [basic migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/basic).
- **Module Maintainers**: If you own and maintain modules to manage your private link resources, see the [module maintainer example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/module_maintainer) to learn how to update your module to support port-mapped architecture while maintaining backward compatibility.
- **Module Users**: If you are using a Terraform module to manage your private link resources, see the [module user example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_private_link_to_port_mapped_architecture/module_user) to learn how to upgrade to a module version that supports port-mapped architecture.

For a working example of the port-mapped architecture, see the [port-mapped architecture example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped).

### 1) Initial Configuration (Legacy Architecture)

Your existing legacy configuration typically includes multiple GCP addresses and forwarding rules (one per Atlas node). The count defaults to 50 based on your Atlas project's `privateServiceConnectionsPerRegionGroup` setting.

```hcl
resource "mongodbatlas_privatelink_endpoint" "legacy" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
  # port_mapping_enabled not set (defaults to false)
}

resource "google_compute_address" "legacy" {
  count        = var.legacy_endpoint_count
  project      = var.gcp_project_id
  name         = "legacy-address-${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  region       = var.gcp_region
}

resource "google_compute_forwarding_rule" "legacy" {
  count                 = var.legacy_endpoint_count
  project               = var.gcp_project_id
  name                  = google_compute_address.legacy[count.index].name
  target                = mongodbatlas_privatelink_endpoint.legacy.service_attachment_names[count.index]
  ip_address            = google_compute_address.legacy[count.index].id
  network               = google_compute_network.default.id
  region                = var.gcp_region
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "legacy" {
  project_id          = mongodbatlas_privatelink_endpoint.legacy.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.legacy.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = "legacy-endpoint-group"
  gcp_project_id      = var.gcp_project_id

  dynamic "endpoints" {
    for_each = google_compute_address.legacy
    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.legacy[endpoints.key].name
    }
  }
}
```

### 2) Create Port-Mapped Endpoint (Parallel Setup)

Add the port-mapped resources alongside your existing legacy resources. Use different resource names (e.g., `port_mapped` vs `legacy`) to avoid conflicts.

```hcl
resource "mongodbatlas_privatelink_endpoint" "port_mapped" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true
}

resource "google_compute_address" "port_mapped" {
  project      = var.gcp_project_id
  name         = "port-mapped-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  region       = var.gcp_region
}

resource "google_compute_forwarding_rule" "port_mapped" {
  project               = var.gcp_project_id
  name                  = google_compute_address.port_mapped.name
  target                = mongodbatlas_privatelink_endpoint.port_mapped.service_attachment_names[0]
  ip_address            = google_compute_address.port_mapped.id
  network               = google_compute_network.default.id
  region                = var.gcp_region
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "port_mapped" {
  project_id                  = mongodbatlas_privatelink_endpoint.port_mapped.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.port_mapped.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.port_mapped.name
  private_endpoint_ip_address = google_compute_address.port_mapped.address
  gcp_project_id              = var.gcp_project_id
}
```

**Apply and test:**

1. Run `terraform plan` to verify new port-mapped resources will be created and legacy resources remain unchanged.

2. Run `terraform apply` to create the port-mapped resources.

3. **Update your application connection strings.** This is when downtime occurs. Retrieve the new connection string from your cluster's private endpoint configuration.

   -> **Note:** Connection string format changes from `pl-0` (e.g., `cluster0-pl-0.a0b1c2.domain.com`) to `psc-0` (e.g., `cluster0-psc-0.a0b1c2.domain.com`). For single-region and multi-region clusters, the connection string uses `psc-0`. **Exception:** Cross-cloud clusters spanning a region with a port-mapped endpoint continue using `pl-0`. Make sure to update all application connection strings accordingly.

4. Test application connectivity with the port-mapped endpoint.

5. Run `terraform plan` to confirm: `No changes. Your infrastructure matches the configuration.`

### 3) Remove Legacy Resources

Once you have verified that the port-mapped endpoint works correctly and your applications are using it:

1. Remove the legacy resources from your Terraform configuration:
   - `mongodbatlas_privatelink_endpoint.legacy`
   - `mongodbatlas_privatelink_endpoint_service.legacy`
   - `google_compute_address.legacy`
   - `google_compute_forwarding_rule.legacy`

2. Keep the shared resources (`google_compute_network`, `google_compute_subnetwork`) and all port-mapped resources from Step 2.

3. Run `terraform plan` to verify legacy resources will be destroyed and port-mapped resources remain unchanged.

4. Run `terraform apply` to delete the legacy resources.

5. Run `terraform plan` again to confirm: `No changes. Your infrastructure matches the configuration.`

---

## Additional Resources

- [GCP Private Service Connect Documentation](https://www.mongodb.com/docs/atlas/security-private-endpoint/)
- [Private Endpoint Resource Documentation](../resources/privatelink_endpoint.md)
- [Private Endpoint Service Resource Documentation](../resources/privatelink_endpoint_service.md)
- [Port-Mapped Architecture Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)
