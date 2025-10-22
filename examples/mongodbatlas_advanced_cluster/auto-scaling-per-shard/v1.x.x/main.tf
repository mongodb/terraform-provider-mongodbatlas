provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

# Below is the old v1.x schema of mongodbatlas_advanced_cluster. 
# To migrate to v2.0.0+, see the main.tf in the parent directory. Refer README.md for more details.
resource "mongodbatlas_advanced_cluster" "this" {
  project_id   = mongodbatlas_project.project.id
  name         = "AutoScalingCluster"
  cluster_type = "SHARDED"
  disk_size_gb = 10 # removed in v2.0.0+, this can now be set per shard for inner specs

  # replication_specs are updated to a list of objects instead of blocks in v2.0.0+
  replication_specs { # first shard
    region_configs {  # region_configs are updated to a list of objects instead of blocks in v2.0.0+
      auto_scaling {  # auto_scaling are updated to an attribute instead of a block in v2.0.0+
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling { # analytics_auto_scaling are updated to an attribute instead of a block in v2.0.0+
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      electable_specs { # electable_specs are updated to an attribute instead of a block in v2.0.0+
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }

  replication_specs { # second shard
    region_configs {
      auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      analytics_auto_scaling {
        compute_enabled           = true
        compute_max_instance_size = "M60"
      }
      electable_specs {
        instance_size = "M40"
        node_count    = 3
      }
      analytics_specs {
        instance_size = "M40"
        node_count    = 1
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "Zone 1"
  }

  lifecycle { # avoid plans as autoscaling changes instance_size
    ignore_changes = [
      # in v2.0.0+, electable_specs and analytics_specs are updated to attributes and will no longer require index to access instance_size
      replication_specs[0].region_configs[0].electable_specs[0].instance_size,
      replication_specs[0].region_configs[0].analytics_specs[0].instance_size,
      replication_specs[1].region_configs[0].electable_specs[0].instance_size,
      replication_specs[1].region_configs[0].analytics_specs[0].instance_size,
    ]
  }
}

resource "mongodbatlas_project" "project" {
  name   = "AutoScalingPerShardCluster"
  org_id = var.atlas_org_id
}
