---
page_title: "Migration Guide: Advanced Cluster New Sharding Configurations"
---

# Migration Guide: Advanced Cluster New Sharding Configurations

**Objective**: Use this guide to migrate your existing `advanced_cluster` resources to support new sharding configurations introduced in version 1.18.0. The new sharding configurations allow you to scale shards independently. Existing sharding configurations continue to work, but you will receive deprecation messages if you continue to use them.

- [Migration Guide: Advanced Cluster New Sharding Configurations](#migration-guide-advanced-cluster-new-sharding-schema)
  - [Changes Overview](#changes-overview)
    - [Migrate advanced\_cluster type `SHARDED`](#migrate-advanced_cluster-type-sharded)
    - [Migrate advanced\_cluster type `GEOSHARDED`](#migrate-advanced_cluster-type-geosharded)
    - [Migrate advanced\_cluster type `REPLICASET`](#migrate-advanced_cluster-type-replicaset)
    - [Use Independent Shard Scaling](#use-independent-shard-scaling)

<a id="overview"></a>
## Changes Overview

`replication_specs` attribute now represents each individual cluster's shard with a unique replication spec element.
When you use the new sharding configurations, it will no longer use the existing attribute `num_shards`, and instead the number of shards are defined by the number of `replication_specs` elements.

<a id="migration-sharded"></a>
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

Note: The first time `terraform apply` command is run **after** updating the configuration, you may receive a `500 Internal Server Error (Error code: "SERVICE_UNAVAILABLE")` error. This is a known temporary issue. If you encounter this, please re-run `terraform apply` and this time the update should succeed. 

<a id="migration-geosharded"></a>
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
        instance_size = "M10"
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
        instance_size = "M10"
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
        instance_size = "M10"
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
        instance_size = "M10"
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
        instance_size = "M10"
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
        instance_size = "M10"
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

Note: The first time `terraform apply` command is run **after** updating the configuration, you may receive a `500 Internal Server Error (Error code: "SERVICE_UNAVAILABLE")` error. This is a known temporary issue. If you encounter this, please re-run `terraform apply` and this time the update should succeed. 

<a id="migration-replicaset"></a>
### Migrate advanced_cluster type `REPLICASET`

To learn more, see the documentation on [transitioning from a replicaset to a sharded cluster](https://www.mongodb.com/docs/atlas/scale-cluster/#convert-a-replica-set-to-a-sharded-cluster).

Consider the following replica set configuration:
```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "REPLICASET"

    replication_specs {
        region_configs {
            electable_specs {
                instance_size = "M10"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

To transition a replica set to sharded cluster 2 separate updates must be applied. First, update the `cluster_type` to SHARDED, and apply this change to the cluster.

```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "SHARDED"

    replication_specs {
        region_configs {
            electable_specs {
                instance_size = "M10"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

Once the cluster type is adjusted accordingly, we can proceed to add a new shard:

```
resource "mongodbatlas_advanced_cluster" "test" {
    project_id   = var.project_id
    name         = "ReplicaSetTransition"
    cluster_type = "SHARDED"

    replication_specs { # first shard
        region_configs {
            electable_specs {
                instance_size = "M10"
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
                instance_size = "M10"
                node_count    = 3
            }
            provider_name = "AZURE"
            priority      = 7
            region_name   = "US_EAST"
        }
    }
}
```

Note: The first time `terraform apply` command is run **after** updating the configuration, you may receive a `500 Internal Server Error (Error code: "SERVICE_UNAVAILABLE")` error. This is a known temporary issue. If you encounter this, please re-run `terraform apply` and this time the update should succeed. 

<a id="use-iss"></a>
### Use Independent Shard Scaling 

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

With each shard's `replication_specs` defined independently, we can now define distinct `instance_size`, and `disk_iops` (only for AWS) values for each shard in the cluster. In the following example, we define an upgraded instance size of M40 only for the first shard in the cluster.

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

-> **NOTE:** For any cluster leveraging the new sharding configurations and defining independently scaled shards, users should also update corresponding `mongodbatlas_cloud_backup_schedule` resource & data sources. This involves updating any existing Terraform configurations of the resource to use `copy_settings.#.zone_id` instead of `copy_settings.#.replication_spec_id`. This is needed as `mongodbatlas_advanced_cluster` resource and data source will no longer have `replication_specs.#.id` present when shards are scaled independently. To learn more, review the [1.18.0 Migration Guide](1.18.0-upgrade-guide.md#transition-cloud-backup-schedules-for-clusters-to-use-zones).
