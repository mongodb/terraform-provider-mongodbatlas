---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint_link"
sidebar_current: "docs-mongodbatlas-datasource-private-endpoint-link"
description: |-
    Describes a Private Endpoint Link.
---

# mongodbatlas_privatelink_endpoint_service

`mongodbatlas_privatelink_endpoint_service` describe a Private Endpoint Link. This represents a Private Endpoint Link Connection that wants to retrieve details in an Atlas project.

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

## Example Usage

```hcl
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT_ID>"
  provider_name = "AWS"
  region        = "us-east-1"
}

resource "aws_vpc_endpoint" "ptfe_service" {
  vpc_id             = "vpc-7fc0a543"
  service_name       = "${mongodbatlas_privatelink_endpoint.test.endpoint_service_name}"
  vpc_endpoint_type  = "Interface"
  subnet_ids         = ["subnet-de0406d2"]
  security_group_ids = ["sg-3f238186"]
}

resource "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id            = "${mongodbatlas_privatelink_endpoint.test.project_id}"
  private_link_id       = "${mongodbatlas_privatelink_endpoint.test.private_link_id}"
  endpoint_service_id   = "${aws_vpc_endpoint.ptfe_service.id}"
  provider_name         ="AWS"
}

data "mongodbatlas_privatelink_endpoint_service" "test" {
  project_id            = "${mongodbatlas_privatelink_endpoint_service.test.project_id}"
  private_link_id       = "${mongodbatlas_privatelink_endpoint_service.test.private_link_id}"
  interface_endpoint_id = "${mongodbatlas_privatelink_endpoint_service.test.interface_endpoint_id}"
}
```

## Argument Reference

* `project_id` - (Required) Unique identifier for the project.
* `private_link_id` - (Required) Unique identifier of the private endpoint service for which you want to retrieve a private endpoint.
* `endpoint_service_id` - (Required) Unique identifier of the `AWS` or `AZURE` resource.
* `provider_name` - (Required) Cloud provider for which you want to create a private endpoint. Atlas accepts `AWS` or `AZURE`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `interface_endpoint_id` - Unique identifier of the interface endpoint.
* `private_endpoint_connection_name` - Name of the connection for this private endpoint that Atlas generates.
* `private_endpoint_ip_address` - Private IP address of the private endpoint network interface.
* `private_endpoint_resource_id` - Unique identifier of the private endpoint.
* `delete_requested` - Indicates if Atlas received a request to remove the interface endpoint from the private endpoint connection.
* `error_message` - Error message pertaining to the interface endpoint. Returns null if there are no errors.
* `connection_status` - Status of the interface endpoint.
  Returns one of the following values:
    * `NONE` - Atlas created the network load balancer and VPC endpoint service, but AWS hasn’t yet created the VPC endpoint.
    * `PENDING_ACCEPTANCE` - AWS has received the connection request from your VPC endpoint to the Atlas VPC endpoint service.
    * `PENDING` - AWS is establishing the connection between your VPC endpoint and the Atlas VPC endpoint service.
    * `AVAILABLE` - Atlas VPC resources are connected to the VPC endpoint in your VPC. You can connect to Atlas clusters in this region using AWS PrivateLink.
    * `REJECTED` - AWS failed to establish a connection between Atlas VPC resources to the VPC endpoint in your VPC.
    * `INITIATING` - Atlas has not yet accepted the connection to your private endpoint.
    * `FAILED` - Atlas failed to accept the connection your private endpoint.
    * `DELETING` - Atlas is removing the interface endpoint from the private endpoint connection.

See [MongoDB Atlas API](https://docs.atlas.mongodb.com/reference/api/private-endpoints-endpoint-get-one/) Documentation for more information.
