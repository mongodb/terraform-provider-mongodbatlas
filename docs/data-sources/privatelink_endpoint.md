---
subcategory: "Private Endpoint Services"
---

# Data Source: mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` describes a Private Endpoint. This represents a Private Endpoint Connection to retrieve details regarding a private endpoint by id in an Atlas project

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

~> **IMPORTANT:** Before configuring a private endpoint for a new region in your cluster,
ensure that you review the [Multi-Region Private Endpoints](https://www.mongodb.com/docs/atlas/troubleshoot-private-endpoints/#multi-region-private-endpoints) troubleshooting documentation.

## Example Usage

```terraform
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS"
  region        = "US_EAST_1"
}

data "mongodbatlas_privatelink_endpoint" "test" {
	project_id      = mongodbatlas_privatelink_endpoint.test.project_id
	private_link_id = mongodbatlas_privatelink_endpoint.test.private_link_id
    provider_name = "AWS"
}
```

### Available complete examples
- [Setup private connection to a MongoDB Atlas Cluster with AWS VPC](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.2.0/examples/mongodbatlas_privatelink_endpoint/aws/cluster)

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the private endpoint service that you want to retrieve.
* `provider_name` - (Required) Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts `AWS`, `AZURE` or `GCP`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
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
* `endpoint_group_names` - For GCP legacy architecture (when `port_mapping_enabled` is not set or `false` on the endpoint resource): A list of the endpoint group names associated with the private endpoint service. For GCP port-based architecture (when `port_mapping_enabled = true` on the endpoint resource): A list of private endpoint names associated with the private endpoint service.
* `region_name` - GCP region for the Private Service Connect endpoint service.
* `service_attachment_names` - For GCP legacy architecture (when `port_mapping_enabled` is not set or `false` on the endpoint resource): A list of service attachments connected to the private endpoint service (one per Atlas node). For GCP port-based architecture (when `port_mapping_enabled = true` on the endpoint resource): A list of one service attachment connected to the private endpoint service.
* `port_mapping_enabled` - Flag that indicates whether this endpoint service uses GCP port-mapping. When `true`, the endpoint service uses the new GCP port-based architecture (requires 1 endpoint). When `false`, it uses the legacy architecture. Only applicable for GCP provider.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/) Documentation for more information.