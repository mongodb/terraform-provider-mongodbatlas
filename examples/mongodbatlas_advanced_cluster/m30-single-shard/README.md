# MongoDB Atlas Provider — Simple M30 Single-Shard Cluster

This example creates an Atlas project and a **single-shard M30 sharded cluster** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. Starting with a single-shard `SHARDED` cluster (rather than a `REPLICASET`) allows you to scale horizontally by adding more shards later without migrating data or recreating the cluster.

This is a good fit for applications that expect significant data growth or write volume over time and want to avoid a future migration from a replica set to a sharded cluster.

## Topology

The cluster runs a single shard as follows:

- Shard 1 — `US_EAST_1` (AWS), three electable nodes, instance size M30

Each shard runs as a three-node replica set internally, so all the automatic failover and data redundancy guarantees of a replica set apply to each shard.

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/)

## Procedure

To run this example, perform the following steps:

1. Set Up your MongoDB Atlas Credentials.

   Create a `terraform.tfvars` file with your credentials:

   ```hcl
   atlas_org_id        = "<MONGODB_ATLAS_ORG_ID>"
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

   - `name`: The name of the cluster (default: `"m30-single-shard"`)
   - `cluster_type`: The type of cluster. Must be `"SHARDED"`
   - `instance_size`: The instance size for all nodes. Set to `"M30"`
   - `node_count`: The number of electable nodes per shard. Set to `3`
   - `provider_name`: The cloud service provider hosting the cluster. Set to `"AWS"`
   - `region_name`: The cloud region where the cluster is deployed. Set to `"US_EAST_1"`
   - `backup_enabled`: Flag that specifies whether cloud backup is enabled. Set to `true`
   - `termination_protection_enabled`: Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, you can't delete the cluster using `terraform destroy`. Set to `false` for development. Set as `true` before moving to a production environment.

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

- `atlas_org_id`: The ID of the Atlas organization. To learn how to retrieve an organization's details, see [View Organizations](https://www.mongodb.com/docs/atlas/access/orgs-create-view-edit-delete/#view-organizations) in the Atlas documentation.
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `project_name`: The Atlas project name (default: `"m30-single-shard-project"`)
- `cluster_name`: The Atlas cluster name (default: `"m30-single-shard"`)
