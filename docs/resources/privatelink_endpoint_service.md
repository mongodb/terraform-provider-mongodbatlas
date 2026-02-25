---
subcategory: "Private Endpoint Services"
---

# Resource: mongodbatlas_privatelink_endpoint_service

`mongodbatlas_privatelink_endpoint_service` provides a Private Endpoint Interface Link resource. This represents a Private Endpoint Interface Link, which adds one [Interface Endpoint](https://www.mongodb.com/docs/atlas/security-private-endpoint/#private-endpoint-concepts) to a private endpoint connection in an Atlas project.

~> **IMPORTANT:** This resource links your cloud provider's Private Endpoint to the MongoDB Atlas Private Endpoint Service. It does not create the service itself (this is done by `mongodbatlas_privatelink_endpoint`). You first create the service in Atlas with `mongodbatlas_privatelink_endpoint`, then the endpoint is created in your cloud provider, and you link them together with the `mongodbatlas_privatelink_endpoint_service` resource.

The [private link Terraform module](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/private-endpoint/mongodbatlas/latest) makes use of this resource and simplifies its use.

-> **NOTE:** You must have Organization Owner or Project Owner role. Create and delete operations wait for all clusters on the project to IDLE to ensure the latest connection strings can be retrieved (default timeout: 2hrs).

~> **IMPORTANT:** For GCP, MongoDB encourages customers to use the port-mapped architecture by setting `port_mapping_enabled = true` on the `mongodbatlas_privatelink_endpoint` resource. This architecture uses a single set of resources to support up to 150 nodes. The legacy architecture requires dedicated resources for each Atlas node, which can lead to IP address exhaustion. For migration guidance, see the [GCP Private Service Connect to Port-Mapped Architecture](../guides/gcp-privatelink-port-mapping-migration.md).

## Example with AWS

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = "<PROJECT_ID>"
  provider_name = "AWS"
  region        = "US_EAST_1"
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = "vpc-7fc0a543"
  service_name       = mongodbatlas_privatelink_endpoint.this.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
  endpoint_service_id = aws_vpc_endpoint.ptfe_service.id
  provider_name       = "AWS"
}
```

## Example with Azure

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AZURE"
  region        = "eastus2"
}

resource "azurerm_private_endpoint" "this" {
  name                = "endpoint-this"
  location            = data.azurerm_resource_group.this.location
  resource_group_name = var.resource_group_name
  subnet_id           = azurerm_subnet.this.id
  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint.this.private_link_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint.this.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Azure Private Link this"
  }

}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id                  = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.this.private_link_id
  endpoint_service_id         = azurerm_private_endpoint.this.id
  private_endpoint_ip_address = azurerm_private_endpoint.this.private_service_connection.0.private_ip_address
  provider_name               = "AZURE"
}
```

## Example with GCP (Legacy Architecture)

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
}

# Create a Google Network
resource "google_compute_network" "default" {
  project                 = var.gcp_project_id
  name                    = "my-network"
  auto_create_subnetworks = false
}

# Create a Google Sub Network
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Create Google 50 Addresses (required for GCP legacy private endpoint architecture)
resource "google_compute_address" "default" {
  count        = 50
  project      = google_compute_subnetwork.default.project
  name         = "tf-this${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = var.gcp_region

  depends_on = [mongodbatlas_privatelink_endpoint.this]
}

# Create 50 Forwarding rules (required for GCP legacy private endpoint architecture)
resource "google_compute_forwarding_rule" "default" {
  count                 = 50
  target                = mongodbatlas_privatelink_endpoint.this.service_attachment_names[count.index]
  project               = google_compute_address.default[count.index].project
  region                = google_compute_address.default[count.index].region
  name                  = google_compute_address.default[count.index].name
  ip_address            = google_compute_address.default[count.index].id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}


resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.this.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = google_compute_network.default.name
  gcp_project_id      = var.gcp_project_id

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

## Example with GCP (Port-Mapped Architecture)

The port-mapped architecture uses port mapping to reduce resource provisioning. In the GCP legacy private endpoint architecture, service attachments were mapped 1:1 with Atlas nodes (one service attachment per node). In the port-mapped architecture, regardless of cloud provider, one service attachment can be mapped to up to 150 nodes via ports designated per node, enabling direct targeting of specific nodes using only one customer IP address. Enable it by setting `port_mapping_enabled = true` on the `mongodbatlas_privatelink_endpoint` resource.

