---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint"
sidebar_current: "docs-mongodbatlas-datasource-private-endpoint"
description: |-
    Describes a Private Endpoint.
---

# mongodbatlas_private_endpoint

`mongodbatlas_private_endpoint` describe a Private Endpoint. This represents a Private Endpoint Connection to retrieve details regarding a private endpoint by id in an Atlas project

!> **WARNING:** This datasource is deprecated and will be removed in the next major version
                Please transition to privatelink_endpoint as soon as possible. [PrivateLink Endpoints](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-get-one/)

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_private_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS"
  region        = "us-east-1"
}

data "mongodbatlas_private_endpoint" "test" {
	project_id      = mongodbatlas_private_endpoint.test.project_id
	private_link_id = mongodbatlas_private_endpoint.test.private_link_id
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the AWS PrivateLink connection.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `status` - Status of the AWS PrivateLink connection or Status of the Azure Private Link Service. Atlas returns one of the following values:
  AWS:
    * `AVAILABLE` 	Atlas is creating the network load balancer and VPC endpoint service.
    * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests.
    * `FAILED` 	A system failure has occurred.
    * `DELETING` 	The AWS PrivateLink connection is being deleted.
  AZURE:
    * `AVAILABLE` 	Atlas created the load balancer and the Private Link Service.
    * `INITIATING` 	Atlas is creating the load balancer and the Private Link Service.
    * `FAILED` 	Atlas failed to create the load balancer and the Private Link service.
    * `DELETING` 	Atlas is deleting the Private Link service.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoint-get-one-private-endpoint-connection/) Documentation for more information.