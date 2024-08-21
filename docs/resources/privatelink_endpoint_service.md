# Resource: mongodbatlas_privatelink_endpoint_service

`mongodbatlas_privatelink_endpoint_service` provides a Private Endpoint Interface Link resource. This represents a Private Endpoint Interface Link, which adds one [Interface Endpoint](https://www.mongodb.com/docs/atlas/security-private-endpoint/#private-endpoint-concepts) to a private endpoint connection in an Atlas project.

The [private link Terraform module](https://registry.terraform.io/modules/terraform-mongodbatlas-modules/private-endpoint/mongodbatlas/latest) makes use of this resource and simplifies its use.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** Create and delete wait for all clusters on the project to IDLE in order for their operations to complete. This ensures the latest connection strings can be retrieved following creation or deletion of this resource. Default timeout is 2hrs.

## Example with AWS

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT_ID>"
  provider_name = "AWS"
  region        = "US_EAST_1"
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = "vpc-7fc0a543"
  service_name       = mongodbatlas_privatelink_endpoint.test.endpoint_service_name
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id          = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test.private_link_id
  endpoint_service_id = aws_vpc_endpoint.ptfe_service.id
  provider_name       = "AWS"
}
```

## Example with Azure

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = var.project_id
  provider_name = "AZURE"
  region        = "eastus2"
}

resource "azurerm_private_endpoint" "test" {
  name                = "endpoint-test"
  location            = data.azurerm_resource_group.test.location
  resource_group_name = var.resource_group_name
  subnet_id           = azurerm_subnet.test.id
  private_service_connection {
    name                           = mongodbatlas_privatelink_endpoint.test.private_link_service_name
    private_connection_resource_id = mongodbatlas_privatelink_endpoint.test.private_link_service_resource_id
    is_manual_connection           = true
    request_message                = "Azure Private Link test"
  }

}

resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id                  = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id             = mongodbatlas_privatelink_endpoint.test.private_link_id
  endpoint_service_id         = azurerm_private_endpoint.test.id
  private_endpoint_ip_address = azurerm_private_endpoint.test.private_service_connection.0.private_ip_address
  provider_name               = "AZURE"
}
```

## Example with GCP

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = var.project_id
  provider_name = "GCP"
  region        = var.gcp_region
}

# Create a Google Network
resource "google_compute_network" "default" {
  project = var.gcp_project
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

# Create Google 50 Addresses
resource "google_compute_address" "default" {
  count        = 50
  project      = google_compute_subnetwork.default.project
  name         = "tf-test${count.index}"
  subnetwork   = google_compute_subnetwork.default.id
  address_type = "INTERNAL"
  address      = "10.0.42.${count.index}"
  region       = var.gcp_region

  depends_on = [mongodbatlas_privatelink_endpoint.test]
}

# Create 50 Forwarding rules
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
  project_id          = mongodbatlas_privatelink_endpoint.test.project_id
  private_link_id     = mongodbatlas_privatelink_endpoint.test.private_link_id
  provider_name       = "GCP"
  endpoint_service_id = google_compute_network.default.name
  gcp_project_id      = var.gcp_project

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

### Available complete examples
- [Setup private connection to a MongoDB Atlas Cluster with AWS VPC](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/aws-privatelink-endpoint/cluster)

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the `AWS` or `AZURE` PrivateLink connection which is created by `mongodbatlas_privatelink_endpoint` resource.
* `endpoint_service_id` - (Required) Unique identifier of the interface endpoint you created in your VPC with the `AWS`, `AZURE` or `GCP` resource.
* `provider_name` - (Required) Cloud provider for which you want to create a private endpoint. Atlas accepts `AWS`, `AZURE` or `GCP`.
* `private_endpoint_ip_address` - (Optional) Private IP address of the private endpoint network interface you created in your Azure VNet. Only for `AZURE`.
* `gcp_project_id` - (Optional) Unique identifier of the GCP project in which you created your endpoints. Only for `GCP`.
* `endpoints` - (Optional) Collection of individual private endpoints that comprise your endpoint group. Only for `GCP`. See below.
* `timeouts`- (Optional) The duration of time to wait for Private Endpoint Service to be created or deleted. The timeout value is defined by a signed sequence of decimal numbers with an time unit suffix such as: `1h45m`, `300s`, `10m`, .... The valid time units are:  `ns`, `us` (or `µs`), `ms`, `s`, `m`, `h`. The default timeout for Private Endpoint create & delete is `2h`. Learn more about timeouts [here](https://www.terraform.io/plugin/sdkv2/resources/retries-and-customizable-timeouts).

### `endpoints`
* `ip_address` - (Optional) Private IP address of the endpoint you created in GCP.
* `endpoint_name` - (Optional) Forwarding rule that corresponds to the endpoint you created in GCP.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `interface_endpoint_id` - Unique identifier of the interface endpoint.
* `private_endpoint_connection_name` - Name of the connection for this private endpoint that Atlas generates.
* `private_endpoint_ip_address` - Private IP address of the private endpoint network interface.
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
* `gcp_status` - Status of the interface endpoint for GCP.
  Returns one of the following values:
  * `INITIATING` - Atlas has not yet accepted the connection to your private endpoint.
  * `AVAILABLE` - Atlas approved the connection to your private endpoint.
  * `FAILED` - Atlas failed to accept the connection your private endpoint.
  * `DELETING` - Atlas is removing the connection to your private endpoint from the Private Link service.
* `endpoint_group_name` - (Optional) Unique identifier of the endpoint group. The endpoint group encompasses all of the endpoints that you created in GCP.
* `endpoints` - Collection of individual private endpoints that comprise your network endpoint group.
  * `status` - Status of the endpoint. Atlas returns one of the [values shown above](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-create-one/#std-label-ref-status-field).

## Import
Private Endpoint Link Connection can be imported using project ID and username, in the format `{project_id}--{private_link_id}--{endpoint_service_id}--{provider_name}`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint_service.test 1112222b3bf99403840e8934--3242342343112--vpce-4242342343--AWS
```

See detailed information for arguments and attributes: [MongoDB API Private Endpoint Link Connection](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-create-one/)
