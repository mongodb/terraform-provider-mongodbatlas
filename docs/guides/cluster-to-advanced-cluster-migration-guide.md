---
page_title: "Migration Guide: Cluster to Advanced Cluster"
---

# Migration Guide: Cluster to Advanced Cluster

**Objective**: This guide explains how to replace the `mongodbatlas_cluster` resource with the `mongodbatlas_advanced_cluster` resource. The data source(s) migration only requires [output changes](#output-changes) as data sources only read clusters.

## Why do we have both `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` resources?

Both `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster` resources currently allow customers to manage MongoDB Atlas Clusters. Initially, only `mongodbatlas_cluster` existed. When MongoDB Atlas added support for important features such as multi-cloud deployments and the [MongoDB Atlas Admin API](https://www.mongodb.com/docs/atlas/reference/api-resources-spec/v2/) introduced versioning support, it became impossible to update this resource without causing many breaking changes for our customers. As a result, we decided to keep the resource as-is and created the `mongodbatlas_advanced_cluster` resource. Despite its name, `mongodbatlas_advanced_cluster` is intended for everyone and encompasses all functionalities, including the basic ones offered by `mongodbatlas_cluster`.
More information about the main changes between the two resources can be found [in the below section](#main-changes-between-mongodbatlas_cluster-and-mongodbatlas_advanced_cluster).

## If I am using `mongodbatlas_cluster`, why should I move to `mongodbatlas_advanced_cluster`?

Due to its schema simplicity, `mongodbatlas_cluster` resource is unable to support most of the latest MongoDB Atlas features, such as [Multi-Cloud Clusters](https://www.mongodb.com/blog/post/introducing-multicloud-clusters-on-mongodb-atlas), [Asymmetric Sharding](https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema), [Independent Scaling of Analytics Node Tiers](https://www.mongodb.com/blog/post/introducing-ability-independently-scale-atlas-analytics-node-tiers) and more.
On the other hand, not only `mongodbatlas_advanced_cluster` covers all what `mongodbatlas_cluster` can do, but it offers all existing MongoDB Atlas functionalities and will continue to do so going forward.

Having that in mind, to access all the latest functionalities and stay up to date with our best offering we recommend you to start planning your move to `mongodbatlas_advanced_cluster`.

### What is the `mongodbatlas_advanced_cluster` Preview of MongoDB Atlas Provider 2.0.0?

To facilitate the operation of moving to `mongodbatlas_advanced_cluster` for our customers, we decided to enable support for the [`moved` block](https://developer.hashicorp.com/terraform/language/moved). This functionality needs the resource to be implemented using the [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework), whereas our existing implementation [is using the SDKv2](https://developer.hashicorp.com/terraform/plugin/sdkv2). Given the considerable changes between the two frameworks and the breaking changes it causes, we have decided to release a preview version of the `mongodbatlas_advanced_cluster` usable under an environment variable and keep the existing implementation as-is to avoid breaking existing users. 
Once the MongoDB Atlas Provider 2.0.0 is released, only the new version will remain and the environment variable won't be needed.
More information about the preview version of `mongodbatlas_advanced_cluster` can be found in the [resource documentation page](../resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529).

## How should I move to `mongodbatlas_advanced_cluster`?

To move from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` we offer two alternatives:
1. (Recommended) Use the [`moved` block](https://developer.hashicorp.com/terraform/language/moved) using the Preview of MongoDB Atlas Provider 2.0.0 for `mongodbatlas_advanced_cluster`
1. Manually use the import command with the [`mongodbatlas_advanced_cluster` resource](../resources/advanced_cluster)

### Best Practices Before Migrating

Before doing any migration, create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

## Migration using the Moved block (recommended)

This is our recommended method to migrate from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster`. The [moved block](https://developer.hashicorp.com/terraform/language/moved) is a Terraform feature that allows to move between resource types. It's conceptually similar to execute `removed` and `import` commands separately but it brings the convenience of doing it in one step. It also works for `modules` and does not require direct access to the Terraform state file.

We recommend that you follow this migration process unless you can't meet the following requirements:
 - Terraform version 1.8 or later is required, more information in the [State Move page](https://developer.hashicorp.com/terraform/plugin/framework/resources/state-move).
 - MongoDB Atlas Provider version 1.29 or later is required.
 - Enable the preview for MongoDB Atlas Provider 2.0.0 of `mongodbatlas_advanced_cluster` by setting the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true`.
   - More information can be found in the [resource documentation page](../resources/advanced_cluster%2520%2528preview%2520provider%2520v2%2529).

The process to migrate from `mongodbatlas_cluster` to `mongodbatlas_advanced_cluster` using the `moved` block varies if you are using `modules` or the resource directly. Module maintainers can upgrade their implementation to `mongodbatlas_advanced_cluster` by making this operation transparent to their users. To learn how, we recommend you to review the examples from a [module maintainer](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_cluster_to_advanced_cluster/module_maintainer) and [module user](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_cluster_to_advanced_cluster/module_user) point of view.

If you are managing the resource directly, you can [refer to this example](https://github.com/mongodb/terraform-provider-mongodbatlas/tree/master/examples/migrate_cluster_to_advanced_cluster/basic).

The basic experience when using the `moved` block is as follows:
1. Before starting, run `terraform plan` to make sure that there are no planned changes.
2. Add the `mongodbatlas_advanced_cluster` resource definition.
  - Set the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true` in order to use the Preview for MongoDB Atlas Provider 2.0.0. You can also define the environment variable in your local development environment so your tools can use the new format and help you with linting and auto-completion.
  - You can use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform) to generate the `mongodbatlas_advanced_cluster` resource definition. This is the recommended method as it will generate a clean configuration keeping the original Terraform expressions. Please be aware of the [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform#limitations) and always review the generated configuration.
  - Alternatively, you can use the `terraform plan -generate-config-out=adv_cluster.tf` method to create the `mongodbatlas_advanced_cluster` resource definition, or create it manually.
3. Comment out or delete the `mongodbatlas_cluster` resource definition.
4. Add the `moved` block to your configuration file, e.g.:
```terraform
moved {
  from = mongodbatlas_cluster.this
  to   = mongodbatlas_advanced_cluster.this
}
```
5. Run `terraform plan` and make sure that there are no planned changes, only the moved block should be shown. If it shows other changes, update the `mongodbatlas_advanced_cluster` configuration until it matches the original `mongodbatlas_cluster` configuration. This is an example of no planned changes except the move:
```text
 # mongodbatlas_cluster.this has moved to mongodbatlas_advanced_cluster.this
     resource "mongodbatlas_advanced_cluster" "this" {
         name                                 = "my-cluster"
         # (24 unchanged attributes hidden)
     }

 Plan: 0 to add, 0 to change, 0 to destroy.
```

6. Run `terraform apply` to apply the changes. The `mongodbatlas_cluster` resource will be removed from the Terraform state and the `mongodbatlas_advanced_cluster` resource will be added.
7. At this moment you can delete the `moved` block from your configuration file, although it's recommended to keep it to help track the migrations.

## Migration using import

Like the previous appraoch, this method uses only [Terraform native tools](https://developer.hashicorp.com/terraform/language/import/generating-configuration) and works if you:
1. Have an existing cluster without any Terraform configuration and want to import and manage your cluster with Terraform.
2. Have existing `mongodbatlas_cluster` resource(s) but you can't use the [recommended approach](#migration-using-the-moved-block-recommended).

**Note**: We recommend the [`moved` block](#migration-using-the-moved-block-recommended) method as it's more convenient and less error-prone. If you continue with this method, you can still use the Preview for MongoDB Atlas Provider 2.0.0 by setting the environment variable `MONGODB_ATLAS_PREVIEW_PROVIDER_V2_ADVANCED_CLUSTER=true` to avoid having to migrate again in the future.

The process works as follow:
1. If you have an existing `mongodbatlas_cluster` resource, remove it from your configuration and delete it from the state file, e.g.: `terraform state rm mongodbatlas_cluster.this`. Alternatively a `removed block` (available in Terraform 1.7 and later, don't confuse with `moved block`) can be used to delete it from the state file, e.g.:
```terraform
  removed {
    from = mongodbatlas_cluster.this

    lifecycle {
      destroy = false
    }
  }
```
2. Find the import IDs of the clusters you want to migrate: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-legacy-cluster`
3. Import it using the `terraform import` command, e.g.: `terraform import mongodbatlas_advanced_cluster.this 664619d870c247237f4b86a6-legacy-cluster`. Alternatively an `import block` can be used (available in Terraform 1.5 and later, can't be used inside modules), e.g.:
```terraform
import {
  to = mongodbatlas_advanced_cluster.this
  id = "664619d870c247237f4b86a6-legacy-cluster" # from step 2
}
```
4. Run `terraform plan -generate-config-out=adv_cluster.tf`. This should generate a `adv_cluster.tf` file and display a message similar to `Plan: 1 to import, 0 to add, 0 to change, 0 to destroy`:
  ```terraform
  resource "mongodbatlas_advanced_cluster" "this" {
    # ... most attributes are removed for readability of this guide
    # ....
    backup_enabled                       = true
    cluster_type                         = "REPLICASET"
    disk_size_gb                         = 10
    name                                 = "legacy-cluster"
    project_id                           = "664619d870c247237f4b86a6"
    state_name                           = "IDLE"
    termination_protection_enabled       = false
    version_release_system               = "LTS"

    advanced_configuration {
      default_read_concern                 = null
      default_write_concern                = null
      fail_index_key_too_long              = false
      javascript_enabled                   = true
      minimum_enabled_tls_protocol         = "TLS1_2"
      no_table_scan                        = false
      oplog_min_retention_hours            = 0
      oplog_size_mb                        = 0
      sample_refresh_interval_bi_connector = 0
      sample_size_bi_connector             = 0
      transaction_lifetime_limit_seconds   = 0
    }

    replication_specs {
      container_id = {
        "AWS:US_EAST_1" = "669644ae01bf814e3d25b963"
      }
      id         = "66978026668b7619f6f48cf2"
      zone_name  = "ZoneName managed by Terraform"

      region_configs {
        priority              = 7
        provider_name         = "AWS"
        region_name           = "US_EAST_1"

        auto_scaling {
          compute_enabled            = false
          compute_max_instance_size  = null
          compute_min_instance_size  = null
          compute_scale_down_enabled = false
          disk_gb_enabled            = false
        }

        electable_specs {
          disk_iops       = 3000
          ebs_volume_type = null
          instance_size   = "M10"
          node_count      = 3
        }
        analytics_specs {
          disk_iops       = 3000
          ebs_volume_type = null
          instance_size   = "M10"
          node_count      = 1
        }
      }
    }
  }
  ```
Alternatively you can use the [Atlas CLI plugin](https://github.com/mongodb-labs/atlas-cli-plugin-terraform) to generate the `mongodbatlas_advanced_cluster` resource definition from a `mongodbatlas_cluster` definition. This will generate a clean configuration keeping the original Terraform expressions. Please be aware of the [plugin limitations](https://github.com/mongodb-labs/atlas-cli-plugin-terraform#limitations) and always review the generated configuration.

5. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
6. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value or blocks you didn't specify before, for example:
   - `advanced_configuration`
   - `connection_strings`
   - `cluster_id`
   - `bi_connector_config`
7. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration will have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
8. Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`
9. Update the references from your previous cluster resource: `mongodbatlas_cluster.this.XXXX` to the new `mongodbatlas_advanced_cluster.this.XXX`.
   - Double check [output-changes](#output-changes) to ensure the underlying configuration stays unchanged.
10.  Replace your existing clusters with the ones from `adv_cluster.tf` and run `terraform state rm mongodbatlas_cluster.this`. Without this step, Terraform will create a plan to delete your existing cluster.
11.  Remove the import block created in step 2.
12.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

### Terraform Actions
Using the `project_id` and `cluster.name`, Terraform imports your cluster and uses the new `mongodbatlas_advanced_cluster` schema to generate a configuration file. This file includes all configurable values in the schema, but none of the previous configuration defined for your `mongodbatlas_cluster`. Therefore, the new configuration will likely be a lot more verbose and contain none of your original [Terraform expressions.](https://developer.hashicorp.com/terraform/language/expressions)

## Main Changes Between `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster`

1. Replication Spec Configuration: Supports different node types (electable, analytics, read_only) where hardware configuration can differ between node types. `regions_config` is renamed to `region_configs`.
2. Provider Settings: Moved from the top level to the replication spec allowing you to create multi-cloud clusters.
3. Auto Scaling: Moved from the top level to the replication spec allowing you to scale replication specs individually.
4. Backup Configuration: Renamed from `cloud_backup` to `backup_enabled`.
5. See the [Migration Guide: Advanced Cluster New Sharding Configurations](advanced-cluster-new-sharding-schema#migration-sharded) for changes to `num_shards` and the new `zone_id`.

### Example 1: Old Configuration (`mongodbatlas_cluster`)

```terraform
resource "mongodbatlas_cluster" "this" {
  project_id   = var.project_id
  name         = "legacy-cluster"
  cluster_type = "REPLICASET"

  provider_instance_size_name = "M10" # 1 Replication Spec Configuration
  provider_name               = "AWS" # 2 Provider Settings

  auto_scaling_disk_gb_enabled = true # 3 Auto Scaling
  cloud_backup                 = true # 4 Backup Configuration

  replication_specs {
    num_shards = 1
    regions_config {
      region_name     = "US_EAST_1"
      priority        = 7
      electable_nodes = 3 # 1 Replication Spec Configuration
      analytics_nodes = 1 # 1 Replication Spec Configuration
      read_only_nodes = 0 # 1 Replication Spec Configuration
    }
  }
}
```

### Example 2: New Configuration (`mongodbatlas_advanced_cluster`)

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = var.project_id
  name           = "advanced-cluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # 4 Backup Configuration

  replication_specs {
    region_configs {
      auto_scaling { # 3 Auto Scaling
        disk_gb_enabled = true
      }
      region_name   = "US_EAST_1"
      priority      = 7
      provider_name = "AWS" # 2 Provider Settings

      electable_specs { # 1 Replication Spec Configuration
        instance_size = "M10"
        node_count    = 3
      }
      analytics_specs { # 1 Replication Spec Configuration
        instance_size = "M10"
        node_count    = 1
      }
    }
  }
}
```

### Example 3: New Configuration (`mongodbatlas_advanced_cluster`) using the Preview of MongoDB Atlas Provider 2.0.0

```terraform
resource "mongodbatlas_advanced_cluster" "this" {
  project_id     = var.project_id
  name           = "advanced-cluster"
  cluster_type   = "REPLICASET"
  backup_enabled = true # 4 Backup Configuration

  replication_specs = [
    {
      region_configs = [
        {
          auto_scaling = { # 3 Auto Scaling
            disk_gb_enabled = true
          }
          region_name   = "US_EAST_1"
          priority      = 7
          provider_name = "AWS" # 2 Provider Settings

          electable_specs = { # 1 Replication Spec Configuration
            instance_size = "M10"
            node_count    = 3
          }
          analytics_specs = { # 1 Replication Spec Configuration
            instance_size = "M10"
            node_count    = 1
          }
        }
      ]
    }
  ]
}
```

### Output Changes

- Connection strings:
  - Before: `srv_address`, `mongo_uri`, `mongo_uri_with_options` and `mongo_uri_updated` attributes were available.
  - After: They're not available anymore, use attributes in `connection_strings` instead.
- `container_id`:
  - Before: `mongodbatlas_cluster.this.replication_specs[0].container_id` was a flat string, such as: `669644ae01bf814e3d25b963`
  - After: `mongodbatlas_advanced_cluster.this.replication_specs[0].container_id` is a map, such as: `{"AWS:US_EAST_1": "669644ae01bf814e3d25b963"}`
  - If you have a single region you can access the `container_id` directly with: `one(values(mongodbatlas_advanced_cluster.this.replication_specs[0].container_id))`
- `snapshot_backup_policy`:
  - Before: It was deprecated.
  - After: Use `mongodbatlas_cloud_backup_schedule` resource instead.