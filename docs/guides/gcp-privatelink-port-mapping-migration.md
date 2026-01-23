---
page_title: "Migration Guide: GCP Private Link Legacy Architecture to Port-Based Architecture"
---

# Migration Guide: GCP Private Link Legacy Architecture to Port-Based Architecture

## Overview

This guide explains how to migrate from the legacy GCP Private Service Connect architecture to the new GCP port-based architecture for MongoDB Atlas private link endpoints.

**Important:** The steps in this guide are for migrating Terraform-managed GCP private link endpoint resources (e.g., `mongodbatlas_privatelink_endpoint`, `mongodbatlas_privatelink_endpoint_service`). The legacy architecture requires dedicated resources (customer forwarding rule, service attachment, internal forwarding rule, and instance group) for each Atlas node. The new GCP port-based architecture uses a single set of resources to support up to 1000 nodes through port mapping, enabling direct targeting of specific nodes using only one customer IP address.

**Note:** Migration to the port-based architecture is recommended but **not required**. If you are currently using the legacy architecture, you may continue to do so. This guide is for users who wish to adopt the new architecture for simplified management and reduced resource overhead, but existing legacy configurations will continue to work and be supported.

## Before You Begin

- **Backup your Terraform state file** before making any changes.
- **Test the process in a non-production environment** if possible.
- **Plan for downtime** during the migration, as you will need to delete the old endpoint service and create a new one.
- Ensure you have the necessary GCP permissions to create and delete Compute Addresses and Forwarding Rules.

## Architecture Comparison

### Legacy Architecture
- Requires dedicated resources (customer forwarding rule, service attachment, internal forwarding rule, and instance group) for each Atlas node
- Requires **one Google Compute Address per Atlas node**
- Requires **one Google Compute Forwarding Rule per Atlas node**
- Requires **endpoints list** in the `mongodbatlas_privatelink_endpoint_service` resource (one endpoint per Atlas node)
- Uses `endpoint_service_id` (can be any identifier string)
- Does not use `port_mapping_enabled` (defaults to `false`)
- `endpoint_group_names`: A list of endpoint group names associated with the private endpoint service
- `service_attachment_names`: A list of service attachments connected to the private endpoint service (one per Atlas node)

### GCP Port-Based Architecture (New)
- Uses a single set of resources to support up to 1000 nodes
- Requires only **1 Google Compute Address** (total, not per node)
- Requires only **1 Google Compute Forwarding Rule** (total, not per node)
- Uses `endpoint_service_id` (the forwarding rule name) and `private_endpoint_ip_address` (the IP address)
- Does not require the `endpoints` list
- Requires `port_mapping_enabled = true` on the `mongodbatlas_privatelink_endpoint` resource
- `endpoint_group_names`: A list of private endpoint names associated with the private endpoint service
- `service_attachment_names`: A list of one service attachment connected to the private endpoint service

## Resource Mapping

The following table shows the key differences between the legacy and port-based architectures:

| Aspect | Legacy Architecture | GCP Port-Based Architecture |
|--------|---------------------|------------------------|
| `mongodbatlas_privatelink_endpoint.port_mapping_enabled` | Not set (defaults to `false`) | Must be set to `true` |
| `google_compute_address` count | One per Atlas node | 1 address (total, supports up to 1000 nodes) |
| `google_compute_forwarding_rule` count | One per Atlas node | 1 forwarding rule (total, supports up to 1000 nodes) |
| `mongodbatlas_privatelink_endpoint_service.endpoint_service_id` | Required (can be any identifier string) | Required (should be the forwarding rule name) |
| `mongodbatlas_privatelink_endpoint_service.private_endpoint_ip_address` | Not used | Required (the IP address of the forwarding rule) |
| `mongodbatlas_privatelink_endpoint_service.endpoints` | Required (one endpoint per Atlas node) | Not used |
| `mongodbatlas_privatelink_endpoint_service.gcp_project_id` | Required | Required |
| `mongodbatlas_privatelink_endpoint.endpoint_group_names` | A list of endpoint group names associated with the private endpoint service | A list of private endpoint names associated with the private endpoint service |
| `mongodbatlas_privatelink_endpoint.service_attachment_names` | A list of service attachments connected to the private endpoint service (one per Atlas node) | A list of one service attachment connected to the private endpoint service |

