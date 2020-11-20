---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint_link"
sidebar_current: "docs-mongodbatlas-datasource-private-endpoint-link"
description: |-
    Describes a Private Endpoint Link.
---

# mongodbatlas_private_endpoint_link

`mongodbatlas_private_endpoint_interface_link` describe a Private Endpoint Link. This represents a Private Endpoint Link Connection that wants to retrieve details in an Atlas project.

!> **WARNING:** This datasource is deprecated and will be removed in the next major version
                Please transition to privatelink_endpoint_service as soon as possible. [PrivateLink Endpoint Service](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-get-one/)

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_private_endpoint" "test" {
  project_id    = "<PROJECT_ID>"
  provider_name = "AWS"
  region        = "us-east-1"
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = "vpc-7fc0a543"
  service_name       = "${mongodbatlas_private_endpoint.test.endpoint_service_name}"
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_private_endpoint_interface_link" "test" {
  project_id            = "${mongodbatlas_private_endpoint.test.project_id}"
  private_link_id       = "${mongodbatlas_private_endpoint.test.private_link_id}"
  interface_endpoint_id = "${aws_vpc_endpoint.ptfe_service.id}"
}

data "mongodbatlas_private_endpoint_interface_link" "test" {
  project_id            = "${mongodbatlas_private_endpoint_link.test.project_id}"
  private_link_id       = "${mongodbatlas_private_endpoint_link.test.private_link_id}"
  interface_endpoint_id = "${mongodbatlas_private_endpoint_link.test.interface_endpoint_id}"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the AWS PrivateLink connection.
* `interface_endpoints` - (Required) Unique identifiers of the interface endpoints in your VPC.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `delete_requested` - Indicates if Atlas received a request to remove the interface endpoint from the private endpoint connection.
* `error_message` - Error message pertaining to the interface endpoint. Returns null if there are no errors.
* `connection_status` - Status of the interface endpoint.
  Returns one of the following values:
    * `NONE` - Atlas created the network load balancer and VPC endpoint service, but AWS hasn’t yet created the VPC endpoint.
    * `PENDING_ACCEPTANCE` - AWS has received the connection request from your VPC endpoint to the Atlas VPC endpoint service.
    * `PENDING` - AWS is establishing the connection between your VPC endpoint and the Atlas VPC endpoint service.
    * `AVAILABLE` - Atlas VPC resources are connected to the VPC endpoint in your VPC. You can connect to Atlas clusters in this region using AWS PrivateLink.
    * `REJECTED` - AWS failed to establish a connection between Atlas VPC resources to the VPC endpoint in your VPC.
    * `DELETING` - Atlas is removing the interface endpoint from the private endpoint connection.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoint-get-one-interface-endpoint/) Documentation for more information.
