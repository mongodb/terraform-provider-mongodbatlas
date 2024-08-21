---
page_title: "Migration Guide: Cluster to Advanced Cluster"
---

# Migration Guide: Cluster to Advanced Cluster

**Objective**: This guide explains how to replace the `mongodbatlas_cluster` resource with the `mongodbatlas_advanced_cluster` resource. The data source(s) migration only requires [output changes](#output-changes) as data sources only read clusters.

## Main Changes Between `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster`

1. Replication Spec Configuration: Supports different node types (electable, analytics, read_only) where hardware configuration can differ between node types.
2. Provider Settings: Moved from the top level to the replication spec allowing you to create multi-cloud clusters.
3. Auto Scaling: Moved from the top level to the replication spec allowing you to scale replication specs individually.
4. Backup Configuration: Renamed from `cloud_backup` to `backup_enabled`.
5. See the [Migration Guide: Advanced Cluster New Sharding Schema](advanced-cluster-new-sharding-schema#migration-sharded) for changes to `num_shards` and the new `zone_id`.

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

### Output Changes

- `container_id`:
  - Before: `mongodbatlas_cluster.this.replication_specs[0].container_id` was a flat string, such as: `669644ae01bf814e3d25b963`
  - After: `mongodbatlas_advanced_cluster.this.replication_specs[0].container_id` is a map, such as: `{"AWS:US_EAST_1": "669644ae01bf814e3d25b963"}`
  - If you have a single region you can access the `container_id` directly with: `one(values(mongodbatlas_advanced_cluster.this.replication_specs[0].container_id))`

## Best Practices Before Migrating
Before doing any migration create a backup of your [Terraform state file](https://developer.hashicorp.com/terraform/cli/commands/state).

## Migration using `terraform plan -generate-config-out=adv_cluster.tf`
This method uses only [Terraform native tools](https://developer.hashicorp.com/terraform/language/import/generating-configuration) and is ideal if you:
1. Have an existing cluster without any Terraform configuration and want to manage your cluster with Terraform.
2. Have existing `mongodbatlas_cluster` resource(s) and don't want to use an external script for migrating.

### Procedure

1. Find the import IDs of the clusters you want to migrate: `{PROJECT_ID}-{CLUSTER_NAME}`, such as `664619d870c247237f4b86a6-legacy-cluster`
2. Add an import block per cluster to one of your `.tf` files:
  ```terraform
  import {
    to = mongodbatlas_advanced_cluster.this
    id = "664619d870c247237f4b86a6-legacy-cluster" # from step 1
  }
  ```
3. Run `terraform plan -generate-config-out=adv_cluster.tf`. This should generate a `adv_cluster.tf` file and display a message similar to `Plan: 1 to import, 0 to add, 0 to change, 0 to destroy`:
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
4. Run `terraform apply`. You should see the resource(s) imported: `Apply complete! Resources: 1 imported, 0 added, 0 changed, 0 destroyed.`
5. Remove the "default" fields. Many fields of this resource are optional. Look for fields with a `null` or `0` value or blocks you didn't specify before, for example:
   - `advanced_configuration`
   - `connection_strings`
   - `cluster_id`
   - `bi_connector_config`
6. Re-use existing [Terraform expressions](https://developer.hashicorp.com/terraform/language/expressions). All fields in the generated configuration will have static values. Look in your previous configuration for:
   - variables, for example: `var.project_id`
   - Terraform keywords, for example: `for_each`, `count`, and `depends_on`
7. Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`
8. Update the references from your previous cluster resource: `mongodbatlas_cluster.this.XXXX` to the new `mongodbatlas_advanced_cluster.this.XXX`.
   - Double check [output-changes](#output-changes) to ensure the underlying configuration stays unchanged.
9.  Replace your existing clusters with the ones from `adv_cluster.tf` and run `terraform state rm mongodbatlas_cluster.this`. Without this step, Terraform will create a plan to delete your existing cluster.
1.  Remove the import block created in step 2.
2.  Re-run `terraform apply` to ensure you have no planned changes: `No changes. Your infrastructure matches the configuration.`

### Terraform Actions
Using the `project_id` and `cluster.name`, Terraform imports your cluster and uses the new `mongodbatlas_advanced_cluster` schema to generate a configuration file. This file includes all configurable values in the schema, but none of the previous configuration defined for your `mongodbatlas_cluster`. Therefore, the new configuration will likely be a lot more verbose and contain none of your original [Terraform expressions.](https://developer.hashicorp.com/terraform/language/expressions)
