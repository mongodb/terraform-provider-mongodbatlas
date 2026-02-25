---
subcategory: "Network Peering"
---

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

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_4"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]

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

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AZURE"
      region_name   = "US_EAST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]

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

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AWS"
      region_name   = "US_EAST_1"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
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

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "GCP"
      region_name   = "US_EAST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
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

  replication_specs = [{
    region_configs = [{
      priority      = 7
      provider_name = "AZURE"
      region_name   = "US_EAST_2"
      electable_specs = {
        instance_size = "M10"
        node_count    = 3
      }
    }]
  }]
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

### Further Examples
- [AWS Network Peering](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_network_peering/aws)
- [Azure Network Peering](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_network_peering/azure)
- [GCP Network Peering](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/v2.7.0/examples/mongodbatlas_network_peering/gcp)

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `container_id` (String)
- `project_id` (String)
- `provider_name` (String)

### Optional

- `accepter_region_name` (String)
- `atlas_cidr_block` (String)
- `atlas_gcp_project_id` (String)
- `atlas_vpc_name` (String)
- `aws_account_id` (String)
- `azure_directory_id` (String)
- `azure_subscription_id` (String)
- `delete_on_create_timeout` (Boolean) Indicates whether to delete the resource being created if a timeout is reached when waiting for completion. When set to `true` and timeout occurs, it triggers the deletion and returns immediately without waiting for deletion to complete. When set to `false`, the timeout will not trigger resource deletion. If you suspect a transient error when the value is `true`, wait before retrying to allow resource deletion to finish. Default is `true`.
- `gcp_project_id` (String)
- `network_name` (String)
- `resource_group_name` (String)
- `route_table_cidr_block` (String)
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `vnet_name` (String)
- `vpc_id` (String)

### Read-Only

- `atlas_id` (String)
- `connection_id` (String)
- `error_message` (String)
- `error_state` (String)
- `error_state_name` (String)
- `id` (String) The ID of this resource.
- `peer_id` (String)
- `status` (String)
- `status_name` (String)

<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)
- `delete` (String)
- `update` (String)

## Import

Network Peering Connections can be imported using project ID and network peering id, in the format `PROJECTID-PEERID-PROVIDERNAME`, e.g.

```
$ terraform import mongodbatlas_network_peering.my_peering 1112222b3bf99403840e8934-5cbf563d87d9d67253be590a-AWS
```

Use the [MongoDB Atlas CLI][https://www.mongodb.com/docs/atlas/cli/current/command/atlas-networking-peering-list/#std-label-atlas-networking-peering-list] to obtain your `project_id` and `peering_id`. Attention gcp and azure users: The `atlas networking peering list` command returns only `AWS` peerings by default. You have to include the `--provider` parameter to list peerings for your cloud provider. Valid values are AWS, AZURE, or GCP.

```
atlas projects list
atlas networking peering list --projectId <projectId> --provider <AZURE|GCP|AWS>
```
See detailed information for arguments and attributes: [MongoDB API Network Peering Connection](https://docs.atlas.mongodb.com/reference/api/vpc-create-peering-connection/)
