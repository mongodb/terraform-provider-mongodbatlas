---
subcategory: "Private Endpoint Services"
---

# Data Source: mongodbatlas_privatelink_endpoint_service

`mongodbatlas_privatelink_endpoint_service` describes a Private Endpoint Link. This represents a Private Endpoint Link Connection that wants to retrieve details in an Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

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

data "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint_service.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint_service.this.private_link_id
  endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
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

data "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint_service.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint_service.this.private_link_id
  endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
  provider_name       = "AZURE"
}
```

## Example with GCP (Legacy Architecture)

```terraform
data "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint_service.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint_service.this.private_link_id
  endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
  provider_name       = "GCP"
}
```

## Example with GCP (Port-Mapped Architecture)

The port-mapped architecture uses port mapping to reduce resource provisioning. In the GCP legacy private endpoint architecture, service attachments were mapped 1:1 with Atlas nodes (one service attachment per node). In the port-mapped architecture, regardless of cloud provider, one service attachment can be mapped to up to 150 nodes via ports designated per node, enabling direct targeting of specific nodes using only one customer IP address. Enable it by setting `port_mapping_enabled = true` on the `mongodbatlas_privatelink_endpoint` resource.

**Important:** For the port-mapped architecture, use `endpoint_service_id` (the forwarding rule name) and `private_endpoint_ip_address` (the IP address). The `endpoints` list is no longer used for the port-mapped architecture.

```terraform
data "mongodbatlas_privatelink_endpoint_service" "this" {
  project_id          = mongodbatlas_privatelink_endpoint_service.this.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint_service.this.private_link_id
  endpoint_service_id = mongodbatlas_privatelink_endpoint_service.this.endpoint_service_id
  provider_name       = "GCP"
}
```

### Available complete examples
- [Setup private connection to a MongoDB Atlas Cluster with AWS VPC](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.6.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)
- [GCP Private Service Connect Endpoint and Service (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.6.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the private endpoint service for which you want to retrieve a private endpoint.
* `endpoint_service_id` - (Required) Unique identifier of the interface endpoint you created in your VPC. For `AWS` and `AZURE`, this is the interface endpoint identifier. For `GCP` port-mapped architecture, this is the forwarding rule name. For `GCP` legacy private endpoint architecture, this is the endpoint group name.
* `provider_name` - (Required) Cloud provider for which you want to retrieve a private endpoint. Atlas accepts `AWS`, `AZURE` or `GCP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `interface_endpoint_id` - Unique identifier of the interface endpoint.
* `private_endpoint_connection_name` - Name of the connection for this private endpoint that Atlas generates.
* `private_endpoint_ip_address` - Private IP address of the private endpoint network interface. For port-mapped architecture, this is required and is the IP address of the forwarding rule. For GCP legacy private endpoint architecture, this is not used.
* `private_endpoint_resource_id` - Unique identifier of the private endpoint.
* `delete_requested` - Indicates if Atlas received a request to remove the interface endpoint from the private endpoint connection.
* `error_message` - Error message pertaining to the interface endpoint. Returns null if there are no errors.
* `aws_connection_status` - Status of the interface endpoint for AWS.
  Returns one of the following values:
  * `NONE` - Atlas created the network load balancer and VPC endpoint service, but AWS hasn’t yet created the VPC endpoint.
  * `PENDING_ACCEPTANCE` - AWS has received the connection request from your VPC endpoint to the Atlas VPC endpoint service.
  * `PENDING` - AWS is establishing the connection between your VPC endpoint and the Atlas VPC endpoint service.
  * `AVAILABLE` - Atlas VPC resources are connected to the VPC endpoint in your VPC. You can connect to Atlas clusters in this region using AWS PrivateLink.
  * `REJECTED` - AWS failed to establish a connection between Atlas VPC resources to the VPC endpoint in your VPC.
  * `DELETING` - Atlas is removing the interface endpoint from the private endpoint connection.
* `azure_status` - Status of the interface endpoint for AZURE.
  Returns one of the following values:
  * `INITIATING` - Atlas has not yet accepted the connection to your private endpoint.
  * `AVAILABLE` - Atlas approved the connection to your private endpoint.
  * `FAILED` - Atlas failed to accept the connection your private endpoint.
  * `DELETING` - Atlas is removing the connection to your private endpoint from the Private Link service.
* `gcp_status` - Status of the interface endpoint.
  Returns one of the following values:
  * `INITIATING` - Atlas has not yet accepted the connection to your private endpoint.
  * `AVAILABLE` - Atlas approved the connection to your private endpoint.
  * `FAILED` - Atlas failed to accept the connection your private endpoint.
  * `DELETING` - Atlas is removing the connection to your private endpoint from the Private Link service.
* `gcp_endpoint_status` - Status of the individual endpoint. Only populated for port-mapped architecture. Returns one of the following values: `INITIATING`, `AVAILABLE`, `FAILED`, `DELETING`.
* `endpoints` - Collection of individual private endpoints that comprise your network endpoint group. Only populated for GCP legacy private endpoint architecture.
  * `endpoint_name` - Forwarding rule that corresponds to the endpoint you created.
  * `ip_address` - Private IP address of the network endpoint group you created.
  * `status` - Status of the endpoint. Atlas returns one of the [values shown above](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-create-one/#std-label-ref-status-field).
* `port_mapping_enabled` - Flag that indicates whether the underlying `privatelink_endpoint` resource uses GCP port-mapping. This is a read-only attribute that reflects the architecture type. When `true`, the endpoint service uses the port-mapped architecture. When `false`, it uses the GCP legacy private endpoint architecture. Only applicable for GCP provider.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-get-one/) Documentation for more information.
