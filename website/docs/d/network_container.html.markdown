---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: network_container"
sidebar_current: "docs-mongodbatlas-datasource-network-container"
description: |-
    Describes a Cluster resource.
---

# Data Source: mongodbatlas_network_container

`mongodbatlas_network_container` describes a Network Peering Container. The resource requires your Project ID and container ID.

~> **IMPORTANT:** This resource creates one Network Peering container into which Atlas can deploy Network Peering connections. An Atlas project can have a maximum of one container for each cloud provider. You must have either the Project Owner or Organization Owner role to successfully call this endpoint.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Basic Example.

```terraform
resource "mongodbatlas_network_container" "test" {
  project_id       = "<YOUR-PROJECT-ID>"
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AWS"
  region_name      = "US_EAST_1"
}

data "mongodbatlas_network_container" "test" {
	project_id   		= mongodbatlas_network_container.test.project_id
	container_id		= mongodbatlas_network_container.test.id
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `container_id` - (Required) The Network Peering Container ID.



## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Network Peering Container ID.
* `atlas_cidr_block` - CIDR block that Atlas uses for your clusters. Atlas uses the specified CIDR block for all other Network Peering connections created in the project. The Atlas CIDR block must be at least a /24 and at most a /21 in one of the following [private networks](https://tools.ietf.org/html/rfc1918.html#section-3).
* `provider_name`  - Cloud provider for this Network Peering connection. If omitted, Atlas sets this parameter to AWS.
* `region_name` - The Atlas AWS region name for where this container will exist.
* `region` - The Atlas Azure region name for where this container will exist.
* `azure_subscription_id` - Unique identifer of the Azure subscription in which the VNet resides.
* `provisioned` - Indicates whether the project has Network Peering connections deployed in the container.
* `gcp_project_id` - Unique identifier of the GCP project in which the Network Peering connection resides.
* `network_name` - Name of the Network Peering connection in the Atlas project.
* `vpc_id` - Unique identifier of the project’s VPC.
* `vnet_name` - 	The name of the Azure VNet. This value is null until you provision an Azure VNet in the container.
* `regions` - Atlas GCP regions where the container resides.


See detailed information for arguments and attributes: [MongoDB API Network Peering Container](https://docs.atlas.mongodb.com/reference/api/vpc-create-container/)

-> **NOTE:** If you need to get an existing container ID see the [How-To Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/howto-guide.html).