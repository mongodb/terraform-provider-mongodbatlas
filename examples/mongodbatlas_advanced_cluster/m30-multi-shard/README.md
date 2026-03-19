# MongoDB Atlas Provider — M30 Multi-Shard Cluster (2 Shards)

This example creates an Atlas project and a **two-shard M30 sharded cluster** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. This is the recommended configuration for workloads that have outgrown a single replica set and need horizontal write scaling or want to distribute data across multiple shards.

Both shards are configured symmetrically (same instance size and region), which simplifies capacity planning and ensures even data distribution with a standard hashing shard key.

## Topology

The cluster distributes data across two shards as follows:

- Shard 1 — `US_EAST_1` (AWS), 3 electable nodes, instance size M30
- Shard 2 — `US_EAST_1` (AWS), 3 electable nodes, instance size M30

 **Tip:** To add a third shard, append another `replication_specs` block with the same shape as shards 1 and 2. Atlas performs a live online shard addition — no downtime is required.

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/)

## Procedure

Follow the next steps to run this example:

1. Set Up your MongoDB Atlas Credentials

   Create a `terraform.tfvars` file with your credentials:

   ```hcl
   atlas_org_id        = "<MONGODB_ATLAS_ORG_ID>"
   atlas_client_id     = "<ATLAS_CLIENT_ID>"
   atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
   ```

   Service Accounts are the recommended way to authenticate with the MongoDB Atlas API. To learn how to create a Service Account, see [Manage Service Accounts](https://www.mongodb.com/docs/atlas/configure-api-access/#manage-service-accounts) in the Atlas documentation.

2. Review your Terraform plan

   The following command lists the resources that your configuration will create.

   ```bash
   terraform plan
   ```

   Review the following fields in the plan output before applying:

   - `name`: The name of the cluster (default: `"m30-multi-shard"`)
   - `cluster_type`: The type of cluster. Must be `"SHARDED"`
   - `instance_size`: The instance size for all nodes. Set to `"M30"`
   - `node_count`: The number of electable nodes per shard. Set to `3`
   - `provider_name`: The cloud service provider hosting the cluster. Set to `"AWS"`
   - `region_name`: The cloud region where the cluster is deployed. Set to `"US_EAST_1"`
   - `backup_enabled`: Flags whether cloud backup is enabled. Set to `true`
   - `termination_protection_enabled`: Flags whether termination protection is enabled. If enabled, you cannot destroy the cluster using `terraform destroy`. Set to `false`
   - `replication_specs`: Two entries, one per shard, each with identical region configuration

3. Apply your configuration

   The following command applies your configuration and creates the resources.

   ```bash
   terraform apply
   ```

   Note that the apply might take several minutes.

4. Destroy the resources

   The following command destroys the resources created by `terraform apply`.

   ```bash
   terraform destroy
   ```

## Variables

The `terraform.tfvars` file must contain the following variables for the configuration to work:

- `atlas_org_id`: The ID of the Atlas organization. To learn how to retrieve an organization's details, see [View Organizations](https://www.mongodb.com/docs/atlas/access/orgs-create-view-edit-delete/#view-organizations) in the Atlas documentation.
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `project_name`: The Atlas project name (default: `"m30-multi-shard-project"`)
- `cluster_name`: The Atlas cluster name (default: `"m30-multi-shard"`)
