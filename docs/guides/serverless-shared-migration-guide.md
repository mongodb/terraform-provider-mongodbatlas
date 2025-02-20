---
page_title: "Migration Guide: Transition out of Serverless Instances and Shared-tier clusters"
---

# Migration Guide: Transition out of Serverless Instances and Shared-tier clusters

The goal of this guide is to help users transition from Serverless Instances and Shared-tier clusters (M2/M5) to Free, Flex or Dedicated Clusters. 

Starting in January 2025 or later, all Shared-tier clusters (in both `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster`) will automatically convert to Flex clusters. Similarly, in March 2025 all Serverless instances (`mongodbatlas_serverless_instance`) will be converted into Free/Flex/Dedicated clusters, [depending on your existing configuration](https://www.mongodb.com/docs/atlas/flex-migration/).
If a Serverless instance has $0 MRR, it automatically converts into a Free cluster. Else, if it does not fit the constraints of a Flex cluster, it will convert into a Dedicated cluster, resulting in downtime and workload disruption. Otherwise, it will convert to a Flex cluster.
Some of these conversions will result in configuration drift in Terraform. 


You can migrate from Serverless instances and Shared-tier clusters manually before autoconversion. 

**--> NOTE:** We recommend waiting until March 2025 or later for Serverless instances and Shared-tier clusters to autoconvert.

For Shared-tier clusters, we are working on enhancing the User Experience such that Terraform Atlas Providers users can make even fewer required changes to their scripts from what is shown below. More updates to come over the next few months.

### Jump to:
- [Shared-tier to Flex](#from-shared-tier-clusters-to-flex)
- [Serverless to Free](#from-serverless-to-free)
- [Serverless to Flex](#from-serverless-to-flex)
- [Serverless to Dedicated](#from-serverless-to-dedicated)

## From Shared-tier clusters to Flex 

### Post-Autoconversion Migration Procedure

Shared-tier clusters will automatically convert in January 2025 or later to Flex clusters in Atlas, retaining all data. We recommend that you migrate to a Flex cluster managed by `mongodbatlas_advanced_cluster` resource once the autoconversion is done.

The following steps explain how to move your exising Shared-tier cluster resource to a flex cluster using `mongodbatlas_advanced_cluster` resource and does not affect the underlying cluster infrastructure:

1. Change the configuration of your Shared-tier cluster to a Flex cluster. For more details on the exact changes, see the [Example Tenant Cluster Upgrade to Flex](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/resources/advanced_cluster#example-tenant-cluster-upgrade-to-flex)
3. Run `terraform plan` to see the planned changes.
4. Run `terraform apply`. This will upgrade your Shared-tier cluster to a Flex tier cluster.
10. Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 
11. Change your usages of `mongodbatlas_shared_tier_restore_job`, `mongodbatlas_shared_tier_restore_jobs`, `mongodbatlas_shared_tier_snapshot` and `mongodbatlas_shared_tier_snapshots` data sources to the new `mongodbatlas_flex_restore_job`, `mongodbatlas_flex_restore_jobs`, `mongodbatlas_flex_snapshot` and `mongodbatlas_flex_snapshot` respectively.

### Pre-Autoconversion Migration Procedure

**NOTE:** We recommend waiting until January 2025 or later for Shared-tier clusters to autoconvert. Manually doing the migration can cause downtime and workload disruption.

1. Create a new Flex Cluster directly from your `.tf` file, e.g.:

    ```terraform
    resource "mongodbatlas_flex_cluster" "this" {
        project_id = var.project_id
        name       = "flexClusterName"
        provider_settings = {
            backing_provider_name = "AWS"
            region_name           = "US_EAST_1"
        }
        termination_protection_enabled = true
    }
    ```
2. Run `terraform apply` to create the new resource.
3. Migrate data from your Shared-tier cluster to the Flex cluster using `mongodump` and `mongostore`.

    Please see the following guide on how to retrieve data from one cluster and store it in another cluster: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/)

    Verify that your data is present within the Flex cluster before proceeding.
4. Delete the Shared-tier cluster by running a destroy command against it.
    
    For *mongodbatlas_advanced_cluster*:

    `terraform destroy -target=mongodbatlas_advanced_cluster.this`

    For *mongodbatlas_cluster*:

    `terraform destroy -target=mongodbatlas_cluster.this`

 5. Remove the resource block for the Shared-tier cluster from your `.tf` file.

## From Serverless to Free 

**Please ensure your Serverless instance meets the following requirements to migrate to Free:**
- $0 MRR

### Post-Autoconversion Migration Procedure

Given your Serverless Instance has $0 MRR, it will automatically convert in March 2025 into a Free cluster in Atlas, retaining all data.

The following steps resolve the configuration drift in Terraform without affecting the underlying cluster infrastructure:

1. Find the import IDs of the Free clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-freeClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_advanced_cluster.this
    id = "664619d870c247237f4b86a6-freeClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=free_cluster.tf`. This should generate a `free_cluster.tf` file with your Free cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value.
6. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
7. Update the references from your previous cluster resource: `mongodbatlas_serverless_instance.this.X` to the new `mongodbatlas_advanced_cluster.this.X`.
8. Update any shared-tier data source blocks to refer to `mongodbatlas_advanced_cluster`.
9. Replace your existing clusters with the ones from `free_cluster.tf` and run `terraform state rm mongodbatlas_serverless_instance.this`. Without this step, Terraform creates a plan to delete your existing cluster.
10.  Remove the import block created in step 2.
11.  Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 

### Pre-Autoconversion Migration Procedure

**NOTE:** We recommend waiting until March 2025 or later for Serverless instances to autoconvert. Manually doing the migration can cause downtime and workload disruption.

1. Create a new Free Cluster directly from your `.tf` file, e.g.:

    ```terraform
    resource "mongodbatlas_advanced_cluster" "this" {
        project_id   = var.atlas_project_id
        name         = "freeClusterName"
        cluster_type = "REPLICASET"

        replication_specs {
            region_configs {
            electable_specs {
                instance_size = "M0"
            }
            provider_name         = "TENANT"
            backing_provider_name = "AWS"
            region_name           = "US_EAST_1"
            priority              = 7
            }
        }
    }
    ```
2.  Run `terraform apply` to create the new resource.
3. Migrate data from your Serverless Instance to the Free cluster using `mongodump` and `mongostore`.

    Please see the following guide on how to retrieve data from one cluster and store it in another cluster: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/)

    Verify that your data is present within the Free cluster before proceeding.
4. Delete the Serverless Instance by running a destroy command against the Serverless Instance:

    `terraform destroy -target=mongodbatlas_serverless_instance.this`

 5. Remove the resource block for the Serverless Instance from your `.tf` file.

## From Serverless to Flex

**Please ensure your Serverless instance meets the following requirements to migrate to Flex:**
- <= 5GB of data
- no privatelink or continuous backup
- < 500 ops/sec consistently.

### Post-Autoconversion Migration Procedure

Given your Serverless Instance fits the constraints of a Flex cluster, it will automatically convert in March 2025 into a Flex cluster in Atlas, retaining all data. We recommend to migrate to `mongodbatlas_flex_cluster` resource once the autoconversion is done.

The following steps explain how to move your exising Serverless instance resource to a Flex cluster using `mongodbatlas_advanced_cluster` resource and does not affect the underlying cluster infrastructure:

1. Find the import IDs of the Flex clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-flexClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_advanced_cluster.flex
    id = "664619d870c247237f4b86a6-flexClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=flex_cluster.tf`. This should generate a `flex_cluster.tf` file with your Flex cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value.
6. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
7. Update the references from your previous cluster resource: `mongodbatlas_serverless_instance.this.X` to the new `mongodbatlas_advanced_cluster.flex.X`.
8. Replace your existing clusters with the ones from `flex_cluster.tf` and run `terraform state rm mongodbatlas_serverless_instance.this`. **IMPORTANT**: Without this step, Terraform creates a plan to delete your existing cluster.
9.  Remove the import block created in step 2.
10.  Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 

### Pre-Autoconversion Migration Procedure

**NOTE:** We recommend waiting until March 2025 or later for Serverless instances to autoconvert. Manual migration can cause downtime and workload disruption.

1. Create a new Flex Cluster directly from your `.tf` file, e.g.:

    ```terraform
    resource "mongodbatlas_flex_cluster" "this" {
        project_id = var.project_id
        name       = "flexClusterName"
        provider_settings = {
            backing_provider_name = "AWS"
            region_name           = "US_EAST_1"
        }
        termination_protection_enabled = true
    }
    ```
2.  Run `terraform apply` to create the new resource.
3. Migrate data from your Serverless Instance to the Flex cluster using `mongodump` and `mongostore`.

    Please see the following guide on how to retrieve data from one cluster and store it in another cluster: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/)

    Verify that your data is present within the Flex cluster before proceeding.
4. Delete the Serverless Instance by running a destroy command against it:

    `terraform destroy -target=mongodbatlas_serverless_instance.this`

 5. You may now safely remove the resource block for the Serverless Instance from your `.tf` file.

## From Serverless to Dedicated 

**Please note your Serverless instance will need to migrate to Decidated if it meets the following requirements:**
- \>= 5GB of data
- needs privatelink or continuous backup
- \> 500 ops/sec consistently.

You cannot migrate from Serverless to Dedicated using the Terraform provider.

### Pre-Autoconversion Migration Procedure

**NOTE:** In early 2025, we will release a UI-based tool for migrating your workloads from Serverless instances to Dedicated clusters. This tool will ensure correct migration with little downtime. You won't need to change connection strings.

To migrate from Serverless to Dedicated prior to early 2025, please see the following guide: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/). **NOTE:** Manual migration can cause downtime and workload disruption.

### Post-Autoconversion Migration Procedure

**NOTE:** Auto-conversion from Serverless to Dedicated will cause downtime and workload disruption. This guide is only valid after the auto-conversion is done.

Given your Serverless Instance doesn't fit the constraints of a Flex cluster, it will automatically convert in March 2025 into a Dedicated cluster in Atlas, retaining all data.

The following steps resolve the configuration drift in Terraform and does not affect the underlying cluster infrastructure:

1. Find the import IDs of the Dedicated clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-advancedClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_advanced_cluster.this
    id = "664619d870c247237f4b86a6-advancedClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=dedicated_cluster.tf`. This should generate a `dedicated_cluster.tf` file with your Dedicated cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value.
6. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
7. Update the references from your previous cluster resource: `mongodbatlas_serverless_instance.this.X` to the new `mongodbatlas_advanced_cluster.this.X`.
8. Update any shared-tier data source blocks to refer to `mongodbatlas_advanced_cluster`.
9. Replace your existing clusters with the ones from `dedicated_cluster.tf` and run `terraform state rm mongodbatlas_serverless_instance.this`. Without this step, Terraform creates a plan to delete your existing cluster.
10.  Remove the import block created in step 2.
11.  Re-run `terraform plan` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 