**Important:** For the port-mapped architecture, use `endpoint_service_id` (the forwarding rule name) and `private_endpoint_ip_address` (the IP address). The `endpoints` list is no longer used for the port-mapped architecture.

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id           = var.project_id
  provider_name        = "GCP"
  region               = var.gcp_region
  port_mapping_enabled = true # Enable port-mapped architecture
}

# Create a Google Network
resource "google_compute_network" "default" {
  project                 = var.gcp_project_id
  name                    = "my-network"
  auto_create_subnetworks = false
}

# Create a Google Sub Network
resource "google_compute_subnetwork" "default" {
  project       = google_compute_network.default.project
  name          = "my-subnet"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.gcp_region
  network       = google_compute_network.default.id
}

# Create Google Address (1 address for port-mapped architecture)
resource "google_compute_address" "default" {
  project      = google_compute_subnetwork.default.project
  name         = "tf-this-psc-endpoint"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.1"
  region       = google_compute_subnetwork.default.region

  depends_on = [mongodbatlas_privatelink_endpoint.this]
}

# Create Forwarding Rule (1 rule for port-mapped architecture)
resource "google_compute_forwarding_rule" "default" {
  target                = mongodbatlas_privatelink_endpoint.this.service_attachment_names[0]
  project               = google_compute_address.default.project
  region                = google_compute_address.default.region
  name                  = google_compute_address.default.name
  ip_address            = google_compute_address.default.id
  network               = google_compute_network.default.id
  load_balancing_scheme = ""
}

resource "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id                  = mongodbatlas_privatelink_endpoint.this.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.this.private_link_id
  provider_name               = "GCP"
  endpoint_service_id         = google_compute_forwarding_rule.default.name
  private_endpoint_ip_address = google_compute_address.default.address
  gcp_project_id              = var.gcp_project_id

  depends_on = [google_compute_forwarding_rule.default]
}

```

### Further Examples
- [AWS PrivateLink Endpoint and Service](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)
- [Azure Private Link Endpoint and Service](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/azure)
- [GCP Private Service Connect Endpoint and Service (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `endpoint_service_id` (String)
- `private_link_id` (String)
- `project_id` (String)
- `provider_name` (String)

### Optional

- `delete_on_create_timeout` (Boolean) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.
- `endpoints` (Block List) (see [below for nested schema](#nestedblock--endpoints))
- `gcp_project_id` (String)
- `private_endpoint_ip_address` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))

### Read-Only

- `aws_connection_status` (String)
- `azure_status` (String)
- `delete_requested` (Boolean)
- `endpoint_group_name` (String)
- `error_message` (String)
- `gcp_endpoint_status` (String) Status of the GCP endpoint. Only populated for port-mapped architecture.
- `gcp_status` (String)
- `id` (String) The ID of this resource.
- `interface_endpoint_id` (String)
- `port_mapping_enabled` (Boolean) Flag that indicates whether the underlying `privatelink_endpoint` resource uses GCP port-mapping. This is a read-only attribute that reflects the architecture type. When `true`, the endpoint service uses the port-mapped architecture. When `false`, it uses the GCP legacy private endpoint architecture. Only applicable for GCP provider.
- `private_endpoint_connection_name` (String)
- `private_endpoint_resource_id` (String)

<a id="nestedblock--endpoints"></a>
### Nested Schema for `endpoints`

Optional:

- `endpoint_name` (String)
- `ip_address` (String)

Read-Only:

- `status` (String)


<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)

## Import
Private Endpoint Link Connection can be imported using project ID, private link ID, endpoint service ID, and provider name, in the format `{project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint_service.this 1112222b3bf99403840e8934--3242342343112--vpce-4242342343--AWS
```

For more information, see:
- [MongoDB API Private Endpoint Link Connection](https://www.mongodb.com/docs/api/doc/atlas-admin-api-v2/operation/operation-creategroupprivateendpointendpointserviceendpoint) for detailed arguments and attributes.
- [Set Up a Private Endpoint](https://www.mongodb.com/docs/atlas/security-private-endpoint/) for general guidance on private endpoints in MongoDB Atlas.
