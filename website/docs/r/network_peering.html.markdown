---
layout: "mongodbatlas"
page_title: "MongoDB Atlas: network_peering"
sidebar_current: "docs-mongodbatlas-resource-network-peering"
description: |-
    Provides a Network Peering resource.
---

# mongodb_atlas_network_peering

`mongodb_atlas_network_peering` provides a Network Peering Connection resource. The resource lets you create, edit and delete network peering connections. The resource requires your Project ID.


~> **GCP AND AZURE ONLY:** You must enable Connect via Peering Only mode to use network peering.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Example with AWS.

```hcl
	resource "mongodbatlas_network_peering" "test" {
		accepter_region_name	  = "us-east-1"	
		project_id    		     	= "<YOUR-PROJEC-ID>"
		container_id            = "507f1f77bcf86cd799439011"
		provider_name           = "AWS"
		route_table_cidr_block  = "192.168.0.0/24"
		vpc_id					        = "vpc-abc123abc123"
		aws_account_id		    	= "abc123abc123"
	}
```

### Example with GCP

```hcl
	resource "mongodbatlas_network_peering" "test" {	
		project_id    	  = "<YOUR-PROJEC-ID>"
		container_id      = "507f1f77bcf86cd799439011"
		provider_name     = "GCP"
			gcp_project_id  = "my-sample-project-191923"
			network_name    = "test1"	
	}
```

### Example with Azure

```hcl
	resource "mongodbatlas_network_peering" "test" {	
		project_id    			  = "<YOUR-PROJEC-ID>"
		atlas_cidr_block      = "192.168.0.0/21"
		container_id          = "507f1f77bcf86cd799439011"
		provider_name         = "AZURE"
		azure_directory_id    = "35039750-6ebd-4ad5-bcfe-cb4e5fc2d915"
		azure_subscription_id = "g893dec2-d92e-478d-bc50-cf99d31bgeg9"
		resource_group_name   = "atlas-azure-peering"
		vnet_name	            = "azure-peer"
	}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.
* `container_id` - (Required) Unique identifier of the Atlas VPC container for the region. You can create an Atlas VPC container using the Create Container endpoint. You cannot create more than one container per region. To retrieve a list of container IDs, use the Get list of VPC containers endpoint.
* `accepter_region_name` - (Optional | **AWS Required**) Specifies the region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - (Optional | **AWS Required**) Account ID of the owner of the peer VPC.
* `provider_name` - (Optional) Cloud provider for this VPC peering connection. If omitted, Atlas sets this parameter to AWS. (Possible Values `AWS`, `AZURE`, `GCP`).
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
* `network_name` - Name of the network peer to which Atlas connects.
* `error_message` - When `"status" : "FAILED"`, Atlas provides a description of the error.


## Import

Clusters can be imported using project ID and network peering peering id, in the format `PROJECTID-PEER-ID`, e.g.

```
$ terraform import mongodbatlas_network_peering.my_peering 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a
```

See detailed information for arguments and attributes: [MongoDB API Network Peering Connection](https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/)