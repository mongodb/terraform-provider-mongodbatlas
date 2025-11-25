provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}


# Below is the old v1.x schema of mongodbatlas_advanced_cluster. 
# To migrate to v2.0.0+, see the main.tf in the parent directory. Refer README.md for more details.
resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = "ClusterToUpgrade"
  cluster_type = "REPLICASET"

  # replication_specs are updated to a list of objects instead of blocks in v2.0.0+
  replication_specs { # replication_specs are updated to a list of objects instead of blocks in v2.0.0+
    num_shards = 1    # removed in v2.0.0+

    region_configs { # region_configs are updated to a list of objects instead of blocks in v2.0.0+
      electable_specs {
        instance_size = var.provider_instance_size_name
        node_count    = 3
      }
      provider_name = var.provider_name
      region_name   = "US_EAST_1"
      priority      = 7
    }
  }

  tags { # tags and labels are updated to maps instead of blocks in v2.0.0+
    key   = "environment"
    value = "dev"
  }
}

resource "mongodbatlas_project" "project" {
  name   = "TenantUpgradeTest"
  org_id = var.atlas_org_id
}
