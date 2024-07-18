# Resource: mongodbatlas_network_peering

`mongodbatlas_network_peering` provides a Network Peering Connection resource. The resource lets you create, edit and delete network peering connections. The resource requires your Project ID.  

Ensure you have first created a network container if it is required for your configuration.  See the network_container resource documentation to determine if you need a network container first.  Examples for creating both container and peering resource are shown below as well as examples for creating the peering connection only.

~> **GCP AND AZURE ONLY:** Connect via Peering Only mode is deprecated, so no longer needed.  See [disable Peering Only mode](https://docs.atlas.mongodb.com/reference/faq/connection-changes/#disable-peering-mode) for details

~> **AZURE ONLY:** To create the peering request with an Azure VNET, you must grant Atlas the following permissions on the virtual network.
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/read
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/write
    Microsoft.Network/virtualNetworks/virtualNetworkPeerings/delete
    Microsoft.Network/virtualNetworks/peer/action
For more information see https://docs.atlas.mongodb.com/security-vpc-peering/ and https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/

-> **Create a Whitelist:** Ensure you whitelist the private IP ranges of the subnets in which your application is hosted in order to connect to your Atlas cluster.  See the project_ip_whitelist resource.

-> **NOTE:** Groups and projects are synonymous terms. You may find **group_id** in the official documentation.


## Example Usage - Container & Peering Connection

### Global configuration for the following examples
```terraform
locals {
  project_id        = <your-project-id>

  # needed for GCP only
  GCP_PROJECT_ID = <your-google-project-id>

  # needed for Azure Only
  AZURE_DIRECTORY_ID = <your-azure-directory-id>
  AZURE_SUBSCRIPTION_ID = <Unique identifer of the Azure subscription in which the VNet resides>
  AZURE_RESOURCES_GROUP_NAME = <Name of your Azure resource group>
  AZURE_VNET_NAME = <Name of your Azure VNet>
}
```

### Example with AWS

```terraform
# Container example provided but not always required, 
# see network_container documentation for details. 
resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "AWS"
  region_name      = "US_EAST_1"
}

# Create the peering connection request
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
# Accept the peering connection request
resource "aws_vpc_peering_connection_accepter" "peer" {
  vpc_peering_connection_id = mongodbatlas_network_peering.test.connection_id
  auto_accept = true
}

```

### Example with GCP

```terraform

# Container example provided but not always required, 
# see network_container documentation for details. 
resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "10.8.0.0/21"
  provider_name    = "GCP"
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id     = local.project_id
  container_id   = mongodbatlas_network_container.test.container_id
  provider_name  = "GCP"
  gcp_project_id = local.GCP_PROJECT_ID
  network_name   = "default"
}

# the following assumes a GCP provider is configured
data "google_compute_network" "default" {
  name = "default"
}

# Create the GCP peer
resource "google_compute_network_peering" "peering" {
  name         = "peering-gcp-terraform-test"
  network      = data.google_compute_network.default.self_link
  peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
}

# Create the cluster once the peering connection is completed
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = local.project_id
  name           = "terraform-manually-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_4"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }

  depends_on = [ google_compute_network_peering.peering ]
}

#  Private connection strings are not available w/ GCP until the reciprocal
#  connection changes to available (i.e. when the status attribute changes
#  to AVAILABLE on the 'mongodbatlas_network_peering' resource, which
#  happens when the google_compute_network_peering and and
#  mongodbatlas_network_peering make a reciprocal connection).  Hence
#  since the cluster can be created before this connection completes
#  you may need to run `terraform refresh` to obtain the private connection strings.

```

### Example with Azure

```terraform

# Ensure you have created the required Azure service principal first, see
# see https://docs.atlas.mongodb.com/security-vpc-peering/

# Container example provided but not always required, 
# see network_container documentation for details. 
resource "mongodbatlas_network_container" "test" {
  project_id       = local.project_id
  atlas_cidr_block = local.ATLAS_CIDR_BLOCK
  provider_name    = "AZURE"
  region           = "US_EAST_2"
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id            = local.project_id
  container_id          = mongodbatlas_network_container.test.container_id
  provider_name         = "AZURE"
  azure_directory_id    = local.AZURE_DIRECTORY_ID
  azure_subscription_id = local.AZURE_SUBSCRIPTION_ID
  resource_group_name   = local.AZURE_RESOURCES_GROUP_NAME
  vnet_name             = local.AZURE_VNET_NAME
}

# Create the cluster once the peering connection is completed
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = local.project_id
  name           = "terraform-manually-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AZURE"
      region_name   = "US_EAST_2"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }

  depends_on = [ mongodbatlas_network_peering.test ]
}
```

## Example Usage - Peering Connection Only, Container Exists
You can create a peering connection if an appropriate container for your cloud provider already exists in your project (see the network_container resource for more information).  A container may already exist if you have already created a cluster in your project, if so you may obtain the `container_id` from the cluster resource as shown in the examples below.

### Example with AWS
```terraform
# Create an Atlas cluster, this creates a container if one
# does not yet exist for this AWS region
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = local.project_id
  name           = "terraform-manually-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

# the following assumes an AWS provider is configured
resource "aws_default_vpc" "default" {
  tags = {
    Name = "Default VPC"
  }
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "mongo_peer" {
  accepter_region_name   = "us-east-2"
  project_id             = local.project_id
  container_id           = one(values(mongodbatlas_advanced_cluster.test.container_id))
  provider_name          = "AWS"
  route_table_cidr_block = "172.31.0.0/16"
  vpc_id                 = aws_default_vpc.default.id
  aws_account_id         = local.AWS_ACCOUNT_ID
}

# Accept the connection 
resource "aws_vpc_peering_connection_accepter" "aws_peer" {
  vpc_peering_connection_id = mongodbatlas_network_peering.mongo_peer.connection_id
  auto_accept               = true

  tags = {
    Side = "Accepter"
  }
}
```

### Example with GCP
```terraform
# Create an Atlas cluster, this creates a container if one
# does not yet exist for this GCP 
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = local.project_id
  name           = "terraform-manually-test"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_2"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id       = local.project_id
  atlas_cidr_block = "192.168.0.0/18"

  container_id   = one(values(mongodbatlas_advanced_cluster.test.replication_specs[0].container_id))
  provider_name  = "GCP"
  gcp_project_id = local.GCP_PROJECT_ID
  network_name   = "default"
}

# the following assumes a GCP provider is configured
data "google_compute_network" "default" {
  name = "default"
}

# Create the GCP peer
resource "google_compute_network_peering" "peering" {
  name         = "peering-gcp-terraform-test"
  network      = data.google_compute_network.default.self_link
  peer_network = "https://www.googleapis.com/compute/v1/projects/${mongodbatlas_network_peering.test.atlas_gcp_project_id}/global/networks/${mongodbatlas_network_peering.test.atlas_vpc_name}"
}
```

### Example with Azure

```terraform

# Ensure you have created the required Azure service principal first, see
# see https://docs.atlas.mongodb.com/security-vpc-peering/

# Create an Atlas cluster, this creates a container if one
# does not yet exist for this AZURE region
resource "mongodbatlas_advanced_cluster" "test" {
  project_id     = local.project_id
  name           = "cluster-azure"
  cluster_type   = "REPLICASET"
  backup_enabled = true

  replication_specs {
    region_configs {
      priority      = 7
      provider_name = "AZURE"
      region_name   = "US_EAST_2"
      electable_specs {
        instance_size = "M10"
        node_count    = 3
      }
    }
  }
}

# Create the peering connection request
resource "mongodbatlas_network_peering" "test" {
  project_id            = local.project_id
  container_id          = one(values(mongodbatlas_advanced_cluster.test.replication_specs[0].container_id))
  provider_name         = "AZURE"
  azure_directory_id    = local.AZURE_DIRECTORY_ID
  azure_subscription_id = local.AZURE_SUBSCRIPTION_ID
  resource_group_name   = local.AZURE_RESOURCE_GROUP_NAME
  vnet_name             = local.AZURE_VNET_NAME
}
```

## Argument Reference

* `project_id` - (Required) The unique ID for the MongoDB Atlas project to create the database user.
* `container_id` - (Required) Unique identifier of the MongoDB Atlas container for the provider (GCP) or provider/region (AWS, AZURE). You can create an MongoDB Atlas container using the network_container resource or it can be obtained from the cluster returned values if a cluster has been created before the first container.
* `provider_name` - (Required) Cloud provider to whom the peering connection is being made. (Possible Values `AWS`, `AZURE`, `GCP`).

**AWS ONLY:**

* `accepter_region_name` - (Required - AWS) Specifies the AWS region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - (Required - AWS) AWS Account ID of the owner of the peer VPC.
* `vpc_id` - (Required) Unique identifier of the AWS peer VPC (Note: this is **not** the same as the Atlas AWS VPC that is returned by the network_container resource).
* `route_table_cidr_block` - (Required - AWS) AWS VPC CIDR block or subnet.

**GCP ONLY:**

* `gcp_project_id` - (Required - GCP) GCP project ID of the owner of the network peer.
* `network_name` - (Required - GCP) Name of the network peer to which Atlas connects.
  
**AZURE ONLY:** 

* `azure_directory_id` - (Required - AZURE) Unique identifier for an Azure AD directory.
* `azure_subscription_id` - (Required - AZURE) Unique identifier of the Azure subscription in which the VNet resides.
* `resource_group_name` - (Required - AZURE) Name of your Azure resource group.
* `vnet_name` - (Required - AZURE) Name of your Azure VNet.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `peer_id` - Unique identifier of the Atlas network peer.
* `id` - Terraform's unique identifier used internally for state management.
* `connection_id` -  Unique identifier of the Atlas network peering container.
* `provider_name` - Cloud provider to whom the peering connection is being made. (Possible Values `AWS`, `AZURE`, `GCP`).

**AWS ONLY:**

* `accepter_region_name` - Specifies the region where the peer VPC resides. For complete lists of supported regions, see [Amazon Web Services](https://docs.atlas.mongodb.com/reference/amazon-aws/).
* `aws_account_id` - Account ID of the owner of the peer VPC.
* `route_table_cidr_block` - Peer VPC CIDR block or subnet.
* `vpc_id` - Unique identifier of the peer VPC.
* `error_state_name` - Error state, if any. The VPC peering connection error state value can be one of the following: `REJECTED`, `EXPIRED`, `INVALID_ARGUMENT`.
* `status_name` - (AWS Only) The VPC peering connection status value can be one of the following: `INITIATING`, `PENDING_ACCEPTANCE`, `FAILED`, `FINALIZING`, `AVAILABLE`, `TERMINATING`.

**AZURE/GCP ONLY:**

* `status` - Status of the Atlas network peering connection.  Azure/GCP: `ADDING_PEER`, `AVAILABLE`, `FAILED`, `DELETING` GCP Only:  `WAITING_FOR_USER`.
  
**GCP ONLY:**

* `gcp_project_id` - GCP project ID of the owner of the network peer.
* `error_message` - When `"status" : "FAILED"`, Atlas provides a description of the error.
* `network_name` - Name of the network peer to which Atlas connects.
* `atlas_gcp_project_id` - The Atlas GCP Project ID for the GCP VPC used by your atlas cluster that is needed to set up the reciprocal connection.
* `atlas_vpc_name` - Name of the GCP VPC used by your atlas cluster that is needed to set up the reciprocal connection.
  
**AZURE ONLY:**

* `azure_directory_id` - Unique identifier for an Azure AD directory.
* `azure_subscription_id` - Unique identifer of the Azure subscription in which the VNet resides.
* `error_state` - Description of the Atlas error when `status` is `Failed`, Otherwise, Atlas returns `null`.
* `resource_group_name` - Name of your Azure resource group.
* `vnet_name` - Name of your Azure VNet.


## Import

Network Peering Connections can be imported using project ID and network peering id, in the format `PROJECTID-PEERID-PROVIDERNAME`, e.g.

```
$ terraform import mongodbatlas_network_peering.my_peering 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a-AWS
```

See detailed information for arguments and attributes: [MongoDB API Network Peering Connection](https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/)

-> **NOTE:** If you need to get an existing container ID see the [How-To Guide](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/howto-guide.html).
