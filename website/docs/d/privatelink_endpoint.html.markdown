---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint"
sidebar_current: "docs-mongodbatlas-datasource-private-endpoint"
description: |-
    Describes a Private Endpoint.
---

# mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` describe a Private Endpoint. This represents a Private Endpoint Connection to retrieve details regarding a private endpoint by id in an Atlas project

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS"
  region        = "us-east-1"
}

data "mongodbatlas_privatelink_endpoint" "test" {
	project_id      = "${mongodbatlas_privatelink_endpoint.test.project_id}"
	private_link_id = "${mongodbatlas_privatelink_endpoint.test.private_link_id}"
    provider_name = "AWS"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the private endpoint service that you want to retrieve.
* `provider_name` - (Required) Cloud provider for which you want to retrieve a private endpoint service. Atlas accepts `AWS` or `AZURE`.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `status` - Status of the AWS PrivateLink connection.
* `private_endpoints` - All private endpoints that you have added to this Azure Private Link Service.
* `private_link_service_name` - Name of the Azure Private Link Service that Atlas manages.
* `private_link_service_resource_id` - Resource ID of the Azure Private Link Service that Atlas manages.
  Returns one of the following values:
  * `AVAILABLE` 	Atlas created the load balancer and the Private Link Service.
  * `INITIATING` 	Atlas is creating the network load balancer and VPC endpoint service.
  * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
  * `FAILED` 	A system failure has occurred.
  * `DELETING` 	The Private Link service is being deleted.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/) Documentation for more information.