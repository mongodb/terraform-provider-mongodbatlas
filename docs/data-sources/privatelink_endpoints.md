---
subcategory: "Private Endpoint Services"
---

# Data Source: mongodbatlas_privatelink_endpoints

`mongodbatlas_privatelink_endpoints` describes all Private Endpoints for a given cloud provider in an Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** Before configuring a private endpoint for a new region in your cluster,
ensure that you review the [Multi-Region Private Endpoints](https://www.mongodb.com/docs/atlas/troubleshoot-private-endpoints/#multi-region-private-endpoints) troubleshooting documentation.

## Example Usage

```terraform
resource "mongodbatlas_privatelink_endpoint" "this" {
  project_id    = var.project_id
  provider_name = "AWS"
  region        = "US_EAST_1"
}

data "mongodbatlas_privatelink_endpoints" "this" {
  project_id    = mongodbatlas_privatelink_endpoint.this.project_id
  provider_name = "AWS"
  depends_on    = [mongodbatlas_privatelink_endpoint.this]
}
```

### Further Examples
- [AWS PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/aws)
- [Azure PrivateLink Endpoint](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/azure)
- [GCP Private Service Connect Endpoint (Port-Mapped Architecture)](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_privatelink_endpoint/gcp-port-mapped)

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `provider_name` - (Required) Cloud provider for which you want to retrieve private endpoint services. Atlas accepts `AWS`, `AZURE` or `GCP`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `results` - A list of Private Endpoints. (see [below for nested schema](#nestedatt--results))

<a id="nestedatt--results"></a>
### Nested Schema for `results`

* `private_link_id` - Unique identifier of the private endpoint.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `status` - Status of the AWS PrivateLink connection.
  Returns one of the following values:
  * `AVAILABLE` 	Atlas created the load balancer and the Private Link Service.
  * `INITIATING` 	Atlas is creating the network load balancer and VPC endpoint service.
  * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
  * `FAILED` 	A system failure has occurred.
  * `DELETING` 	The Private Link service is being deleted.
* `private_endpoints` - All private endpoints that you have added to this Azure Private Link Service.
* `private_link_service_name` - Name of the Azure Private Link Service that Atlas manages.
* `private_link_service_resource_id` - Resource ID of the Azure Private Link Service that Atlas manages.
* `endpoint_group_names` - For port-mapped architectures, this is a list of private endpoint names associated with the private endpoint service. For GCP legacy private endpoint architectures, this is a list of the endpoint group names associated with the private endpoint service.
* `region_name` - Region for the Private Service Connect endpoint service.
* `service_attachment_names` - For port-mapped architecture, this is a list containing one service attachment connected to the private endpoint service. For GCP legacy private endpoint architecture, this is a list of service attachments connected to the private endpoint service (one per Atlas node).
* `port_mapping_enabled` - Flag that indicates whether this resource uses GCP port-mapping. When `true`, it uses the port-mapped architecture. When `false` or unset, it uses the GCP legacy private endpoint architecture. Only applicable for GCP provider.

See [MongoDB Atlas API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/#tag/Private-Endpoint-Services/operation/listGroupPrivateEndpointEndpointService) Documentation for more information.