---

## Migration Steps

For complete migration examples showing the step-by-step transition from legacy to port-based architecture, see the [GCP Private Link migration example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_gcp_privatelink_legacy_to_port_based).

For working examples of each architecture, see the [legacy architecture example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp) and the [port-based architecture example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-based).

### Step 1: Initial Configuration - Legacy Architecture Only

Original configuration with legacy architecture (50 endpoints):

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  # port_mapping_enabled is not set (defaults to false for legacy architecture)
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
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

# Create Google 50 Addresses (required for legacy architecture)
resource "google_compute_address" "default" {
  count        = 50
  project      = google_compute_subnetwork.default.project
  name         = "tf-test${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test]
}

# Create 50 Forwarding rules (required for legacy architecture)
resource "google_compute_forwarding_rule" "default" {
  count                 = 50
  target                = mongodbatlas_privatelink_endpoint.test.service_attachment_names[count.index]
  project               = google_compute_address.default[count.index].project
  region                = google_compute_address.default[count.index].region
  name                  = google_compute_address.default[count.index].name
  ip_address            = google_compute_address.default[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id               = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id          = mongodbatlas_privatelink_endpoint.test.private_link_id
  provider_name            = "GCP"
  # Note: endpoint_service_id can be any identifier string for legacy architecture.
  # It's used only as an identifier and doesn't need to match any GCP resource name.
  endpoint_service_id      = "legacy-endpoint-group"
  gcp_project_id           = var.gcp_project_id
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
  # Legacy architecture requires the endpoints list with all 50 endpoints
  dynamic "endpoints" {
    for_each = google_compute_address.default

    content {
      ip_address    = endpoints.value["address"]
      endpoint_name = google_compute_forwarding_rule.default[endpoints.key].name
    }
  }

  depends_on = [google_compute_forwarding_rule.default]
}
```

### Step 2: Create New Port-Based Endpoint (Parallel Setup)

**Important:** You cannot modify an existing `mongodbatlas_privatelink_endpoint` to enable port mapping. You must create a new endpoint with `port_mapping_enabled = true`. During this step, you will have both architectures running in parallel.

1. **Create a new private link endpoint with port mapping enabled:**

```terraform
# New endpoint with port-based architecture
resource "mongodbatlas_privatelink_endpoint" "test_new" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true # Enable new GCP port-based architecture
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
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

# Create Google Address (1 address for new GCP port-based architecture)
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-based-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_new]
}

# Create Forwarding Rule (1 rule for new GCP port-based architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.test_new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

# Create MongoDB Atlas Private Endpoint Service for new architecture
resource "mongodbatlas_privatelink_endpoint_service" "test_new" {
  project_id                = mongodbatlas_privatelink_endpoint.test_new.project_id
  private_link_id           = mongodbatlas_privatelink_endpoint.test_new.private_link_id
  provider_name             = "GCP"
  endpoint_service_id       = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id            = var.gcp_project_id
  delete_on_create_timeout  = true
  timeouts {
    create = "10m"
    delete = "10m"
  }

  depends_on = [google_compute_forwarding_rule.new]
}
```

**Apply and test:**

1. Run `terraform plan` to review the changes. You should see:
   - A new `mongodbatlas_privatelink_endpoint.test_new` resource being created
   - New GCP resources (1 address, 1 forwarding rule) being created
   - A new `mongodbatlas_privatelink_endpoint_service.test_new` resource being created
   - The old resources remain unchanged

2. Run `terraform apply` to create the new port-based endpoint resources.

3. **Update your application connection strings** to use the new endpoint. You can retrieve the connection string from your cluster's private endpoint configuration.

4. Test your application connectivity with the new endpoint to ensure everything works correctly.

5. Re-run `terraform plan` to ensure you have no unexpected changes: `No changes. Your infrastructure matches the configuration.`

### Step 3: Final State - Remove Legacy Resources

Once you have verified that the new port-based endpoint works correctly and your applications are using it, remove the legacy resources from your configuration:

```terraform
# New endpoint with port-based architecture
resource "mongodbatlas_privatelink_endpoint" "test_new" {
  project_id               = var.project_id
  provider_name            = "GCP"
  region                   = var.gcp_region
  port_mapping_enabled     = true # Enable new GCP port-based architecture
  delete_on_create_timeout = true
  timeouts {
    create = "10m"
    delete = "10m"
  }
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

# Create Google Address (1 address for new GCP port-based architecture)
resource "google_compute_address" "new" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-test-port-based-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.100"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.test_new]
}

