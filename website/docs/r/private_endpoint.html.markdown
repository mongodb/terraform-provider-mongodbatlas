---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint"
sidebar_current: "docs-mongodbatlas-resource-private_endpoint"
description: |-
    Provides a Private Endpoint resource.
---

# mongodbatlas_private_endpoint

`mongodbatlas_private_endpoint` provides a Private Endpoint resource. This represents a Private Endpoint Connection that can be created in an Atlas project.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.


## Example Usage

```hcl
resource "mongodbatlas_private_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS"
  region        = "us-east-1"
}
```

## Argument Reference

* `project_id` - Required 	Unique identifier for the project.
* `providerName` - (Required) Name of the cloud provider you want to create the private endpoint connection for. Must be AWS.
* `region` - (Required) Cloud provider region in which you want to create the private endpoint connection.
Accepted values are:
  * `us-east-1`
  * `us-east-2`
  * `us-west-1`
  * `us-west-2`
  * `ca-central-1`
  * `sa-east-1`
  * `eu-north-1`
  * `eu-west-1`
  * `eu-west-2`
  * `eu-west-3`
  * `eu-central-1`
  * `me-south-1`
  * `ap-northeast-1`
  * `ap-northeast-2`
  * `ap-south-1`
  * `ap-southeast-1`
  * `ap-southeast-2`
  * `ap-east-1`


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `private_link_id` - Unique identifier of the AWS PrivateLink connection.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `status` - Status of the AWS PrivateLink connection.
  Returns one of the following values:
  * `INITIATING` 	Atlas is creating the network load balancer and VPC endpoint service.
  * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
  * `FAILED` 	A system failure has occurred.
  * `DELETING` 	The AWS PrivateLink connection is being deleted.

## Import
Private Endpoint Connection can be imported using project ID and username, in the format `{project_id}-{private_link_id}`, e.g.

```
$ terraform import mongodbatlas_private_endpoint.test 1112222b3bf99403840e8934-3242342343112
```

See detailed information for arguments and attributes: [MongoDB API Private Endpoint Connection](https://docs.atlas.mongodb.com/reference/api/private-endpoint-create-one-private-endpoint-connection/)