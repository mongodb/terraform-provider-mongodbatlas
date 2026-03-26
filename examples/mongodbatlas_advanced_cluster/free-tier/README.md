# MongoDB Atlas Provider — Free Tier (M0) Cluster

This example creates an Atlas project and a **free tier M0 cluster** using the [`mongodbatlas_advanced_cluster`](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster) resource. M0 clusters are permanently free and are ideal for learning, prototyping, and early-stage development.

> **NOTE**: M0 is a shared-tier cluster with limited resources. To learn more, see [Atlas M0 (Free Cluster) Limits](https://www.mongodb.com/docs/atlas/reference/free-shared-limitations/).

## Prerequisites

To run this example, ensure you have the following tools:

- Terraform MongoDB Atlas Provider v2.0.0 or later
- A [MongoDB Atlas account](https://www.mongodb.com/docs/atlas/tutorial/create-atlas-account/).

## Procedure

To run this example, perform the following steps:

1. Set up your MongoDB Atlas credentials.

   Create a ``terraform.tfvars`` file with your credentials:

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

   - `name`: Name of the cluster, set via `var.cluster_name`. Default for this example: `"free-tier-cluster"`
   - `cluster_type`: Type of cluster. Must be `"REPLICASET"` for M0 clusters
   - `instance_size`: Instance size of the cluster. For M0 clusters, value must be `"M0"`
   - `provider_name`: Cloud service provider on which MongoDB Cloud provisions the hosts. For M0 clusters, must be `"TENANT"`
   - `backing_provider_name`: Cloud provider hosting the cluster. Default for this example: `"AWS"`
   - `region_name`: Cloud region where the cluster is deployed. Default for this example: `"US_EAST_1"`)
   - `termination_protection_enabled`: Flag that indicates whether termination protection is enabled on the cluster. If set to `true`, you can't delete the cluster using `terraform destroy`. Set to `false` for development. Set as `true` before moving to a production environment

3. Apply your configuration.

   The following command applies your configuration and creates the resources.

   ```bash
   terraform apply
   ```

4. (Optional) Destroy the resources.

   The following command destroys the resources created by `terraform apply`.

   ```bash
   terraform destroy
   ```

## Variables

The ``terraform.tfvars`` file must contain the following variables for the configuration to work:

- `org_id`: ID of the MongoDB Atlas organization. To learn how to retrieve an organization's details, see [View Organizations](https://www.mongodb.com/docs/atlas/access/orgs-create-view-edit-delete/#view-organizations) in the Atlas documentation.
- `atlas_client_id`: MongoDB Atlas Service Account Client ID
- `atlas_client_secret`: MongoDB Atlas Service Account Client Secret
