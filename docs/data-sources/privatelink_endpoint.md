---
subcategory: "Private Endpoint Services"
---

# Data Source: mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` describes a Private Endpoint. This represents a Private Endpoint Connection to retrieve details regarding a private endpoint by id in an Atlas project

-> **NOTE:** Groups and projects are synonymous terms. The official documentation uses `group_id`.

~> **IMPORTANT:** Before configuring a private endpoint for a new region in your cluster, review the [Multi-Region Private Endpoints](https://www.mongodb.com/docs/atlas/troubleshoot-private-endpoints/#multi-region-private-endpoints) troubleshooting documentation.

## Example Usage

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = "US_EAST_1"
}

data "mongodbatlas_privatelink_endpoint" "this" {
	project_id      = mongodbatlas_privatelink_endpoint.this.project_id
	private_link_id = mongodbatlas_privatelink_endpoint.this.private_link_id
    provider_name = "AWS"
}
```

### Further Examples
- [AWS PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws)
- [Azure PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/azure)
- [GCP Private Service Connect Endpoint (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the private endpoint that you want to retrieve.
* `provider_name` - (Required) Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts `AWS`, `AZURE`, or `GCP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns `null` while Atlas creates the endpoint service.
* `error_message` - Error message for the private endpoint connection. Returns `null` if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `status` - Status of the AWS PrivateLink connection.
  Returns one of the following values:
  * `AVAILABLE` - Atlas created the load balancer and the Private Link Service.
  * `INITIATING` - Atlas is creating the network load balancer and VPC endpoint service.
  * `WAITING_FOR_USER` - The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
  * `FAILED` - A system failure occurred.
  * `DELETING` - Atlas is deleting the Private Link service.
* `private_endpoints` - All private endpoints that you have added to this Azure Private Link Service.
* `private_link_service_name` - Name of the Azure Private Link Service that Atlas manages.
* `private_link_service_resource_id` - Resource ID of the Azure Private Link Service that Atlas manages.
* `endpoint_group_names` - List of private endpoint names associated with the private endpoint service for port-mapped architectures. For GCP legacy private endpoint architectures, this is a list of endpoint group names associated with the private endpoint service.
* `region_name` - Region for the Private Service Connect endpoint service.
* `service_attachment_names` - List containing one service attachment connected to the private endpoint service for port-mapped architecture. For GCP legacy private endpoint architecture, this is a list of service attachments connected to the private endpoint service (one per Atlas node). Returns an empty list while Atlas creates the service attachments.
* `port_mapping_enabled` - Flag that indicates whether this resource uses GCP port-mapping. When `true`, the resource uses port-mapped architecture. When `false` or unset, the resource uses GCP legacy private endpoint architecture. Only applicable for GCP provider.

See the [MongoDB Atlas API documentation](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/) for more information.