# Create Forwarding Rule (1 rule for new GCP port-based architecture)
resource "google_compute_forwarding_rule" "new" {
  target                = mongodbatlas_privatelink_endpoint.test_new.service_attachment_names[0]
  project               = google_compute_address.new.project
  region                = google_compute_address.new.region
  name                  = google_compute_address.new.name
  ip_address            = google_compute_address.new.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "test_new" {
  project_id                = mongodbatlas_privatelink_endpoint.test_new.project_id
  private_link_id           = mongodbatlas_privatelink_endpoint.test_new.private_link_id
  provider_name             = "GCP"
  endpoint_service_id       = google_compute_forwarding_rule.new.name
  private_endpoint_ip_address = google_compute_address.new.address
  gcp_project_id            = var.gcp_project_id
  delete_on_create_timeout  = true
  timeouts {
    create = "10m"
    delete = "10m"
  }

  depends_on = [google_compute_forwarding_rule.new]
}
```

1. Run `terraform plan` to verify:
   - Legacy endpoint resources are planned for destruction
   - Legacy GCP resources (50 addresses, 50 forwarding rules) are planned for destruction
   - Only the port-based architecture resources remain
   - No unexpected changes

2. Run `terraform apply` to finalize the migration. This will:
   - Delete the legacy `mongodbatlas_privatelink_endpoint_service` resource
   - Delete the legacy `mongodbatlas_privatelink_endpoint` resource
   - Delete the 50 legacy Google Compute Addresses
   - Delete the 50 legacy Google Compute Forwarding Rules

3. Verify that your applications and infrastructure continue to work with the new port-based endpoint.

4. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

---

## Important Considerations

### Cannot Modify Existing Endpoint

**You cannot modify an existing `mongodbatlas_privatelink_endpoint` to enable port mapping.** The `port_mapping_enabled` attribute must be set when the endpoint is first created. If you need to migrate, you must:

1. Create a new endpoint with `port_mapping_enabled = true`
2. Create new GCP resources (1 address, 1 forwarding rule)
3. Create a new endpoint service linking to the new endpoint
4. Update your application connection strings
5. Delete the old endpoint and resources

### Connection String Updates & Downtime Considerations

After migrating to the new architecture, you will need to update your application connection strings. The connection string will reference the new endpoint. You can retrieve the updated connection string from your cluster's private endpoint configuration. This may result in downtime. Plan for this accordingly.

### Resource Naming

When creating the new resources, consider using different names to avoid conflicts during the parallel setup phase. For example:
- Legacy: `google_compute_address.default` (with count)
- New: `google_compute_address.new` (single resource)

---

## Additional Resources

- [GCP Private Service Connect Documentation](https://www.mongodb.com/docs/atlas/security-private-endpoint/)
- [Private Endpoint Resource Documentation](../resources/privatelink_endpoint.md)
- [Private Endpoint Service Resource Documentation](../resources/privatelink_endpoint_service.md)
- [Legacy Architecture Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp)
- [Port-Based Architecture Example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/mongodbatlas_privatelink_endpoint/gcp-port-based)
