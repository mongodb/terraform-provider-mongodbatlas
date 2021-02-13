---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: private_endpoint"
sidebar_current: "docs-mongodbatlas-resource-private_endpoint"
description: |-
    Provides a Private Endpoint resource.
---

# mongodbatlas_privatelink_endpoint

`mongodbatlas_privatelink_endpoint` provides a Private Endpoint resource. This represents a Private Endpoint Service that can be created in an Atlas project.

~> **IMPORTANT:**You must have one of the following roles to successfully handle the resource:
  * Organization Owner
  * Project Owner

-> **NOTE:** Groups and projects are synonymous terms. You may find group_id in the official documentation.

-> **NOTE:** A network container is created for a private endpoint to reside in if one does not yet exist in the project.  


## Example Usage

```hcl
resource "mongodbatlas_privatelink_endpoint" "test" {
  project_id    = "<PROJECT-ID>"
  provider_name = "AWS/AZURE"
  region        = "us-east-1"
}
```

## Argument Reference

* `project_id` - Required 	Unique identifier for the project.
* `providerName` - (Required) Name of the cloud provider for which you want to create the private endpoint service. Atlas accepts `AWS` or `AZURE`.
* `region` - (Required) Cloud provider region in which you want to create the private endpoint connection.
Accepted values are: [AWS regions](https://docs.atlas.mongodb.com/reference/amazon-aws/#amazon-aws) and [AZURE regions](https://docs.atlas.mongodb.com/reference/microsoft-azure/#microsoft-azure)
  
  AWS:
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

AZURE:
  * `centralus`
  * `eastus`
  * `eastus2`
  * `northcentralus`
  * `westus`
  * `southcentralus`
  * `westus2`
  * `westcentralus`
  * `brazilsouth`
  * `canadaeast`
  * `canadacentral`
  * `northeurope`
  * `westeurope`
  * `uksouth`
  * `ukwest`
  * `francecentral`
  * `germanywestcentral`
  * `germanynorth`
  * `switzerlandnorth`
  * `switzerlandwest`
  * `norwayeast`
  * `eastasia`
  * `southeastasia`
  * `australiaeast`
  * `australiasoutheast`
  * `centralindia`
  * `southindia`
  * `westindia`
  * `japaneast`
  * `japanwest`
  * `koreacentral`
  * `koreasouth`
  * `southafricanorth`
  * `uaenorth`
  * `uaecentral`


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Terraform's unique identifier used internally for state management.
* `private_link_id` - Unique identifier of the AWS PrivateLink connection.
* `endpoint_service_name` - Name of the PrivateLink endpoint service in AWS. Returns null while the endpoint service is being created.
* `error_message` - Error message pertaining to the AWS PrivateLink connection. Returns null if there are no errors.
* `interface_endpoints` - Unique identifiers of the interface endpoints in your VPC that you added to the AWS PrivateLink connection.
* `private_endpoints` - All private endpoints that you have added to this Azure Private Link Service.
* `private_link_service_name` - Name of the Azure Private Link Service that Atlas manages.
* `private_link_service_resource_id` - Resource ID of the Azure Private Link Service that Atlas manages.
* `status` - Status of the AWS PrivateLink connection or Status of the Azure Private Link Service. Atlas returns one of the following values:
  AWS:
    * `AVAILABLE` 	Atlas is creating the network load balancer and VPC endpoint service.
    * `WAITING_FOR_USER` The Atlas network load balancer and VPC endpoint service are created and ready to receive connection requests. When you receive this status, create an interface endpoint to continue configuring the AWS PrivateLink connection.
    * `FAILED` 	A system failure has occurred.
    * `DELETING` 	The AWS PrivateLink connection is being deleted.
  AZURE:
    * `AVAILABLE` 	Atlas created the load balancer and the Private Link Service.
    * `INITIATING` 	Atlas is creating the load balancer and the Private Link Service.
    * `FAILED` 	Atlas failed to create the load balancer and the Private Link service.
    * `DELETING` 	Atlas is deleting the Private Link service.

## Import
Private Endpoint Service can be imported using project ID, private link ID, provider name and region, in the format `{project_id}-{private_link_id}-{provider_name}-{region}`, e.g.

```
$ terraform import mongodbatlas_privatelink_endpoint.test 1112222b3bf99403840e8934-3242342343112-AWS-us-east-1
```

See detailed information for arguments and attributes: [MongoDB API Private Endpoint Service](https://docs.atlas.mongodb.com/reference/api/private-endpoints-service-create-one//)