# MongoDB Atlas Provider — M30 Multi-Shard Cluster (2 Shards)

This example creates an Atlas project and a **two-shard M30 sharded cluster** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. This is the recommended configuration for workloads that have outgrown a single replica set and need horizontal write scaling or want to distribute data across multiple shards.

Both shards are configured symmetrically (same instance size and region), which simplifies capacity planning and ensures even data distribution with a standard hashing shard key.

## Topology

The cluster distributes data across two shards as follows:

- Shard 1 — `US_EAST_1` (AWS), 3 electable nodes, instance size M30
- Shard 2 — `US_EAST_1` (AWS), 3 electable nodes, instance size M30

 **Tip:** To add a third shard, append another `replication_specs` block with the same shape as shards 1 and 2. Atlas performs a live online shard addition with no downtime.

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/)

## Procedure

To run this example, perform the following steps:

1. Set up your MongoDB Atlas credentials.

   Create a `terraform.tfvars` file with your credentials:

   ```hcl
   org_id              = "<MONGODB_ATLAS_ORG_ID>"
   atlas_client_id     = "<ATLAS_CLIENT_ID>"
   atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
   ```

   Service Accounts are the recommended way to authenticate with the MongoDB Atlas API. To learn more, see [Authentication Methods](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/provider-configuration#authentication-methods) in the MongoDB Atlas Provider documentation.

2. Review your Terraform plan.

   The following command lists the resources that your configuration creates.

   ```bash
   terraform plan
   ```

   Review the following fields in the plan output before applying:

   - `name`: Name of the cluster (default: `"m30-multi-shard"`)
   - `cluster_type`: Type of cluster. Must be `"SHARDED"`
   - `instance_size`: Instance size for all nodes. Set to `"M30"`
   - `node_count`: Number of electable nodes per shard. Set to `3`
   - `provider_name`: Cloud service provider hosting the cluster. Set to `"AWS"`
   - `region_name`: Cloud region where the cluster is deployed. Set to `"US_EAST_1"`
   - `backup_enabled`: Flag that specifies whether cloud backup is enabled. Set to `true`
   - `termination_protection_enabled`: Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, you can't delete the cluster using `terraform destroy`. Set to `true` in this example
   - `replication_specs`: Two entries, one per shard, each with identical region configuration

3. Apply your configuration.

   The following command applies your configuration and creates the resources.

   ```bash
   terraform apply
   ```

   This operation might take several minutes to complete.

4. (Optional) Destroy the resources.

   The following command destroys the resources created by `terraform apply`.

   ```bash
   terraform destroy
   ```

## Variables

The `terraform.tfvars` file must contain the following variables for the configuration to work:

- `org_id`: ID of the MongoDB Atlas organization. To learn how to retrieve an organization's details, see [View Organizations](https://www.mongodb.com/docs/atlas/access/orgs-create-view-edit-delete/#view-organizations) in the Atlas documentation
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
- `project_name`: MongoDB Atlas project name (default: `"m30-multi-shard-project"`
- `cluster_name`: MongoDB Atlas cluster name (default: `"m30-multi-shard"`
