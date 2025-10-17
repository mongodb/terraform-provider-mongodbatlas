provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

# Below is the old v1.x schema of mongodbatlas_advanced_cluster. 
# To migrate to v2.0.0+, see the main.tf in the parent directory. Refer README.md for more details.
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id     = mongodbatlas_project.project.id
  name           = var.cluster_name
  cluster_type   = "SHARDED"
  backup_enabled = true
  disk_size_gb   = 10 # removed in v2.0.0+, this can now be set per shard for inner specs

  # replication_specs are updated to a list of objects instead of blocks in v2.0.0+
  replication_specs {

    # this cluster has 2 M30 shards
    # num_shards has been removed in v2.0.0+, upgrading to v2.0.0+ will require adding a new identical `replication_specs` element for each shard
    # refer https://registry.terraform.io/providers/mongodb/mongodbatlas/latest/docs/guides/advanced-cluster-new-sharding-schema for details.
    num_shards = 2

    region_configs {    # region_configs are updated to a list of objects instead of blocks in v2.0.0+
      electable_specs { # electable_specs are updated to an attribute instead of a block in v2.0.0+
        instance_size = "M30"
        disk_iops     = 3000
        node_count    = 3
      }
      provider_name = "AWS"
      priority      = 7
      region_name   = "EU_WEST_1"
    }
    zone_name = "zone n1"
  }
  advanced_configuration { # advanced_configuration are updated to an attribute instead of a block in v2.0.0+
    javascript_enabled                   = true
    oplog_size_mb                        = 999
    sample_refresh_interval_bi_connector = 300
  }

  tags { # tags and labels are updated to maps instead of blocks in v2.0.0+
    key   = "environment"
    value = "dev"
  }
}

resource "mongodbatlas_project" "project" {
  name   = "Symmetric Geosharded Cluster"
  org_id = var.atlas_org_id
}
