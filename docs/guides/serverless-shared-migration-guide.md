---
page_title: "Migration Guide: Transition out of Serverless Instances and Shared-tier clusters"
---

# Migration Guide: Transition out of Serverless Instances and Shared-tier clusters

The goal of this guide is to help users transition out of Serverless Instances (`mongodbatlas_serverless_instance`) and Shared-tier clusters (M2 and M5) in favor of Flex or Dedicated Clusters. 

After March 2025, all Serverless instances and Shared-tier clusters will be automatically converted to Flex/Dedicated clusters, results in configuration drift in Terraform. If a Serverless instance/Shared-tier cluster does not fit the constraints of a Flex cluster, it will be converted into a Dedicated cluster which will result in downtime and workload disruption. It will otherwise be converted into a Flex cluster.

Migration from Serverless instances and Shared-tier clusters to Flex/Dedicated clusters can be done manually before March 2025. 

**Note:** We recommend to wait until March 2025 for Serverless instances and Shared-tier clusters to be autoconverted.

## From M2/M5 clusters to Flex 
### Pre-Autoconversion Migration Procedure

1. Create a new Flex Cluster directly from your `.tf` file, e.g.:

    ```terraform
    resource "mongodbatlas_flex_cluster" "flex_cluster" {
        project_id = var.project_id
        name       = "clusterName"
        provider_settings = {
            backing_provider_name = "AWS"
            region_name           = "US_EAST_1"
        }
        termination_protection_enabled = true
    }
    ```
2.  Run `terraform apply` to create the new resource.
3. Migrate data from your Shared-tier cluster to the Flex cluster using `mongodump` and `mongostore`.

    Please see the following guide on how to retrieve data from one cluster and store it in another cluster: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/)

    Verify that your data is present within the Flex cluster before proceeding.
4. Delete the Shared-tier cluster by running a destory command targetting the Shared-tier cluster:

    `terraform destroy -target=mongodbatlas_advanced_cluster.<shared-tier-cluster-name>`

 5. You may now safely remove the resource block for the Shared-tier cluster from your `.tf` file.

### Post-Autoconversion Migration Procedure

Given your Shared-tier cluster fits the constraints of a Flex cluster, it alongisde all its data will have been automatically converted into a Flex cluster in Atlas. The following will resolve the configuration drift in Terraform.

1. Find the import IDs of the Flex clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-flexClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_flex_cluster.flex_cluster
    id = "664619d870c247237f4b86a6-flexClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=flex_cluster.tf`. This should generate a `flex_cluster.tf` file with your Flex cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration will have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
6. Replace your existing clusters with the ones from `flex_cluster.tf` and run `terraform state rm mongodbatlas_advanced_cluster.advClusterName`. Without this step, Terraform will create a plan to delete your existing cluster.
7.  Remove the import block created in step 2.
8.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 


## From Serverless to Flex

### Pre-Autoconversion Migration Procedure

1. Create a new Flex Cluster directly from your `.tf` file, e.g.:

    ```terraform
    resource "mongodbatlas_flex_cluster" "flex_cluster" {
        project_id = var.project_id
        name       = "clusterName"
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
4. Delete the Shared-tier cluster by running a destory command targetting the Shared-tier cluster:

    `terraform destroy -target=mongodbatlas_serverless_instance.<serverless-instance-cluster-name>`

 5. You may now safely remove the resource block for the Serverless Instance from your `.tf` file.

### Post-Autoconversion Migration Procedure

Given your Serverless Instance fits the constraints of a Flex cluster, it alongisde all its data will have been automatically converted into a Flex cluster in Atlas. The following will resolve the configuration drift in Terraform.

1. Find the import IDs of the Flex clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-flexClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_flex_cluster.flex_cluster
    id = "664619d870c247237f4b86a6-flexClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=flex_cluster.tf`. This should generate a `flex_cluster.tf` file with your Flex cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration will have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
6. Replace your existing clusters with the ones from `flex_cluster.tf` and run `terraform state rm mongodbatlas_serverless_instance.serverlessInstanceName`. Without this step, Terraform will create a plan to delete your existing cluster.
7.  Remove the import block created in step 2.
8.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 


## From Serverless to Dedicated 

Migration from Serverless to Dedicated cannot be done using the Terraform provider.

**Note:** We recommend waiting until January 2025 for a UI based migration tool.

### Pre-Autoconversion Migration Procedure

To migrate from Serverless to Dedicated prior to January 2025, please see the following guide: [Convert a Serverless Instance to a Dedicated Cluster](https://www.mongodb.com/docs/atlas/tutorial/convert-serverless-to-dedicated/)

### Post-Autoconversion Migration Procedure
Given your Serverless Instance did not fit the constraints of a Flex cluster, it alongisde all its data will have been automatically converted into a Dedicated cluster in Atlas. The following will resolve the configuration drift in Terraform.

1. Find the import IDs of the Dedicated clusters: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-advancedClusterName`
2. Add an import block per cluster to one of your `.tf` files:
    ```terraform
    import {
    to = mongodbatlas_advanced_cluster.advanced_cluster
    id = "664619d870c247237f4b86a6-advancedClusterName" # from step 1
    }
    ```
3. Run `terraform plan -generate-config-out=adv_cluster.tf`. This should generate a `adv_cluster.tf` file with your Dedicated cluster in it.
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration will have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
6. Replace your existing clusters with the ones from `adv_cluster.tf` and run `terraform state rm mongodbatlas_serverless_instance.serverlessInstanceName`. Without this step, Terraform will create a plan to delete your existing cluster.
7.  Remove the import block created in step 2.
8.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.` 


## From Serverless to Free 
