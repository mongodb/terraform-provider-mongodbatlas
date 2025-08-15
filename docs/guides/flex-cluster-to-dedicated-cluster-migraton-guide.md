---
page_title: "Migration Guide: Flex Cluster to Dedicated Cluster"
---

# Migration Guide: Flex Cluster to Dedicated Cluster

**Objective**: This guide explains how to replace the `mongodbatlas_flex_cluster` resource with the `mongodbatlas_advanced_cluster` resource.

Currently, the only method to migrate your Flex cluster to a Dedicated cluster is via the Atlas UI.

<!-- Noting that implementation of flex_cluster as a part of mongodb_advanced_cluster in January 2025 will create new migration journey -->

## Best Practices Before Migrating
Before doing any migration, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

### Procedure


See [Modify a Cluster](https://www.mongodb.com/docs/atlas/scale-cluster/) for how to migrate via the Atlas UI.

Complete the following procedure to resolves the configuration drift in Terraform. This does not affect the underlying cluster infrastructure.

1. Find the import IDs of the new Dedicated cluster your Flex cluster has migrated to: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-clusterName`
2. Add an import block to one of your `.tf` files:
  ```terraform
  import {
    to = mongodbatlas_advanced_cluster.this
    id = "664619d870c247237f4b86a6-clusterName" # from step 1
  }
  ```
  3. Run `terraform plan -generate-config-out=adv_cluster.tf`. This should generate a `adv_cluster.tf` file.
  4. Run `terraform apply`. You should see the resource imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
  5. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value or blocks you didn't specify before. Required fields have been outlined in the below example resource block:
      ``` terraform
      resource "mongodbatlas_advanced_cluster" "this" {
         cluster_type = "REPLICASET"
         name         = "clusterName"
         project_id   = "664619d870c247237f4b86a6"
         replication_specs = [{
            zone_name = "Zone 1"
            region_configs = [{
               priority      = 7
               provider_name = "AWS"
               region_name   = "EU_WEST_1"
               analytics_specs = {
                  instance_size = "M10"
                  node_count    = 0
               }
               electable_specs = {
                  instance_size = "M10"
                  node_count    = 3
               }
            }]
         }]
      }
      ```
   6. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration have static values. Look in your previous configuration for:
      - variables, for example: `var.project_id`
      - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
  7. Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`
  8. Update the references from your previous cluster resource: `mongodbatlas_flex_cluster.this.X` to the new `mongodbatlas_advanced_cluster.this.X`.
  9. Update any data source blocks to refer to `mongodbatlas_advanced_cluster`.
  10. Replace your existing clusters with the ones from `adv_cluster.tf` and run `terraform state rm mongodbatlas_flex_cluster.this`. Without this step, Terraform creates a plan to delete your existing cluster.
  11.  Remove the import block created in step 2.
  12.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`
