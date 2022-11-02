provider "mongodbatlas" {
  public_key  = var.public_key
  private_key = var.private_key
}

resource "mongodbatlas_advanced_cluster" "cluster" {
  project_id   = mongodbatlas_project.project.id
  name         = "ClusterToUpgrade"
  cluster_type = "REPLICASET"

  replication_specs {
    num_shards = 1

    region_configs {
      electable_specs {
        instance_size = var.provider_instance_size_name
      }
      provider_name         = var.provider_name
      backing_provider_name = var.backing_provider_name
      region_name           = "US_EAST_1"
      priority              = 7
    }
  }
}

resource "mongodbatlas_project" "project" {
  name   = "TenantUpgradeTest"
  org_id = var.atlas_org_id
}