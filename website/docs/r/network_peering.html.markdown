---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: network_peering"
sidebar_current: "docs-mongodbatlas-resource-network-peering"
description: |-
    Provides a Network Peering resource.
---

# mongodbatlas_network_peering

`mongodbatlas_network_peering` provides a Network Peering Connection resource. The resource lets you create, edit and delete network peering connections. The resource requires your Project ID.  Ensure you have first created a Network Container.  See the network_container resource and examples below.

~> **GCP AND AZURE ONLY:** You must enable Connect via Peering Only mode to use network peering.

~> **AZURE ONLY:** To create the peering request with an Azure VNET, you must grant Atlas the following permissions on the virtual network.
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/read
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/delete
    Microsoft.Network/virtualNetworks/peer/action
For more information see https://docs.atlas.mongodb.com/security-vpc-peering/

-> **Create a Whitelist:** Ensure you whitelist the private IP ranges of the subnets in which your application is hosted in order to connect to your Atlas cluster.  See the project_ip_whitelist resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Global configuration for the following examples
```hcl
locals {
  project_id        = <your-project-id>

  # needed for GCP only
  google_project_id = <your-google-project-id>
}
```

### Example with AWS.

```hcl
resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AWS"
  region_name      = "US_EAST_1"
}

resource "mongodbatlas_network_peering" "test" {
  accepter_region_name   = "us-east-1"
  project_id             = local.project_id
  container_id           = "507f1f77bcf86cd799439011"
  provider_name          = "AWS"
  route_table_cidr_block = "192.168.0.0/24"
  vpc_id                 = "vpc-abc123abc123"
  aws_account_id         = "abc123abc123"
}

# the following assumes an AWS provider is configured  
resource "aws_vpc_peering_connection_accepter" "peer" {
  vpc_peering_connection_id = "${mongodbatlas_network_peering.test.connection_id}"
  auto_accept = true
}

```

### Example with GCP

```hcl

resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "192.168.192.0/18"
  provider_name    = "GCP"
}

resource "mongodbatlas_private_ip_mode" "my_private_ip_mode" {
  project_id = local.project_id
  enabled    = true
}

resource "mongodbatlas_network_peering" "test" {
  project_id     = local.project_id
  container_id   = mongodbatlas_network_container.test.container_id
  provider_name  = "GCP"
  network_name   = "myNetWorkPeering"
  gcp_project_id = local.google_project_id

  depends_on = [mongodbatlas_private_ip_mode.my_private_ip_mode]
}

resource "google_compute_network" "vpc_network" {
  name = "vpcnetwork"
}

resource "google_compute_network_peering" "gcp_main_atlas_peering" {
  name         = "atlas-gcp-main"
  network      = google_compute_network.vpc_network.self_link
  peer_network = "projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
}
```

### Example with Azure

```hcl

resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AZURE"
  region           = "US_WEST"
}

resource "mongodbatlas_private_ip_mode" "my_private_ip_mode" {
  project_id = "${mongodbatlas_project.my_project.id}"
  enabled    = true
}

resource "mongodbatlas_network_peering" "test" {
  project_id            = local.project_id
  atlas_cidr_block      = "10.8.0.0/21"
  container_id          = mongodbatlas_network_container.test.container_id
  provider_name         = "AZURE"
  azure_directory_id    = "35039750-6ebd-4ad5-bcfe-cb4e5fc2d915"
  azure_subscription_id = "g893dec2-d92e-478d-bc50-cf99d31bgeg9"
  resource_group_name   = "atlas-azure-peering"
  vnet_name             = "azure-peer"

  depends_on = [mongodbatlas_private_ip_mode.my_private_ip_mode]
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `container_id` - (Required) Unique identifier of the Atlas VPC container for the region. You can create an Atlas VPC container using the Create Container endpoint. You cannot create more than one container per region. To retrieve a list of container IDs, use the Get list of VPC containers endpoint.
* `provider_name` - (Required) Cloud provider for this VPC peering connection. (Possible Values `AWS`, `AZURE`, `GCP`).
* `accepter_region_name` - (Optional | **AWS Required**) Specifies the region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - (Optional | **AWS Required**) Account ID of the owner of the peer VPC.
* `route_table_cidr_block` - (Optional | **AWS Required**) Peer VPC CIDR block or subnet.
* `vpc_id` - (Optional | **AWS Required**) Unique identifier of the peer VPC.
* `atlas_cidr_block` - (Optional | **AZURE Required**) Unique identifier for an Azure AD directory.
* `azure_directory_id` - (Optional | **AZURE Required**) Unique identifier for an Azure AD directory.
* `azure_subscription_id` - (Optional | **AZURE Required**) Unique identifer of the Azure subscription in which the VNet resides.
* `resource_group_name` - (Optional | **AZURE Required**) Name of your Azure resource group.
* `vnet_name` - (Optional | **AZURE Required**) Name of your Azure VNet.
* `gcp_project_id` - (Optinal | **GCP Required**) GCP project ID of the owner of the network peer.
* `network_name` - (Optional | **GCP Required**) Name of the network peer to which Atlas connects.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `peer_id` - The Network Peering Container ID.
* `id` -	The Terraform's unique identifier used internally for state management.
* `connection_id` -  Unique identifier for the peering connection.
* `accepter_region_name` - Specifies the region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - Account ID of the owner of the peer VPC.
* `provider_name` - Cloud provider for this VPC peering connection. If omitted, Atlas sets this parameter to AWS. (Possible Values `AWS`, `AZURE`, `GCP`).
* `route_table_cidr_block` - Peer VPC CIDR block or subnet.
* `vpc_id` - Unique identifier of the peer VPC.
* `error_state_name` - Error state, if any. The VPC peering connection error state value can be one of the following: `REJECTED`, `EXPIRED`, `INVALID_ARGUMENT`.
* `status_name` - The VPC peering connection status value can be one of the following: `INITIATING`, `PENDING_ACCEPTANCE`, `FAILED`, `FINALIZING`, `AVAILABLE`, `TERMINATING`.
* `atlas_cidr_block` - Unique identifier for an Azure AD directory.
* `azure_directory_id` - Unique identifier for an Azure AD directory.
* `azure_subscription_id` - Unique identifer of the Azure subscription in which the VNet resides.
* `resource_group_name` - Name of your Azure resource group.
* `vnet_name` - Name of your Azure VNet.
* `error_state` - Description of the Atlas error when `status` is `Failed`, Otherwise, Atlas returns `null`.
* `status` - Status of the Atlas network peering connection: `ADDING_PEER`, `AVAILABLE`, `FAILED`, `DELETING`, `WAITING_FOR_USER`.
* `gcp_project_id` - GCP project ID of the owner of the network peer.
* `atlas_gcp_project_id` - The Atlas GCP Project ID for the GCP VPC used by your atlas cluster that it is need to set up the reciprocal connection.
* `atlas_vpc_name` - The Atlas VPC Name is used by your atlas clister that it is need to set up the reciprocal connection.
* `network_name` - Name of the network peer to which Atlas connects.
* `error_message` - When `"status" : "FAILED"`, Atlas provides a description of the error.


## Import

Clusters can be imported using project ID and network peering peering id, in the format `PROJECTID-PEERID-PROVIDERNAME`, e.g.

```
$ terraform import mongodbatlas_network_peering.my_peering 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a-AWS
```

See detailed information for arguments and attributes: [MongoDB API Network Peering Connection](https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/)
