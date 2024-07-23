# Cluster to Advanced Cluster Migration Guide

**Objective**: Guide users to replace the `mongodbatlas_cluster` resource with the `mongodbatlas_advanced_cluster` resource.

-> **NOTE:** This guide focus on the resource migration as the data source migration is only requiring a resource_type change from `data.mongodbatlas_cluster` to `data.mongodbatlas_advanced_cluster`.  However, pay attention to the [output changes.](#output-changes)

## Main Changes Between `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster`

Key differences in the configuration:

1. Replication Spec Configuration: Supports different node types (electable, analytics, read_only) where hardware configuration can differ between node types.
2. Provider Settings: Moved from the top level to the replication spec allowing you to create multi-cloud clusters.
3. Auto Scaling: Moved from the top level to the replication spec allowing you to scale replication specs individually.
4. Backup Configuration: Renamed from `cloud_backup` to `backup_enabled`.
5. See also ["Migration to new sharding schema and leveraging Independent Shard Scaling"](/guides/advanced-cluster-new-sharding-schema#migration-sharded)

### Example 1: Old Configuration (`mongodbatlas_cluster`)

```terraform
resource "mongodbatlas_cluster" "this" {
    project_id                   = var.project_id
    name                         = "legacy-cluster"
    cluster_type                 = "REPLICASET"

    # Provider Settings "block"
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
    project_id             = var.project_id
    name                   = "advanced-cluster"
    cluster_type           = "REPLICASET"
    backup_enabled         = true # 4 Backup Configuration

    # Replication specs
    replication_specs {
      auto_scaling { # 3 Auto Scaling
        disk_gb_enabled = true
      }
      region_configs {
        region_name   = "US_EAST_1"
        priority    = 7
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
```

### Output Changes

- container_id:
  - before: `mongodbatlas_cluster.this.replication_specs[0].container_id` was a flat string, e.g., `669644ae01bf814e3d25b963`
  - after: `mongodbatlas_advanced_cluster.this.replication_specs[0].container_id` is a map, e.g., `{"AWS:US_EAST_1": "669644ae01bf814e3d25b963"}`

## How to Change
Before doing any migration it is recommended to make a backup of your Terraform state files. (ADD_LINK_HERE)

### Method 1: `terraform plan -generate-config-out=cluster.tf`
This method uses only Terraform native tools and are ideal for customers who:
1. Have an existing cluster without any Terraform configuration and want to manage their cluster with Terraform.
2. Have existing `mongodbatlas_cluster` resource(s) and don't want to use an external script for migrating.

#### Steps

1. Find the import IDs of the clusters you want to migrate: `{PROJECT_ID}-{CLUSTER_NAME}`, e.g., `664619d870c247237f4b86a6-legacy-cluster`
2. Add an import block per cluster to one of your `.tf` files:
  ```terraform
  import {
    to = mongodbatlas_advanced_cluster.this
    id = "664619d870c247237f4b86a6-legacy-cluster" # from step 1
  }
  ```
3. Run `terraform plan -generate-config-out=cluster.tf`, should generate a `cluster.tf` file and display a message similar to `Plan: 1 to import, 0 to add, 0 to change, 0 to destroy`:
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
4. Run a `terraform apply`, you should expect to see the resource imported.
5. Remove the "default" fields. Many fields of this resource are optional, look for fields with a `null` or `0` value or blocks you didn't specify before, e.g:
   - `advanced_configuration`
   - `connection_strings`
   - `cluster_id`
   - `bi_connector_config`
6. Re-run `terraform apply` to ensure you have no plan changes.
7. Update the references from your old cluster resource: `mongodbatlas_cluster.this.XXXX` to the new `mongodbatlas_advanced_cluster.this.XXX`.
   - Double check [output-changes](#output-changes) to ensure the meaning stays unchanged
8. Replace your existing clusters with the ones from `cluster.tf` and run `terraform state rm mongodbatlas_cluster.this`.
9. Re-run `terraform apply` to ensure you have no plan changes.
