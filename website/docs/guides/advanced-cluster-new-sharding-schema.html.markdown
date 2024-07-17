---
page_title: "advanced_cluster - Migration to new sharding schema and leveraging Independent Shard Scaling"
---

**Objective**: Guide users to migrate their existing advanced_cluster configurations to use the new sharding schema which was introduced in version `1.18.0`. Additionally, a section is included to describe how Independent Shard Scaling can be used once the new sharding schema is adopted. Exiting sharding configurations will continue to work but will have deprecation messages if not using the new sharding schema.

- [Overview of schema changes](#overview)
- [Migrating existing advanced_cluster type SHARDED](#migration-sharded)
- [Migrating existing advanced_cluster type GEOSHARDED](#migration-geosharded)
- [Migrating existing advanced_cluster type REPLICASET](#migration-replicaset)
- [Leveraging Independent Shard Scaling](#leveraging-iss)

<a id="overview"></a>
# Overview of schema changes

`replication_specs` attribute has been modified to now be able to represent each individual shard of a cluster with a unique replication spec element. This implies that when using the new sharding schema the existing attribute `num_shards` will no longer be defined, and instead the number of shards will be defined by the number of `replication_specs` elements.

<a id="migration-sharded"></a>
## Migrating existing advanced_cluster type SHARDED

Considering the following configuration of a SHARDED cluster using the deprecated `num_shards`:
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

In order to update our configuration to the new schema we will remove the use of `num_shards` and add a new identical `replication_specs` element for each shard. Note that these 2 changes must be done at the same time.

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

This  updated configuration will trigger a terraform update plan, however the underlying cluster will not face any changes after the apply as both configurations represent a sharded cluster composed of 2 shards.

<a id="migration-geosharded"></a>
## Migrating existing advanced_cluster type GEOSHARDED

Considering the following configuration of a GEOSHARDED cluster using the deprecated `num_shards`:

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

In order to update our configuration to the new schema we will remove the use of `num_shards` and add a new identical `replication_specs` element for each shard. Note that these 2 changes must be done at the same time.

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



This  updated configuration will trigger a terraform update plan, however the underlying cluster will not face any changes after the apply as both configurations represent a geo sharded cluster with 2 zones and 2 shards in each one.

<a id="migration-replicaset"></a>
## Migrating existing advanced_cluster type REPLICASET

-> **NOTE:**  Please consider the following complementary documentation providing details on transitioning from a replicaset to a sharded cluster: https://www.mongodb.com/docs/atlas/scale-cluster/#convert-a-replica-set-to-a-sharded-cluster.

Considering the following replica set configuration:
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

To transition a replica set to sharded cluster 2 separate updates must be applied. First the cluster type must be adjusted to SHARDED, and apply this change without any other additional changes.

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

Once the cluster type is adjusted accordingly, we can proceed to add a new shard using the new schema:

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

<a id="leveraging-iss"></a>
## Leveraging Independent Shard Scaling 

Prerequisite: The new sharding schema must be used, meaning the advanced_cluster configuration is not using `num_shards` and therefore each shard is represented with a unique `replication_specs` element. Please refer to documentation above to transition into this schema if not already.

Considering the following configuration of a SHARDED cluster that has a symmetric configuration for its 2 shards:

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "SymmetricShardedCluster"
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

It is now possible to define a different `instance_size`, and `disk_iops` in the case of AWS, for each individual shard. In the following update we will define an upgraded instance size of M40 only for the first shard. One criteria for defining which shard to scaling independently is by looking into the metrics dashboard in the UI (e.g. https://cloud.mongodb.com/v2/<PROJECT-ID>#/clusters/detail/SymmetricShardedCluster)

```
resource "mongodbatlas_advanced_cluster" "test" {
  project_id   = var.project_id
  name         = "SymmetricShardedCluster"
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
