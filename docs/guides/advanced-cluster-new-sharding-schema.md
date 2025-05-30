---
page_title: "Migration Guide: Advanced Cluster New Sharding Configurations"
---

# Migration Guide: Advanced Cluster New Sharding Configurations

**Objective**: Use this guide to migrate your existing `advanced_cluster` resources to support new sharding configurations introduced in version 1.18.0. The new sharding configurations allow you to scale shards independently. Additionally, as of version 1.23.0, compute auto-scaling supports scaling instance sizes independently for each shard when using the new sharding configuration. Existing sharding configurations continue to work, but you will receive deprecation messages if you continue to use them.

Note: Once applied, the `advanced_cluster` resource making use of the new sharding configuration will not be able to transition back to the old sharding configuration.

- [Migration Guide: Advanced Cluster New Sharding Configurations](#migration-guide-advanced-cluster-new-sharding-schema)
  - [Changes Overview](#changes-overview)
    - [Migrate advanced\_cluster type `SHARDED`](#migrate-advanced_cluster-type-sharded)
    - [Migrate advanced\_cluster type `GEOSHARDED`](#migrate-advanced_cluster-type-geosharded)
    - [Migrate advanced\_cluster type `REPLICASET`](#migrate-advanced_cluster-type-replicaset)
- [Use Independent Shard Scaling](#use-independent-shard-scaling)
- [Data Source Transition for Asymmetric Clusters](#data-source-transition-for-asymmetric-clusters)
- [Use Auto-Scaling Per Shard](#use-auto-scaling-per-shard)

## Changes Overview

`replication_specs` attribute now represents each individual cluster's shard with a unique replication spec element.
When you use the new sharding configurations, it will no longer use the existing attribute `num_shards`, and instead the number of shards are defined by the number of `replication_specs` elements.

### Migrate advanced_cluster type `SHARDED`

Consider the following configuration of a `SHARDED` cluster using the deprecated `num_shards`:
```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "SymmetricShardedCluster"
  cluster_type = "SHARDED"

  replication_specs {
    # deprecation warning will be encoutered for using num_shards
    num_shards = 2 
    region_configs {
      electable_specs {
        instance_size = "M30"
        disk_iops = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}
```

In order to use our new sharding configurations, we will remove the use of `num_shards` and add a new identical `replication_specs` element for each shard. Note that these 2 changes must be done at the same time.

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "SymmetricShardedCluster"
  cluster_type = "SHARDED"

  replication_specs { # first shard
    region_configs {
      electable_specs {
        instance_size = "M30"
        disk_iops = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # second shard
    region_configs {
      electable_specs {
        instance_size = "M30"
        disk_iops = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}
```

This updated configuration will trigger a Terraform update plan. However, the underlying cluster will not face any changes after the `apply` command, as both configurations represent a sharded cluster composed of two shards.

### Migrate advanced_cluster type `GEOSHARDED`

Consider the following configuration of a `GEOSHARDED` cluster using the deprecated `num_shards`:

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id = var.project_id
  name = "GeoShardedCluster"
  cluster_type   = "GEOSHARDED"

  replication_specs {
    zone_name  = "zone n1"
    num_shards = 2
    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_EAST_1"
    }
  }

  replication_specs {
    zone_name  = "zone n2"
    num_shards = 2

    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "EU_WEST_1"
    }
  }
}
```

In order to use our new sharding configurations, we will remove the use of `num_shards` and add a new identical `replication_specs` element for each shard. Note that these two changes must be done at the same time.

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id = var.project_id
  name = "GeoShardedCluster"
  cluster_type   = "GEOSHARDED"

  replication_specs { # first shard for zone n1
    zone_name  = "zone n1"
    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_EAST_1"
    }
  }

  replication_specs { # second shard for zone n1
    zone_name  = "zone n1"
    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "US_EAST_1"
    }
  }

  replication_specs { # first shard for zone n2
    zone_name  = "zone n2"
    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # second shard for zone n2
    zone_name  = "zone n2"
    region_configs {
    electable_specs {
        instance_size = "M30"
        node_count    = 3
    }
    provider_name = "AWS"
    priority      = 7
    region_name   = "EU_WEST_1"
    }
  }
}
```

This updated configuration triggers a Terraform update plan. However, the underlying cluster will not face any changes after the `apply` command, as both configurations represent a geo sharded cluster with two zones and two shards in each one.

### Migrate advanced_cluster type `REPLICASET`

To learn more, see the documentation on [transitioning from a replica set to a sharded cluster](https://www.mongodb.com/docs/atlas/scale-cluster/#scale-your-replica-set-to-a-sharded-cluster).

Consider the following replica set configuration:
```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "REPLICASET"

    replication_specs {
        region_configs {
            electable_specs {
                instance_size = "M30"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

To upgrade a replica set to a multi-sharded cluster, you must upgrade to a single shard cluster first, restart your application and reconnect to the cluster, and then add additional shards. If you don't restart the application clients, your data might be inconsistent once Atlas begins distributing data across shards.

First, update the `cluster_type` to SHARDED (single shard), and apply this change to the cluster.

```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "SHARDED"

    replication_specs {
        region_configs {
            electable_specs {
                instance_size = "M30"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

Once the cluster type is adjusted accordingly, you must restart the application clients. If you don't reconnect the application clients, your application may suffer from data outages.

We can now proceed to add an additional second shard:

```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "SHARDED"

    replication_specs { # first shard
        region_configs {
            electable_specs {
                instance_size = "M30"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }

    replication_specs { # second shard
        region_configs {
            electable_specs {
                instance_size = "M30"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

## Use Independent Shard Scaling 

Use the new sharding configurations. Each shard must be represented with a unique `replication_specs` element and `num_shards` must not be used, as illustrated in the following example.

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "ShardedCluster"
  cluster_type = "SHARDED"

  replication_specs { # first shard
    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # second shard
    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}
```

With each shard's `replication_specs` defined independently, we can now define distinct `instance_size` and `disk_iops` values for each shard in the cluster. Note that independent `disk_iops` values are only supported for AWS provisioned IOPS, or Azure regions that support Extended IOPS. In the following example, we define an upgraded instance size of M40 only for the first shard in the cluster.

Consider reviewing the Metrics Dashboard in the MongoDB Atlas UI (e.g. https://cloud.mongodb.com/v2/<PROJECT-ID>#/clusters/detail/ShardedCluster) for insight into how each shard within your cluster is currently performing, which will inform any shard-specific resource allocation changes you might require.

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "ShardedCluster"
  cluster_type = "SHARDED"

  replication_specs { # first shard upgraded to M40
    region_configs {
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }

  replication_specs { # second shard preserves M30
    region_configs {
      electable_specs {
        instance_size = "M30"
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
  }
}
```

## Data Source Transition for Asymmetric Clusters

When a cluster transitions to asymmetric shards, customers using data sources must update their Terraform configuration to handle the new sharding schema.

### Scenario: Cluster Becomes Asymmetric

If you have an existing cluster that becomes asymmetric due to independent shard scaling or auto-scaling per shard, you will encounter errors when using the legacy data sources.

**Error Symptoms:**
- `mongodbatlas_cluster` data source will fail with API error code `ASYMMETRIC_SHARD_UNSUPPORTED`
- `mongodbatlas_advanced_cluster` data source without `use_replication_spec_per_shard = true` will return an error asking you to enable this attribute

### Required Changes

**Before (will fail for asymmetric clusters):**
```hcl
# This will fail with ASYMMETRIC_SHARD_UNSUPPORTED error
data "mongodbatlas_cluster" "example" {
  project_id = var.project_id
  name       = "my-cluster"
}

# This will fail and ask you to set use_replication_spec_per_shard = true
data "mongodbatlas_advanced_cluster" "example" {
  project_id = var.project_id
  name       = "my-cluster"
}
```

**After (required for asymmetric clusters):**
```hcl
# Remove mongodbatlas_cluster data source completely
# Replace with mongodbatlas_advanced_cluster and enable the new schema

data "mongodbatlas_advanced_cluster" "example" {
  project_id                     = var.project_id
  name                           = "my-cluster"
  use_replication_spec_per_shard = true  # Required for asymmetric clusters
}
```

### Conditional Data Source Pattern

For modules or configurations that need to support both symmetric and asymmetric clusters, you can use conditional data source creation:

```hcl
# Example: Conditional data source based on cluster configuration
locals {
  # Determine if cluster is likely to be asymmetric based on your configuration
  cluster_uses_new_sharding = length(var.replication_specs_new) > 0
}

# Legacy cluster data source (only for symmetric clusters)
data "mongodbatlas_cluster" "this" {
  count      = local.cluster_uses_new_sharding ? 0 : 1
  name       = mongodbatlas_advanced_cluster.this.name
  project_id = mongodbatlas_advanced_cluster.this.project_id
  depends_on = [mongodbatlas_advanced_cluster.this]
}

# Advanced cluster data source (supports asymmetric clusters)
data "mongodbatlas_advanced_cluster" "this" {
  count                          = local.cluster_uses_new_sharding ? 1 : 0
  name                           = mongodbatlas_advanced_cluster.this.name
  project_id                     = mongodbatlas_advanced_cluster.this.project_id
  use_replication_spec_per_shard = true
  depends_on                     = [mongodbatlas_advanced_cluster.this]
}
```

**Notes:**
- Once a cluster becomes asymmetric, the `mongodbatlas_cluster` data source will permanently fail for that cluster.
- The `use_replication_spec_per_shard = true` attribute is required for clusters with independent shard scaling.
- This transition is necessary to take advantage of the new sharding features and avoid API compatibility issues.
- Ensure all references to the legacy data source are updated in your outputs and other resource configurations.

## Use Auto-Scaling Per Shard

As of version 1.23.0, enabled `compute` auto-scaling (either `auto_scaling` or `analytics_auto_scaling`) will scale the `instance_size` of each shard independently. Each shard must be represented with a unique `replication_specs` element and `num_shards` must not be used. On the contrary, if using deprecated `num_shards` or a lower version, enabled compute auto-scaling will scale uniformily across all shards in the cluster. 

The following example illustrates a configuration that has compute auto-scaling per shard for electable and analytic nodes.

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "AutoScalingCluster"
  cluster_type = "SHARDED"
  replication_specs { # first shard
    region_configs {
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count = 1
      }
      auto_scaling {
        compute_enabled = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling {
        compute_enabled = true
        compute_max_instance_size = "M60"
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }
  replication_specs { # second shard
    region_configs {
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count = 1
      }
      auto_scaling {
        compute_enabled = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling {
        compute_enabled = true
        compute_max_instance_size = "M60"
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }
  lifecycle { # avoids future non-empty plans as instance size start to scale from initial values
    ignore_changes = [ 
      replication_specs[0].region_configs[0].electable_specs[0].instance_size, 
      replication_specs[0].region_configs[0].analytics_specs[0].instance_size, 
      replication_specs[1].region_configs[0].electable_specs[0].instance_size,
      replication_specs[1].region_configs[0].analytics_specs[0].instance_size
    ]
  }
}
```

While the example initially defines 2 symmetric shards, auto-scaling of `electable_specs` or `analytic_specs` can lead to asymmetric shards due to changes in `instance_size`.

-> **NOTE:** In the following scenarios, a `mongodbatlas_advanced_cluster` using the new sharding configuration (single `replication_specs` per shard) might not have shard-level auto-scaling enabled:
1. Configuration was defined prior to version 1.23.0 when auto-scaling per shard feature was released.
2. Cluster was imported from a legacy schema (For example, `mongodbatlas_cluster` or `mongodbatlas_advanced_cluster` using `num_shards` > 1).
In these cases, you must update the cluster configuration to activate the auto-scaling per shard feature. This can be done by temporarily modifying a value like `compute_min_instance_size`.

-> **NOTE:** See the table [below](#resources-and-data-sources-impacted-by-independent-shard-scaling) for other impacted resources when a cluster transitions to independently scaled shards.

## Resources and Data Sources Impacted by Independent Shard Scaling

Name | Changes | Transition Guide
--- | --- | ---
`mongodbatlas_advanced_cluster` | Data source must use the `use_replication_spec_per_shard` attribute. | -
`mongodbatlas_advanced_cluster` | Use `replication_specs.#.zone_id` instead of `replication_specs.#.id`. | -
`mongodbatlas_cluster` | Resource and Data Source will not work. API error code `ASYMMETRIC_SHARD_UNSUPPORTED`. | [cluster-to-advanced-cluster-migration-guide](cluster-to-advanced-cluster-migration-guide.md)
`mongodbatlas_cloud_backup_schedule` | Use `copy_settings.#.zone_id` instead of `copy_settings.#.replication_spec_id` | [1.18.0 Migration Guide](1.18.0-upgrade-guide.md#transition-cloud-backup-schedules-for-clusters-to-use-zones)
`mongodbatlas_global_cluster_config` | `custom_zone_mapping` is no longer populated, `custom_zone_mapping_zone_id` must be used instead. | -
