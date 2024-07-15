# Data Source: mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` describes a Private Endpoint. This represents a Private Endpoint Connection to retrieve details regarding a private endpoint by id in an Atlas project

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

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
- [Setup private connection to a MongoDB Atlas Cluster with AWS VPC](https://github.com/mongodb/terraform-provider-mongodbatlas/blob/master/examples/aws-privatelink-endpoint/cluster)

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
* `endpoint_group_names` - GCP network endpoint groups corresponding to the Private Service Connect endpoint service.
* `region_name` - GCP region for the Private Service Connect endpoint service.
* `service_attachment_names` - Unique alphanumeric and special character strings that identify the service attachments associated with the GCP Private Service Connect endpoint service.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/) Documentation for more information.