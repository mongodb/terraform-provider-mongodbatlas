# MongoDB Atlas Provider — Simple M10 Replica Set

This example creates an Atlas project and a **three-node M10 replica set** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. M10 is the recommended entry-level dedicated cluster for development, internal tooling, and low-volume production workloads.

## Topology

A three-node replica set provides:

- Automatic failover (a new primary is elected if the current one goes down)
- Data redundancy across three copies
- Sufficient votes (three) to always reach a majority

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/)

## Procedure

Follow the next steps to run this example:

1. Set Up your MongoDB Atlas Credentials

   Create a ``terraform.tfvars`` file with your credentials:

   ```hcl
   atlas_org_id        = "<MONGODB_ATLAS_ORG_ID>"
   atlas_client_id     = "<ATLAS_CLIENT_ID>"
   atlas_client_secret = "<ATLAS_CLIENT_SECRET>"
   ```

   Service Accounts are the recommended way to authenticate with the MongoDB Atlas API. To learn how to create a Service Account, see [Manage Service Accounts](https://www.mongodb.com/docs/atlas/configure-api-access/#manage-service-accounts) in the Atlas documentation.

2. Review your Terraform Plan

   The following command lists the resources that your configuration will create.

   ```bash
   terraform plan
   ```

   Review the following fields in the plan output before applying:

   - `name`: The name of the cluster, set via `var.cluster_name` (Default for this example: `"m10-replicaset"`)
   - `cluster_type`: The type of cluster. Must be `"REPLICASET"`
   - `instance_size`: The instance size of the cluster. Set to `"M10"`
   - `node_count`: The number of electable nodes. Set to `3`
   - `provider_name`: The cloud service provider hosting the cluster (Default for this example: `"AWS"`)
   - `region_name`: The cloud region where the cluster is deployed (Default for this example: `"US_EAST_1"`)
   - `backup_enabled`: Flags whether cloud backup is enabled. Set to `true`
   - `termination_protection_enabled`: Flags whether termination protection is enabled. (Default for this example: `false`)

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

The ``terraform.tfvars`` file must contain the following variables for the configuration to work:

- `atlas_org_id`: The ID of the Atlas organization. To learn how to retrieve an organization's details, see [View Organizations](https://www.mongodb.com/docs/atlas/access/orgs-create-view-edit-delete/#view-organizations) in the Atlas documentation.
- `atlas_client_id`: The MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: The MongoDB Atlas Service Account Client Secret
- `project_name`: The Atlas project name (default: `"m10-replicaset-project"`)
- `cluster_name`: The Atlas cluster name (default: `"m10-replicaset"`)
