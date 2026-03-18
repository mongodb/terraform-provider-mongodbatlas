# MongoDB Atlas Provider — Free Tier (M0) Cluster

This example creates an Atlas project and a **free tier M0 cluster** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. M0 clusters are permanently free and are ideal for learning, prototyping, and early-stage development.

> **NOTE**: M0 is a shared-tier cluster with limited resources. To learn more, see [Atlas M0 (Free Cluster) Limits](https://www.mongodb.com/docs/atlas/reference/free-shared-limitations/) for more information in the Atlas documentation.

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/).

## Procedure

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

   - `name`: The name of the cluster, set via `var.cluster_name` (Default for this example: `"free-tier-cluster"`)
   - `cluster_type`: The type of cluster. Must be `"REPLICASET"` for M0 clusters
   - `instance_size`: The instance size of the cluster. For M0 clusters, must be `"M0"`
   - `provider_name`: The cloud service provider on which MongoDB Cloud provisions the hosts. For M0 clusters, must be `"TENANT"`
   - `backing_provider_name`: The cloud provider hosting the cluster (Default for this example: `"AWS"`)
   - `region_name`: The cloud region where the cluster is deployed (Default for this example: `"US_EAST_1"`)
   - `termination_protection_enabled`: Flags whether termination protection is enabled on the cluster. Set to `false` for development; mark as `enabled` before moving to production

3. Apply your configuration

   The following command applies your configuration and creates the resources.

   ```bash
   terraform apply
   ```

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
