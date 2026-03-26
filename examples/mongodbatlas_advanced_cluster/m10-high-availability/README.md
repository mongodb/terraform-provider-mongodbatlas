# MongoDB Atlas Provider — M10 High-Availability Replica Set (2-2-1)

This example creates an Atlas project and a **five-node M10 replica set spread across three AWS regions** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. This pattern is a cost-effective way to achieve high availability with protection against a full regional outage, without the cost of a sharded cluster.

## Topology

The cluster distributes nodes across three regions as follows:

- `US_EAST_1` (AWS) — Two electable nodes, priority 7 (preferred primary region)
- `US_WEST_2` (AWS) — Two electable nodes, priority 6 (secondary failover region)
- `EU_WEST_1` (AWS) — One electable node, priority 5 (tiebreaker)

### Why the odd node (2-2-1)?

A four-node cluster (2+2) would have an even number of votes. If both regions in a 2-node group are simultaneously affected, neither side can reach a simple majority, resulting in an election stalemate. The fifth node in a third region adds one extra vote, ensuring a majority (3 of 5) is always reachable even if any single region goes down.

Failover scenarios with five votes (majority = 3):

- `US_EAST_1` lost — Three remaining votes (US_WEST_2: 2 + EU_WEST_1: 1) — election succeeds
- `US_WEST_2` lost — Three remaining votes (US_EAST_1: 2 + EU_WEST_1: 1) — election succeeds
- `EU_WEST_1` lost — Four remaining votes (US_EAST_1: 2 + US_WEST_2: 2) — election succeeds

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

   - `name`: Name of the cluster. Default for this example: `"m10-high-availability"`
   - `cluster_type`: Type of cluster. Must be `"REPLICASET"`
   - `instance_size`: Instance size for all nodes. Set to `"M10"`
   - `provider_name`: Cloud service provider hosting the cluster. Set to `"AWS"`
   - `backup_enabled`: Flag that specifies whether cloud backup is enabled. Set to `true`
   - `termination_protection_enabled`: Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, you can't delete the cluster using `terraform destroy`. Set to `true` in this example
   - `region_configs`: Cloud provider region configuration. Three entries — `US_EAST_1` (2 nodes, priority 7), `US_WEST_2` (2 nodes, priority 6), `EU_WEST_1` (1 node, priority 5)

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
- `project_name`: MongoDB Atlas project name. Default for this example: `"m10-ha-project"`
- `cluster_name`: MongoDB Atlas cluster name. Default for this example: `"m10-high-availability"`
