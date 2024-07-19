# Cluster to Advanced Cluster Migration Guide

- Outline
  - Motivation
  - Main Changes
  - How to change

**Objective**: Guide users to replace the `mongodbatlas_cluster` resource with the `mongodbatlas_advanced_cluster` resource.

**Note**: This guide focus on the resource migration as the data source migration is only requiring a resource_type change from `data.mongodbatlas_cluster` to `data.mongodbatlas_advanced_cluster`.  However, pay attention to the [output changes.](#output-changes)

## Motivations for migrating
- Access to new features:
  - [Multi Cloud Clusters](https://www.mongodb.com/resources/basics/multicloud) to increase availability of your cluster
  - Auto scaling features to help scale your clusters ADD_LINK_HERE
  - Advanced hardware configuration to reduce costs without loosing performance ADD_LINK_HERE
  - Upcoming new features
- Future proof your cluster:
  - Avoid deprecation warnings
  - No problems when the `mongodbatlas_cluster` reaches end of life

## Main Changes Between `mongodbatlas_cluster` and `mongodbatlas_advanced_cluster`

Key differences in the configuration:

1. Resource Type: mongodbatlas_cluster vs mongodbatlas_advanced_cluster.
2. Replication Configuration: The advanced cluster allows you to specify detailed replication specs, including auto_scaling, node counts, instance sizes, disk_iops, disk_size_gb, etc. for different node types (electable, analytics, read_only).
3. Provider Settings: In the standard cluster, provider settings are specified at the top level. In the advanced cluster, they're part of the region_configs block within replication_specs.
4. Backup Configuration: In the standard cluster, it's a simple boolean (cloud_backup). In the advanced cluster, it's backup_enabled.

### Old Configuration Example

```hcl
resource "mongodbatlas_cluster" "standard_cluster" {
    project_id             = var.project_id
    name                   = "standard-cluster"
    cluster_type           = "REPLICASET"
    mongo_db_major_version = "5.0"
    cloud_backup           = true

    # Provider Settings "block"
    provider_name               = "AWS"
    provider_region_name        = "US_EAST_1"
    provider_instance_size_name = "M10"
    }
```

### New Configuration Example

```hcl
    resource "mongodbatlas_advanced_cluster" "advanced_cluster" {
    project_id             = var.project_id
    name                   = "advanced-cluster"
    cluster_type           = "REPLICASET"
    mongo_db_major_version = "5.0"
    backup_enabled         = true

    # Replication specs
    replication_specs {
        num_shards = 1
        region_configs {
        electable_specs {
            instance_size = "M10"
            node_count    = 3
        }
        analytics_specs {
            instance_size = "M10"
            node_count    = 1
        }
        priority    = 7
        provider_name = "AWS"
        region_name   = "US_EAST_1"
        }
    }
```

### Output Changes


## How to change

TODO: Will continue from here by creating a step-by-step guide

1. find the id of clusters to replace
2. add import blocks and execute command to generate config
3. remove unused default attributes
4. import new cluster
5. change output references
6. delete old cluster block