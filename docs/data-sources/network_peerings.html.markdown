# Data Source: mongodbatlas_network_peerings

`mongodbatlas_network_peerings` describes all Network Peering Connections.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage

### Basic Example (AWS).

```terraform
resource "mongodbatlas_network_peering" "test" {
	accepter_region_name	= "us-east-1"	
	project_id    			= "<YOUR-PROJEC-ID>"
	container_id            = "507f1f77bcf86cd799439011"
	provider_name           = "AWS"
	route_table_cidr_block  = "192.168.0.0/24"
	vpc_id					= "vpc-abc123abc123"
	aws_account_id			= "abc123abc123"
}


data "mongodbatlas_network_peerings" "test" {
    project_id = mongodbatlas_network_peering.test.project_id
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the project to create the database user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Network Peering Connection ID.
* `results` - A list where each represents a Network Peering Connection.

### Network Peering Connection

* `peering_id` - Atlas assigned unique ID for the peering connection.
* `connection_id` - Unique identifier for the peering connection.
* `accepter_region_name` - Specifies the region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - Account ID of the owner of the peer VPC.
* `provider_name` - Cloud provider for this VPC peering connection. If omitted, Atlas sets this parameter to AWS. (Possible Values `AWS`, `AZURE`, `GCP`).
* `route_table_cidr_block` - Peer VPC CIDR block or subnet.
* `vpc_id` - Unique identifier of the peer VPC.
* `error_state_name` - Error state, if any. The VPC peering connection error state value can be one of the following: `REJECTED`, `EXPIRED`, `INVALID_ARGUMENT`.
* `status_name` - The VPC peering connection status value can be one of the following: `INITIATING`, `PENDING_ACCEPTANCE`, `FAILED`, `FINALIZING`, `AVAILABLE`, `TERMINATING`.
* `azure_directory_id` - Unique identifier for an Azure AD directory.
* `azure_subscription_id` - Unique identifer of the Azure subscription in which the VNet resides.
* `resource_group_name` - Name of your Azure resource group. 
* `vnet_name` - Name of your Azure VNet.
* `error_state` - Description of the Atlas error when `status` is `Failed`, Otherwise, Atlas returns `null`.
* `status` - Status of the Atlas network peering connection: `ADDING_PEER`, `AVAILABLE`, `FAILED`, `DELETING`, `WAITING_FOR_USER`.
* `gcp_project_id` - GCP project ID of the owner of the network peer. 
* `network_name` - Name of the network peer to which Atlas connects.
* `error_message` - When `"status" : "FAILED"`, Atlas provides a description of the error.

See detailed information for arguments and attributes: [MongoDB API Network Peering Connection](https://docs.atlas.mongodb.com/reference/api/vpc-get-connections-list/)
