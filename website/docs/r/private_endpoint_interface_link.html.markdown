---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint_link"
sidebar_current: "docs-mongodbatlas-resource-private_endpoint_link"
description: |-
    Provides a Private Endpoint Link resource.
---

# mongodbatlas_private_endpoint_link

`mongodbatlas_private_endpoint_link` provides a Private Endpoint Link resource. This represents a Private Endpoint Link, which adds one interface endpoint to a private endpoint connection in an Atlas project.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

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

resource "mongodbatlas_private_endpoint_link" "test" {
  project_id            = "${mongodbatlas_private_endpoint.test.project_id}"
  private_link_id       = "${mongodbatlas_private_endpoint.test.private_link_id}"
  interface_endpoint_id = "${aws_vpc_endpoint.ptfe_service.id}"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the AWS PrivateLink connection which is created by `mongodbatlas_private_endpoint` resource.
* `interface_endpoint_id` - (Required) Unique identifier of the interface endpoint you created in your VPC with the AWS resource.

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

## Import
Private Endpoint Link Connection can be imported using project ID and username, in the format `{project_id}-{private_link_id}-{interface_endpoint_id}`, e.g.

```
$ terraform import mongodbatlas_private_endpoint_link.test 1112222b3bf99403840e8934-3242342343112-vpce-4242342343
```

See detailed information for arguments and attributes: [MongoDB API Private Endpoint Link Connection](https://docs.atlas.mongodb.com/reference/api/private-endpoint-create-one-interface-endpoint/)