---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: network_container"
sidebar_current: "docs-mongodbatlas-resource-network-container"
description: |-
    Provides a Network Peering resource.
---

# mongodbatlas_network_container

`mongodbatlas_network_container` provides a Network Peering Container resource. The resource lets you create, edit and delete network peering containers. The resource requires your Project ID.

~> **IMPORTANT:**
<br> This resource creates one Network Peering container into which Atlas can deploy Network Peering connections. You must have either the Project Owner or Organization Owner role to successfully call this endpoint.
<br>
<br> The following list outlines the maximum number of Network Peering containers per cloud provider:
<br> &#8226; GCP - One container per project.
<br> &#8226; AWS and Azure - One container per cloud provider region per project.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Example with AWS.

```hcl
  resource "mongodbatlas_network_container" "test" {
    project_id       = "<YOUR-PROJECT-ID>"
    atlas_cidr_block = "10.8.0.0/21"
    provider_name    = "AWS"
    region_name      = "US_EAST_1"
  }

```

### Example with GCP

```hcl
resource "mongodbatlas_network_container" "test" {
  project_id       = "<YOUR-PROJECT-ID>"
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "GCP"
}
```

### Example with Azure

```hcl
resource "mongodbatlas_network_container" "test" {
  project_id       = "<YOUR-PROJECT-ID>"
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AZURE"
  region           = "US_EAST_2"
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `atlas_cidr_block` - (Required) CIDR block that Atlas uses for your clusters. Atlas uses the specified CIDR block for all other Network Peering connections created in the project. The Atlas CIDR block must be at least a /24 and at most a /21 in one of the following [private networks](https://tools.ietf.org/html/rfc1918.html#section-3).
* `provider_name`  - (Optional) Cloud provider for this Network Peering connection. If omitted, Atlas sets this parameter to AWS.
* `region_name` - (Optional | AWS provider only) The Atlas AWS region name for where this container will exist.
* `region` - (Optional | AZURE provider only) The Atlas Azure region name for where this container will exist.


## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `container_id` - The Network Peering Container ID.
* `id` -	The Terraform's unique identifier used internally for state management.
* `region_name` - The Atlas AWS region name for where this container exists.
* `region` - The Atlas Azure region name for where this container exists.
* `azure_subscription_id` - Unique identifer of the Azure subscription in which the VNet resides.
* `provisioned` - Indicates whether the project has Network Peering connections deployed in the container.
* `gcp_project_id` - Unique identifier of the GCP project in which the Network Peering connection resides.
* `network_name` - Name of the Network Peering connection in the Atlas project.
* `vpc_id` - Unique identifier of the project’s VPC.
* `vnet_name` - 	The name of the Azure VNet. This value is null until you provision an Azure VNet in the container.


## Import

Clusters can be imported using project ID and network peering container id, in the format `PROJECTID-CONTAINER-ID`, e.g.

```
$ terraform import mongodbatlas_network_container.my_container 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a
```

See detailed information for arguments and attributes: [MongoDB API Network Peering Container](https://docs.atlas.mongodb.com/reference/api/vpc-create-container/)
